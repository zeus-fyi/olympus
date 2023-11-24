package misc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
)

func DelayedPanic(err error) {
	log.Err(err).Msg("DelayedPanic")
	time.Sleep(10 * time.Second)
	panic(err)
}

func HashParams(hashParams []interface{}) (string, error) {
	hash := sha3.New256()
	for _, v := range hashParams {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		_, _ = hash.Write(b)
	}
	// Get the resulting encoded byte slice
	sha3v := hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(sha3v)), nil
}
