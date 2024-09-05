package passportrootupdate

import rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"

type PassportRootTransferDetails struct {
	Operation *rarimocore.PassportRootUpdate
	Proof     []byte
}
