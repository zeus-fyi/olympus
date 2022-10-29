package keys

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

const AgeKeyTypeID = 1
const PgpKeyTypeID = 2
const EcdsaKeyTypeID = 3
const BlsKeyTypeID = 4
const BearerKeyTypeID = 5
const JwtKeyTypeID = 6

type KeyType struct {
	autogen_bases.KeyTypes
}
