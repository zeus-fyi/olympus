package users

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

type User struct {
	autogen_bases.Users
}

type UserGroup struct {
	Slice autogen_bases.UsersSlice
}

func NewUser() User {
	u := User{autogen_bases.Users{
		UserID:   0,
		Metadata: "{}",
	}}

	return u
}
