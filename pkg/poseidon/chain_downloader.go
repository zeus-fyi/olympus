package poseidon

import "github.com/rs/zerolog/log"

func (p *Poseidon) GzipDecChainData() error {
	err := p.GzipDecompress(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: GzipDecChainData")
		return err
	}
	return err
}

func (p *Poseidon) ZstdDecChainData() error {
	err := p.ZstdDecompress(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: ZstdDecChainData")
		return err
	}
	return err
}
