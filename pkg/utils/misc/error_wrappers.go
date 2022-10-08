package misc

import "github.com/rs/zerolog/log"

func ReturnIfErr(err error, errMsg string) error {
	if err != nil {
		log.Err(err).Msg(errMsg)
		return err
	}
	return nil
}
