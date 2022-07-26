package repositorys

import (
	"context"
	"golang_rest_websockets/models"
)

type Repository interface {
	InsertUser(context context.Context, user *models.User) error
	GetUserById(context context.Context, id string) (*models.User, error)
	GetUserByEmail(context context.Context, email string) (*models.User, error)
	InsertPost(context context.Context, post *models.Post) error
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func InsertUser(context context.Context, user *models.User) error {
	return implementation.InsertUser(context, user)
}

func GetUserById(context context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(context, id)
}

func GetUserByEmail(context context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(context, email)
}

func Close() error {
	return implementation.Close()
}

func InsertPost(context context.Context, post *models.Post) error {
	return implementation.InsertPost(context, post)
}
