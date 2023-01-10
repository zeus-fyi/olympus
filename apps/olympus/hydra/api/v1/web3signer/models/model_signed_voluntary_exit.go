/*
 * Web3Signer ETH2 Api
 *
 * Sign Eth2 Artifacts
 *
 * API version: @VERSION@
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package consensys_eth2_openapi

type SignedVoluntaryExit struct {
	Message VoluntaryExit `json:"message,omitempty"`

	Signature string `json:"signature,omitempty"`
}
