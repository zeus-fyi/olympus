package athena_client

func (t *AthenaClientTestSuite) TestHostDiskInfo() {
	resp, err := t.AthenaTestClient.GetHostDiskInfo(ctx)
	t.Assert().Nil(err)
	t.Assert().NotNil(resp)
}

func (t *AthenaClientTestSuite) TestHostMemInfo() {
	resp, err := t.AthenaTestClient.GetHostMemInfo(ctx)
	t.Assert().Nil(err)
	t.Assert().NotNil(resp)
}
