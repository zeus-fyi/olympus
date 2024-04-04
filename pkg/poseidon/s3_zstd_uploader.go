package poseidon

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (p *Poseidon) S3ZstdCompressAndUpload(ctx context.Context, br S3BucketRequest) error {
	ctx = context.WithValue(ctx, "func", "ZstdCompressAndUpload")
	err := p.ZstCompressFile(ctx, &p.Path)
	if err != nil {
		log.Err(err).Msg("S3ZstdCompressAndUpload: ZstCompressFile")
		return err
	}
	err = p.S3UploadSnapshot(ctx, br)
	return err
}
