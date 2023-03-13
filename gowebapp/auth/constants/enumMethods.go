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

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (d AuthType) EnumIndex() int {
	return int(d)
}

// Enum EnvConst
type EnvConst int

// String - Creating common behavior - give the type a String function
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
	return strings[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (d EnvConst) EnumIndex() int {
	return int(d)
}

// Enum EmailType
type EmailType int

// String - Creating common behavior - give the type a String function
func (d EmailType) String() string {
	strings := [...]string{
		"newMagicLink",
		"successMessage",
	}

	if d < NewMagicLink || d > SuccessMessage {
		return "Unknown"
	}
	return strings[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (d EmailType) EnumIndex() int {
	return int(d)
}

// Enum SuccessType starting with index 1
type SuccessType int

// String - Creating common behavior - give the type a String function
func (d SuccessType) String() string {
	strings := [...]string{
		"passwordReset",
		"onboarding",
		"passwordSet",
	}

	if d < PasswordReset || d > Onboarding {
		return "Unknown"
	}
	return strings[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (d SuccessType) EnumIndex() int {
	return int(d)
}
