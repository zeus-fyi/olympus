package keys

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

const (
	AgeKeyTypeID             = 0
	GpgKeyTypeID             = 1
	PgpKeyTypeID             = 2
	EcdsaKeyTypeID           = 3
	BlsKeyTypeID             = 4
	BearerKeyTypeID          = 5
	JwtKeyTypeID             = 6
	PassphraseKeyTypeID      = 7
	SessionIDKeyTypeID       = 8
	VerifyEmailTokenTypeID   = 9
	PasswordResetTokenTypeID = 10
	StripeCustomerID         = 11
	QuickNodeCustomerID      = 12
	DiscordApiKeyID          = 13
)

type KeyType struct {
	autogen_bases.KeyTypes
}

func NewBearerKeyType() KeyType {
	return KeyType{autogen_bases.KeyTypes{
		KeyTypeID:   BearerKeyTypeID,
		KeyTypeName: "",
	}}
}

func NewQuickNodeKeyType(quickNodeID string) KeyType {
	return KeyType{autogen_bases.KeyTypes{
		KeyTypeID:   QuickNodeCustomerID,
		KeyTypeName: quickNodeID,
	}}
}
