package v1_ethereum_aws

import (
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
}
