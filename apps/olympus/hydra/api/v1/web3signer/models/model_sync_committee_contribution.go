/*
 * Web3Signer ETH2 Api
 *
 * Sign Eth2 Artifacts
 *
 * API version: @VERSION@
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package consensys_eth2_openapi

type SyncCommitteeContribution struct {
	Slot string `json:"slot,omitempty"`

	// Bytes32 hexadecimal
	BeaconBlockRoot string `json:"beacon_block_root,omitempty"`

	SubcommitteeIndex string `json:"subcommittee_index,omitempty"`

	// SSZ hexadecimal
	AggregationBits string `json:"aggregation_bits,omitempty"`

	// Bytes96 hexadecimal
	Signature string `json:"signature,omitempty"`
}
