package database

import (
	"context"
	"database/sql"
	"golang_rest_websockets/models"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPpostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InserUser(context context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(context, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(context context.Context, id int64) (*models.User, error) {
	rows, err := repo.db.QueryContext(context, "SELECT id, email FROM users WHERE id = $1", id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
