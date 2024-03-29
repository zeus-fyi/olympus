package dynamic_secrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/tyler-smith/go-bip32"
	"github.com/wealdtech/go-ed25519hd"
	"github.com/zeus-fyi/gochain/web3/accounts"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	aegis_crypto "github.com/zeus-fyi/olympus/pkg/aegis/crypto"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	zeus_ecdsa "github.com/zeus-fyi/zeus/pkg/aegis/crypto/ecdsa"
)

var MaxZeros = 6

func GetAccount(val zeus_ecdsa.AddressGenerator) (accounts.Account, error) {
	pw := crypto.Keccak256Hash([]byte(val.Mnemonic)).Hex()
	seed, err := ed25519hd.SeedFromMnemonic(val.Mnemonic, pw)
	if err != nil {
		return accounts.Account{}, err
	}
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return accounts.Account{}, err
	}
	child, _ := masterKey.NewChildKey(uint32(val.PathIndex))
	privateKeyECDSA := crypto.ToECDSAUnsafe(child.Key)
	acc, err := accounts.CreateAccountFromPkey(privateKeyECDSA)
	if err != nil {
		return accounts.Account{}, err
	}
	tmp := acc.PrivateKey()
	if tmp == "" {
		return accounts.Account{}, errors.New("private key is nil")
	}
	return *acc, nil
}
func genAddresses(count int) (zeus_ecdsa.AddressGenerator, error) {
	numWorkers := runtime.NumCPU()
	addresses, err := aegis_crypto.GenAddresses(count, numWorkers)
	if err != nil {
		return zeus_ecdsa.AddressGenerator{}, err
	}
	if addresses.LeadingZeroesCount > MaxZeros {
		log.Info().Interface("address", addresses.Address).Msgf("found address with %d leading zeros", addresses.LeadingZeroesCount)
		MaxZeros = addresses.LeadingZeroesCount
		return addresses, nil
	}
	return zeus_ecdsa.AddressGenerator{}, errors.New("no addresses found")
}

func encAddress(age encryption.Age, ag zeus_ecdsa.AddressGenerator) (memfs.MemFS, filepaths.Path, error) {
	p := filepaths.Path{
		DirIn:  "",
		DirOut: "keygen",
	}
	fs := memfs.NewMemFs()
	key, err := json.Marshal(ag)
	if err != nil {
		return fs, p, err
	}

	name := fmt.Sprintf("key-%d.txt", ag.LeadingZeroesCount)
	p.FnOut = name
	err = age.EncryptItem(fs, &p, key)
	if err != nil {
		return fs, p, err
	}
	encOut, err := fs.ReadFileOutPath(&p)
	if err != nil {
		return fs, p, err
	}
	if encOut == nil {
		return fs, p, err
	}
	return fs, p, err
}

func SaveAddress(ctx context.Context, tries int, s3Client s3base.S3Client, age encryption.Age) error {
	ag, err := genAddresses(tries)
	if err != nil {
		return err
	}
	log.Info().Interface("address", ag.Address).Msgf("found address with %d leading zeros", ag.LeadingZeroesCount)
	fs, p, err := encAddress(age, ag)
	if err != nil {
		return err
	}
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(p.FnOut),
	}
	if p.FnOut == "key-4.txt.age" {
		MaxZeros += 1
		return errors.New("key already exists")
	}
	uploader := s3uploader.NewS3ClientUploader(s3Client)
	exists, err := uploader.CheckIfKeyExists(ctx, input)
	if err != nil {
		return err
	}
	if exists {
		MaxZeros += 1
		return errors.New("key already exists")
	}
	err = uploader.UploadFromInMemFs(ctx, p, input, fs)
	if err != nil {
		return err
	}
	return nil
}

func ReadAddress(ctx context.Context, p filepaths.Path, s3Client s3base.S3Client, age encryption.Age) (zeus_ecdsa.AddressGenerator, error) {
	download := s3reader.NewS3ClientReader(s3Client)
	bucketName := "zeus-fyi"
	getObj := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(p.FnIn),
	}
	ag := zeus_ecdsa.AddressGenerator{}
	b := download.ReadBytes(ctx, &p, getObj)
	if b == nil {
		return ag, errors.New("no addresses found")
	}
	fs := memfs.NewMemFs()
	p.DirIn = "keygen"
	err := fs.MakeFileIn(&p, b.Bytes())
	if err != nil {
		return ag, err
	}
	p.DirOut = "keygen"
	err = age.DecryptToMemFsFile(&p, fs)
	if err != nil {
		return ag, err
	}

	out, err := fs.ReadFileOutPath(&p)
	if err != nil {
		return ag, err
	}
	err = json.Unmarshal(out, &ag)
	if err != nil {
		return ag, err
	}
	return ag, nil
}
