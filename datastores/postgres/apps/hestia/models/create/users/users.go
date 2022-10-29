package create_users

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/users"
)

type User struct {
	users.User
}

func NewCreateUser() User {
	return User{users.NewUser()}
}
