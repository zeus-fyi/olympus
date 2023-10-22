package v1_iris

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
)

func (s *IrisV1TestSuite) TestWebsocket() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s.Eg.GET("/mempool", mempoolWebSocketHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = s.E.Start(":9010")
	}()
	var addr = flag.String("addr", "localhost:9010", "ws service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/v1/mempool"}
	ws, _, werr := websocket.DefaultDialer.Dial(u.String(), nil)
	s.Require().Nil(werr)
	defer ws.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			s.Require().Nil(err)
			fmt.Println(message)
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	// todo needs to send expected data to websocket via redis or mock
	go SendDataToWebSocket()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := ws.WriteMessage(websocket.TextMessage, []byte(t.String()))
			s.NoError(err)
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			s.NoError(err)
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (s *IrisV1TestSuite) TestLiveWebsocket() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//var addr = flag.String("addr", "localhost:8080", "ws service address")

	var addr = flag.String("addr", "localhost:8080", "ws service address")

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/v1/mempool"}

	requestHeader := http.Header{}
	requestHeader.Add("Authorization", "Bearer "+s.Tc.ProductionLocalTemporalBearerToken)

	ws, _, werr := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	s.Require().Nil(werr)
	defer ws.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			s.Require().Nil(err)
			fmt.Println(message)
			tx := types.Transaction{}
			err = tx.UnmarshalBinary(message)
			s.Require().Nil(err)

			log.Printf("recv: %s", tx.Hash().String())
			fmt.Println(tx)
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := ws.WriteMessage(websocket.TextMessage, []byte(t.String()))
			s.NoError(err)
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			s.NoError(err)
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
