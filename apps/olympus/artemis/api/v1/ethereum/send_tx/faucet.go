package artemis_eth_txs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LK4D4/trylock"
	"github.com/gochain/gochain/v4/common"
	"github.com/jellydator/ttlcache/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
)

var Faucet *FaucetServer

type FaucetConfig struct {
	network    string
	httpPort   int
	interval   int
	payout     int
	proxyCount int
	queueCap   int
}

func NewFaucetConfig(network string, httpPort, interval, payout, proxyCount, queueCap int) *FaucetConfig {
	return &FaucetConfig{
		network:    network,
		httpPort:   httpPort,
		interval:   interval,
		payout:     payout,
		proxyCount: proxyCount,
		queueCap:   queueCap,
	}
}

type claimRequest struct {
	Address string `json:"address"`
}

type claimResponse struct {
	Message string `json:"msg"`
}

type infoResponse struct {
	Account string `json:"account"`
	Network string `json:"network"`
	Payout  string `json:"payout"`
}

func decodeJSONBody(r *http.Request, dst interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1024))
	defer r.Body.Close()
	if err != nil {
		return &malformedRequest{status: http.StatusBadRequest, message: "Unable to read request body"}
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, message: msg}
		default:
			return err
		}
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	return nil
}

func readAddress(r *http.Request) (string, error) {
	var claimReq claimRequest
	if err := decodeJSONBody(r, &claimReq); err != nil {
		return "", err
	}
	if !IsValidAddress(claimReq.Address, true) {
		return "", &malformedRequest{status: http.StatusBadRequest, message: "invalid address"}
	}

	return claimReq.Address, nil
}

func renderJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func IsValidAddress(address string, checksummed bool) bool {
	if !common.IsHexAddress(address) {
		return false
	}
	return !checksummed || common.HexToAddress(address).Hex() == address
}

type Limiter struct {
	mutex      sync.Mutex
	cache      *ttlcache.Cache
	proxyCount int
	ttl        time.Duration
}

func NewLimiter(proxyCount int, ttl time.Duration) *Limiter {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)
	return &Limiter{
		cache:      cache,
		proxyCount: proxyCount,
		ttl:        ttl,
	}
}

func (l *Limiter) Error() {

	//var mr *malformedRequest
	//if errors.As(err, &mr) {
	//	renderJSON(w, claimResponse{Message: mr.message}, mr.status)
	//} else {
	//	renderJSON(w, claimResponse{Message: http.StatusText(http.StatusInternalServerError)}, http.StatusInternalServerError)
	//}
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	address, err := readAddress(r)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			renderJSON(w, claimResponse{Message: mr.message}, mr.status)
		} else {
			renderJSON(w, claimResponse{Message: http.StatusText(http.StatusInternalServerError)}, http.StatusInternalServerError)
		}
		return
	}

	if l.ttl <= 0 {
		next.ServeHTTP(w, r)
		return
	}

	clintIP := getClientIPFromRequest(l.proxyCount, r)
	l.mutex.Lock()
	if l.limitByKey(w, address) || l.limitByKey(w, clintIP) {
		l.mutex.Unlock()
		return
	}
	l.cache.SetWithTTL(address, true, l.ttl)
	l.cache.SetWithTTL(clintIP, true, l.ttl)
	l.mutex.Unlock()

	next.ServeHTTP(w, r)
	if w.(negroni.ResponseWriter).Status() != http.StatusOK {
		l.cache.Remove(address)
		l.cache.Remove(clintIP)
		return
	}
	log.Info().Interface("address", address).Interface("clientIP", clintIP).Msg("Maximum request limit has been reached")
}

func (l *Limiter) limitByKey(w http.ResponseWriter, key string) bool {
	if _, ttl, err := l.cache.GetWithTTL(key); err == nil {
		errMsg := fmt.Sprintf("You have exceeded the rate limit. Please wait %s before you try again", ttl.Round(time.Second))
		renderJSON(w, claimResponse{Message: errMsg}, http.StatusTooManyRequests)
		return true
	}
	return false
}

func getClientIPFromRequest(proxyCount int, r *http.Request) string {
	if proxyCount > 0 {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			xForwardedForParts := strings.Split(xForwardedFor, ",")
			// Avoid reading the user's forged request header by configuring the count of reverse proxies
			partIndex := len(xForwardedForParts) - proxyCount
			if partIndex < 0 {
				partIndex = 0
			}
			return strings.TrimSpace(xForwardedForParts[partIndex])
		}
	}
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}

type FaucetServer struct {
	mutex trylock.Mutex
	cfg   *FaucetConfig
	queue chan string
}

func NewFaucetServer() *FaucetServer {
	cfg := NewFaucetConfig("ephemeral", 9005, 1440, 32, 0, 1000)
	return &FaucetServer{
		cfg:   cfg,
		queue: make(chan string, cfg.queueCap),
	}
}

func (s *FaucetServer) Run() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			s.consumeQueue()
		}
	}()

	//limiter := NewLimiter(s.cfg.proxyCount, time.Duration(s.cfg.interval)*time.Minute)

}

func (s *FaucetServer) consumeQueue() {
	if len(s.queue) == 0 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for len(s.queue) != 0 {
		address := <-s.queue
		sendEthTransferPayload := web3_actions.TransferArgs{
			Amount:    EtherToWei(int64(s.cfg.payout)),
			ToAddress: common.StringToAddress(address),
		}
		sendEthPayload := web3_actions.SendEtherPayload{
			TransferArgs:   sendEthTransferPayload,
			GasPriceLimits: web3_actions.GasPriceLimits{},
		}
		ctx := context.Background()
		err := artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.ExecuteArtemisSendEthTxWorkflow(ctx, sendEthPayload)
		log.Info().Interface("toAddress", address).Msg("send eth tx in progress")
		if err != nil {
			log.Err(err).Msg("Failed to handle transaction in the queue")
		} else {
			log.Info().Interface("txHash", "txHash").Interface("address", address).Msg("Consume from queue successfully")
		}
	}
}

type malformedRequest struct {
	status  int
	message string
}

func (mr *malformedRequest) Error() string {
	return mr.message
}
func (s *FaucetServer) FaucetHandler(c echo.Context) error {
	request := new(claimRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	if !IsValidAddress(request.Address, true) {
		return c.JSON(http.StatusBadRequest, "invalid address")
	}
	address := request.Address
	// Try to lock mutex if the work queue is empty
	if len(s.queue) != 0 || !s.mutex.TryLock() {
		select {
		case s.queue <- address:
			log.Info().Interface("address", address).Msg("Added to queue successfully")
			resp := claimResponse{Message: fmt.Sprintf("Added %s to the queue", address)}
			return c.JSON(http.StatusAccepted, resp)
		default:
			log.Warn().Msg("Max queue capacity reached")
			resp := claimResponse{Message: "Faucet queue is too long, please try again later"}
			return c.JSON(http.StatusServiceUnavailable, resp)
		}
	}

	sendEthTransferPayload := web3_actions.TransferArgs{
		Amount:    EtherToWei(int64(s.cfg.payout)),
		ToAddress: common.StringToAddress(address),
	}
	sendEthPayload := web3_actions.SendEtherPayload{
		TransferArgs:   sendEthTransferPayload,
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	err := artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.ExecuteArtemisSendEthTxWorkflow(context.Background(), sendEthPayload)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.mutex.Unlock()
	if err != nil {
		log.Err(err).Msg("Failed to send transaction")
		resp := claimResponse{Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, resp)
	}

	log.Ctx(ctx).Info().Interface(
		"txHash", "txHash").Interface("address", address).Msg("Funded directly successfully")
	resp := claimResponse{Message: fmt.Sprintf("Txhash: %s", "")}
	return c.JSON(http.StatusAccepted, resp)
}

func EtherToWei(amount int64) *big.Int {
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return new(big.Int).Mul(big.NewInt(amount), ether)
}

func (s *FaucetServer) handleInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		renderJSON(w, infoResponse{
			Account: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
			Network: s.cfg.network,
			Payout:  strconv.Itoa(s.cfg.payout),
		}, http.StatusOK)
	}
}
