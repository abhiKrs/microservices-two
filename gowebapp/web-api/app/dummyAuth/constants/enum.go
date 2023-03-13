package constants

const (
	Users  DBUserRole = iota // EnumIndex = 0
	Admins                   // EnumIndex = 1
)

const (
	Custom   DBAuthType = iota + 1 // EnumIndex = 1
	Firebase                       // EnumIndex = 2
)

const (
	DBMagicLink DBProviderType = iota + 1 // EnumIndex = 1
	DBPassword                            // EnumIndex = 2
	DBGoogle                              // EnumIndex = 3
	DBGithub                              // EnumIndex = 4
)

const (
	Register      DBMagicLinkType = iota + 1 // EnumIndex = 1
	Login                                    // EnumIndex = 2
	ResetPassword                            // EnumIndex = 3
)

const (
	Bcrypt DBPasswordAlgo = iota + 1 // EnumIndex = 1
)
