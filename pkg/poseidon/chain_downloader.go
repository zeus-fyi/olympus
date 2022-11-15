package poseidon

import "github.com/rs/zerolog/log"

func (p *Poseidon) UnGzipChainData() error {
	err := p.UnGzip(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: UnGzipChainData")
		return err
	}
	return err
}
