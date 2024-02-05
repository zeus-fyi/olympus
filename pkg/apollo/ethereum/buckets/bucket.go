package apollo_buckets

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

var ApolloS3Manager s3base.S3Client

var EthMainnetBucket = BucketRequest{
	Protocol: "ethereum",
	Network:  "mainnet",
}

type BucketRequest struct {
	BucketName string `json:"bucketName"`
	BucketKey  string `json:"bucketKey"`

	Protocol string `json:"protocol"`
	Network  string `json:"network"`
}

func (b *BucketRequest) GetBucketName() string {
	key := []string{"apollo", strings.ToLower(b.Protocol), strings.ToLower(b.Network)}
	return strings.Join(key, "-")
}

func UploadBalancesAtEpoch(ctx context.Context, keyName string, balances []byte) error {
	// for now just hard coded
	ctx = context.WithValue(ctx, "func", "UploadBalanceAtEpoch")
	br := EthMainnetBucket
	p := filepaths.Path{
		PackageName: "",
		DirIn:       ".",
		DirOut:      "./out",
		FnIn:        keyName + ".json",
	}
	uploader := s3uploader.NewS3ClientUploader(ApolloS3Manager)
	input := &s3.PutObjectInput{
		Bucket: aws.String(br.GetBucketName()),
		Key:    aws.String(keyName + ".json.tar.lz4"),
	}

	ok, err := uploader.CheckIfKeyExists(ctx, input)
	if ok || err != nil {
		if err != nil {
			log.Err(err).Msg("UploadBalancesAtEpoch")
		}
		return nil
	}
	// upload to spaces
	inMemFs := memfs.NewMemFs()
	err = inMemFs.MakeFileIn(&p, balances)
	if err != nil {
		log.Err(err).Msgf("UploadBalancesAtEpoch: MakeFileIn %s", keyName)
		return err
	}
	comp := compression.NewCompression()
	inMemFs, err = comp.Lz4CompressInMemFsFile(&p, inMemFs)
	if err != nil {
		log.Err(err).Msgf("UploadBalancesAtEpoch: MakeFileIn %s", keyName)
		return err
	}

	err = uploader.UploadFromInMemFs(ctx, p, input, inMemFs)
	if err != nil {
		log.Err(err).Msgf("UploadBalancesAtEpoch: UploadFromInMemFs %s", keyName)
		return err
	}
	return err
}

func DownloadBalancesAtEpoch(ctx context.Context, keyName string) ([]byte, error) {
	// for now just hard coded
	p := &filepaths.Path{
		PackageName: "",
		DirIn:       ".",
		DirOut:      ".",
		FnIn:        keyName + ".json.tar.lz4",
		FnOut:       keyName,
	}
	// upload to spaces
	inMemFs := memfs.NewMemFs()
	comp := compression.NewCompression()

	br := EthMainnetBucket
	ctx = context.WithValue(ctx, "func", "DownloadBalancesAtEpoch")
	downloader := s3reader.NewS3ClientReader(ApolloS3Manager)
	input := &s3.GetObjectInput{
		Bucket: aws.String(br.GetBucketName()),
		Key:    aws.String(p.FnIn),
	}
	buf, err := downloader.ReadBytesNoPanic(ctx, p, input)
	if err != nil {
		log.Err(err).Msgf("DownloadBalancesAtEpoch: UploadFromInMemFs %s", keyName)
		return nil, err
	}
	if buf.Len() <= 0 {
		err = errors.New("no key found")
		return nil, err
	}
	err = inMemFs.MakeFileIn(p, buf.Bytes())
	if err != nil {
		log.Err(err).Msgf("DownloadBalancesAtEpoch: UploadFromInMemFs %s", keyName)
		return nil, err
	}

	ho, err := downloader.GetHeadObject(ctx, input)
	if err != nil {
		log.Err(err).Msgf("DownloadBalancesAtEpoch: GetHeadObject %s", keyName)
		return nil, err
	}

	p.Metadata = ho.Metadata
	inMemFs, err = comp.Lz4DecompressInMemFsFile(p, inMemFs)
	if err != nil {
		log.Err(err).Msgf("DownloadBalancesAtEpoch: UploadFromInMemFs %s", keyName)
		return nil, err
	}
	b, err := inMemFs.ReadFile(p.FileInPath())
	if err != nil {
		log.Err(err).Msgf("DownloadBalancesAtEpoch: UploadFromInMemFs %s", keyName)
		return nil, err
	}
	return b, err
}
