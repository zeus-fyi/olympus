/*
 * Web3Signer ETH2 Api
 *
 * Sign Eth2 Artifacts
 *
 * API version: @VERSION@
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package consensys_eth2_openapi

type BeaconBlockSigning struct {
	Type string `json:"type"`

	ForkInfo SigningForkInfo `json:"fork_info"`

	SigningRoot string `json:"signingRoot,omitempty"`

	BeaconBlock any `json:"beacon_block"`
}
