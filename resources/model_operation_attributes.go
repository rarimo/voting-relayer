/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type OperationAttributes struct {
	BlockHeight uint64 `json:"block_height"`
	// Address of the contract to which the transaction data should be sent
	DestinationAddress string `json:"destination_address"`
	// Destination chain ID
	DestinationChain string `json:"destination_chain"`
	OperationId      string `json:"operation_id"`
	Proof            string `json:"proof"`
	// Serialized transaction data
	TxHash string `json:"tx_hash"`
}
