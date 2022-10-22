package v0

import "github.com/rs/zerolog/log"

func (l *LibV0) ErrHandler(e error) error {
	log.Error().Err(e).Int("UNIX UTC Time: ", l.UnixTimeStampNow())
	return e
}
