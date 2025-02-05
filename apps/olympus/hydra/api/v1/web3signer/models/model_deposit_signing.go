/*
 * Web3Signer ETH2 Api
 *
 * Sign Eth2 Artifacts
 *
 * API version: @VERSION@
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package consensys_eth2_openapi

type DepositSigning struct {
	Type string `json:"type"`

	// signing root for optional verification if field present
	SigningRoot string `json:"signingRoot,omitempty"`

	Deposit DepositSigningDeposit `json:"deposit"`
}
