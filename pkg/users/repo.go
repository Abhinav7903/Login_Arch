package users

import "LoginArch/factory"

type Repository interface {
	CreateUser(user factory.User) error
}
