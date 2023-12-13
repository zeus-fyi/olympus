package hera_reddit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var ctx = context.Background()

type RedditTestSuite struct {
	test_suites_base.TestSuite
	rc Reddit
}

func (s *RedditTestSuite) SetupTest() {
	s.InitLocalConfigs()

	oauthPub := s.Tc.RedditPublicOAuth2
	oauthSec := s.Tc.RedditSecretOAuth2
	rc, err := InitRedditClient(ctx, oauthPub, oauthSec, s.Tc.RedditUsername, s.Tc.RedditPassword)
	s.Require().Nil(err)
	s.Assert().NotNil(rc)
	s.rc = rc
}

func (s *RedditTestSuite) TestInitOrgRedditClient() {
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}

	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	ou := org_users.NewOrgUserWithID(s.Tc.ProductionLocalTemporalOrgID, s.Tc.ProductionLocalTemporalUserID)
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, ou, "reddit")
	s.Require().Nil(err)
	rc, err := InitOrgRedditClient(ctx, ps.OAuth2Public, ps.OAuth2Secret, ps.Username, ps.Password)
	s.Require().Nil(err)
	s.Assert().NotNil(rc)
	lpo := &reddit.ListOptions{
		Limit:  10,
		After:  "1829cc6",
		Before: "",
	}
	posts, _, err := rc.ReadOnly.Subreddit.NewPosts(ctx, "mlops", lpo)
	s.Require().Nil(err)
	s.Assert().NotZero(posts)
}

func (s *RedditTestSuite) TestReadPosts() {
	lpo := &reddit.ListOptions{
		Limit:  10,
		After:  "1829cc6",
		Before: "",
	}
	posts, _, err := s.rc.ReadOnly.Subreddit.NewPosts(ctx, "ethdev", lpo)
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	searchID := 1700890336466260000

	resp, err := hera_search.InsertIncomingRedditPosts(ctx, searchID, posts)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}

func TestRedditTestSuite(t *testing.T) {
	suite.Run(t, new(RedditTestSuite))
}
