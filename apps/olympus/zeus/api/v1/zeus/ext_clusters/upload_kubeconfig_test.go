package zeus_v1_clusters_api

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/compression"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var ctx = context.Background()

type KubeConfigRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *KubeConfigRequestTestSuite) TestKubeConfigUpload() {
	t.Eg.POST("/kubeconfig", CreateOrUpdateKubeConfigsHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	fp := filepaths.Path{
		DirIn:  "/Users/alex/go/Olympus/olympus/apps/olympus/zeus/api/v1/zeus/ext_clusters/.kube",
		DirOut: "/Users/alex/go/Olympus/olympus/apps/olympus/zeus/api/v1/zeus/ext_clusters/",
		FnIn:   "kube",
	}
	err := ZipKubeConfigChartToPath(&fp)
	t.Require().Nil(err)
	authKeysCfg := t.Tc.ProdLocalAuthKeysCfg
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)

	pubKey := t.Tc.LocalAgePubkey
	privKey := t.Tc.LocalAgePkey
	AgeEnc = encryption.NewAge(privKey, pubKey)

	resp, err := t.ZeusClient.R().
		SetFormData(map[string]string{}).
		SetFile("kubeconfig", fp.FileOutPath()).
		Post("/v1/kubeconfig")
	t.Require().Nil(err)
	t.Require().Equal(200, resp.StatusCode())
}

func ZipKubeConfigChartToPath(p *filepaths.Path) error {
	comp := compression.NewCompression()
	err := comp.GzipCompressDir(p)
	if err != nil {
		log.Err(err).Interface("path", p).Msg("ZeusClient: ZipKubeConfigChartToPath")
		return err
	}
	return err
}

func (t *KubeConfigRequestTestSuite) TestS3EncUploader() {
	t.InitLocalConfigs()
	authKeysCfg := t.Tc.ProdLocalAuthKeysCfg
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	pubKey := t.Tc.LocalAgePubkey
	privKey := t.Tc.LocalAgePkey
	ageEnc := encryption.NewAge(privKey, pubKey)
	in := bytes.Buffer{}
	tmp := []byte("test")
	in.Write(tmp)
	err := EncAndUpload(ctx, t.Tc.ProductionLocalTemporalOrgID, in, ageEnc)
	t.Require().Nil(err)
}

func (t *KubeConfigRequestTestSuite) TestKubeConfigAccess() {
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, t.Tc.ProdLocalAuthKeysCfg)
	t.Tc.ProdLocalAuthKeysCfg.AgePrivKey = t.Tc.LocalAgePkey
	t.Tc.ProdLocalAuthKeysCfg.AgePubKey = t.Tc.LocalAgePubkey
	authCfg := auth_startup.NewExtClustersAuthClient(ctx, t.Tc.ProductionLocalTemporalOrgID, t.Tc.ProdLocalAuthKeysCfg)
	t.Require().NotEmpty(authCfg)

	inMemFs := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(context.Background(), t.Tc.ProductionLocalTemporalOrgID, authCfg)
	k := zeus_core.K8Util{}
	k.ConnectToK8sFromInMemFsCfgPath(inMemFs)

	rawCfg, err := k.GetRawConfigs()
	t.Require().Nil(err)
	t.Require().NotEmpty(rawCfg)

	m := make(map[string]string)

	for ctxName, ai := range rawCfg.Clusters {
		if strings.Contains(ai.Server, "aws") {
			m[ctxName] = "aws"
			fmt.Println("aws command found")
			continue
		}
		if strings.Contains(ai.Server, "digtalocean") {
			m[ctxName] = "do"
			fmt.Println("digital ocean command found")
			continue
		}
		if strings.Contains(ai.Server, "ovh") {
			fmt.Println("ovh server found")
			m[ctxName] = "ovh"
			continue
		}
		if strings.Contains(ai.Server, "gke") || strings.Contains(ctxName, "gke") || strings.Contains(ctxName, "gcp") {
			m[ctxName] = "gcp"
			fmt.Println("gcp command found")
			continue
		}
	}

	ctxNames, err := k.GetContexts()
	t.Require().Nil(err)
	t.Require().NotEmpty(ctxNames)
	for ctxName, _ := range ctxNames {
		fmt.Println(ctxName)
	}
}

func TestKubeConfigRequestTestSuite(t *testing.T) {
	suite.Run(t, new(KubeConfigRequestTestSuite))
}
