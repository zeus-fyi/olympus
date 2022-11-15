package poseidon

import "github.com/rs/zerolog/log"

func (p *Poseidon) GzipChainData() error {
	err := p.CreateTarGzipArchiveDir(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: GzipChainData")
		return err
	}
	return err
}
