package users

import "LoginArch/factory"

type Repository interface {
	CreateUser(user factory.User) error
	Login(data factory.User) (bool, error)
	GetUser(email string) (factory.User, error)
}
