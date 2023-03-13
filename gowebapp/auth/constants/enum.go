package constants

const (
	MagicLink AuthType = iota + 1 // EnumIndex = 1
	Password                      // EnumIndex = 2
	Google
)

// Declare related constants for each EnvConst starting with index 1
const (
	Development EnvConst = iota + 1 // EnumIndex = 1
	Test                            // EnumIndex = 2
	Staging
	Production
)

// Declare related constants for each EmailType starting with index 1
const (
	NewMagicLink EmailType = iota + 1 // EnumIndex = 1
	SuccessMessage
)

// Declare related constants for each SuccessType starting with index 1
const (
	PasswordReset SuccessType = iota + 1 // EnumIndex = 1
	Onboarding
	PasswordSet
)
