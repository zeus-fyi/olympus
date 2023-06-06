package hephaestus_build_actions

import (
	"context"

	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Config struct {
	Host      string
	Port      string
	Name      string
	PGConnStr string
}

func Upload(ctx context.Context) {
	a := poseidon_orchestrations.PoseidonS3Manager
	sb := s3base.S3Client{
		AwsS3Client:    a.AwsS3Client,
		SpacesEndpoint: a.SpacesEndpoint,
	}
	newFp := filepaths.Path{
		PackageName: appName,
		DirIn:       dataDir.DirIn,
		DirOut:      dataDir.DirOut,
		FnIn:        dataDir.FnIn,
		FnOut:       dataDir.FnOut,
		Env:         env,
		FilterFiles: string_utils.FilterOpts{},
	}
	pos := poseidon.NewPoseidon(sb)
	pos.Path = newFp
	err := pos.UploadsBinCompressBuild(ctx, appName)
	if err != nil {
		log.Fatal().Err(err).Msg("Poseidon: UploadsBinCompressBuild failed")
		misc.DelayedPanic(err)
	}
	return
}
