package constants

// Enum AuthType
type AuthType uint8

func (d AuthType) String() string {
	strings := [...]string{
		"magicLink",
		"password",
		"google",
	}

	if d < MagicLink || d > Google {
		return "Unknown"
	}
	return strings[d-1]
}

func (d AuthType) EnumIndex() int {
	return int(d)
}

// Enum EnvConst
type EnvConst int

func (d EnvConst) String() string {
	strings := [...]string{
		"development",
		"test",
		"staging",
		"production",
	}

	if d < Development || d > Production {
		return "Unknown"
	}
	return strings[d]
}

func (d EnvConst) EnumIndex() int {
	return int(d)
}

// Enum EmailType
type EmailType int

func (d EmailType) String() string {
	strings := [...]string{
		"magicLinkMessage",
		"successMessage",
	}

	if d < MagicLinkMessage || d > SuccessMessage {
		return "Unknown"
	}
	return strings[d]
}

func (d EmailType) EnumIndex() int {
	return int(d)
}

// Enum SuccessType
type SuccessType int

func (d SuccessType) String() string {
	strings := [...]string{
		"passwordReset",
		"onboarding",
		"passwordSet",
	}

	if d < PasswordReset || d > Onboarding {
		return "Unknown"
	}
	return strings[d]
}

func (d SuccessType) EnumIndex() int {
	return int(d)
}
