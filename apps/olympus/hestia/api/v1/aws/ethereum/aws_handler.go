package v1_aws_ethereum_automation

import (
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
}
