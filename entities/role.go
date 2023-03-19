package entities

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleBot    Role = "bot"
)
