package constants

// Declare related constants for each EnvConst starting with index 0
const (
	Development EnvConst = iota // EnumIndex = 0
	Test                        // EnumIndex = 1
	Staging                     // EnumIndex = 2
	Production                  // EnumIndex = 3
)

// Declare related constants for each EmailType starting with index 0
const (
	MagicLinkMessage EmailType = iota // EnumIndex = 0
	SuccessMessage                    // EnumIndex = 1
)

// Declare related constants for each SuccessType starting with index 0
const (
	PasswordReset SuccessType = iota // EnumIndex = 0
	Onboarding                       // EnumIndex = 1
	PasswordSet                      // EnumIndex = 2
)

// Declare related constants for each AuthType starting with index 0
const (
	MagicLink AuthType = iota + 1 // EnumIndex = 1
	Password                      // EnumIndex = 2
	Google                        // EnumIndex = 3
)

const (
	Register      DBMagicLinkType = iota + 1 // EnumIndex = 1
	Login                                    // EnumIndex = 2
	ResetPassword                            // EnumIndex = 3
)
