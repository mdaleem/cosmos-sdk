package types

// msg_authorization module event types
const (
	EventGrantAuthorization   = "grant-authorization"
	EventRevokeAuthorization  = "revoke-authorization"
	EventExecuteAuthorization = "execute-authorization"

	AttributeKeyGranteeAddress = "grantee"
	AttributeKeyGranterAddress = "granter"

	AttributeValueCategory = ModuleName
)
