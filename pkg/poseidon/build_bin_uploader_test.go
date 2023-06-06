package poseidon

import "context"

func (s *ChainUploaderTestSuite) TestBuildBinUploader() {
	ctx := context.Background()
	pos := NewPoseidon(s.S3)
	pos.DirIn = "./app/path"
	pos.DirOut = "./build/bin"
	pos.FnIn = "appName"
	err := pos.UploadsBinCompressBuild(ctx, pos.FnIn)
	s.Require().Nil(err)
}
