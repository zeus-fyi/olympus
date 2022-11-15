package poseidon

import "github.com/rs/zerolog/log"

func (p *Poseidon) GzipCompressChainData() error {
	err := p.CreateTarGzipArchiveDir(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: GzipChainData")
		return err
	}
	return err
}

func (p *Poseidon) ZstdCompressChainData() error {
	err := p.CreateTarZstdArchiveDir(&p.Path)
	if err != nil {
		log.Err(err).Msg("Poseidon: CreateTarZstdArchiveDir")
		return err
	}
	return err
}
