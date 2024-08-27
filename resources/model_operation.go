/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Operation struct {
	Key
	Attributes OperationAttributes `json:"attributes"`
}
type OperationRequest struct {
	Data     Operation `json:"data"`
	Included Included  `json:"included"`
}

type OperationListRequest struct {
	Data     []Operation `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustOperation - returns Operation from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustOperation(key Key) *Operation {
	var operation Operation
	if c.tryFindEntry(key, &operation) {
		return &operation
	}
	return nil
}
