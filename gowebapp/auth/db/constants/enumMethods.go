package constants

// Enum DBUserRole
type DBUserRole uint8

func (d DBUserRole) String() string {
	strings := [...]string{
		"users",
		"admins",
	}

	if d < Users || d > Admins {
		return "Unknown"
	}
	return strings[d]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (d DBUserRole) EnumIndex() int {
	return int(d)
}

// Enum DBAuthType
type DBAuthType uint8

func (d DBAuthType) String() string {
	strings := [...]string{
		"custom",
		"firebase",
	}

	if d < Custom || d > Firebase {
		return "Unknown"
	}
	return strings[d-1]
}

func (d DBAuthType) EnumIndex() int {
	return int(d)
}

// Enum DBProviderType
type DBProviderType uint8

func (d DBProviderType) String() string {
	strings := [...]string{
		"magiclink",
		"password",
		"google",
		"github",
	}

	if d < DBPassword || d > DBGithub {
		return "Unknown"
	}
	return strings[d-1]
}

func (d DBProviderType) EnumIndex() int {
	return int(d)
}

// Enum DBMagicLinkType
type DBMagicLinkType uint8

func (m DBMagicLinkType) String() string {
	strings := [...]string{
		"register",
		"login",
		"resetPassword",
	}

	if m < Register || m > ResetPassword {
		return "Unknown"
	}
	return strings[m-1]
}

func (m DBMagicLinkType) EnumIndex() int {
	return int(m)
}

// Enum DBPasswordAlgo
type DBPasswordAlgo uint8

func (m DBPasswordAlgo) String() string {
	strings := [...]string{
		"bcrypt",
	}

	if m < Bcrypt || m > Bcrypt {
		return "Unknown"
	}
	return strings[m-1]
}

func (m DBPasswordAlgo) EnumIndex() int {
	return int(m)
}
