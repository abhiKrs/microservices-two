package constants

const (
	Users  DBUserRole = iota // EnumIndex = 1
	Admins                   // EnumIndex = 2
)

const (
	Custom   DBAuthType = iota + 1 // EnumIndex = 1
	Firebase                       // EnumIndex = 2
)

const (
	DBMagicLink DBProviderType = iota + 1 // EnumIndex = 1
	DBPassword
	DBGoogle
	DBGithub
)

const (
	Register DBMagicLinkType = iota + 1 // EnumIndex = 1
	Login                               // EnumIndex = 2
	ResetPassword
)

const (
	Bcrypt DBPasswordAlgo = iota + 1 // EnumIndex = 1
)
