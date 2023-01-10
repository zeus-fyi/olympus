/*
 * Web3Signer ETH2 Api
 *
 * Sign Eth2 Artifacts
 *
 * API version: @VERSION@
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package consensys_eth2_openapi

type ImportKeystoresResponseData struct {

	// - imported: Keystore successfully decrypted and imported to keymanager permanent storage - duplicate: Keystore's pubkey is already known to the keymanager - error: Any other status different to the above: decrypting error, I/O errors, etc.
	Status string `json:"status"`

	// error message if status == error
	Message string `json:"message,omitempty"`
}
