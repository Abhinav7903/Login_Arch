package postgres

import (
	"LoginArch/factory"
	"LoginArch/pkg/users"
	"fmt"
)

// CreateUser inserts a new user into the database
func (p *Postgres) CreateUser(user factory.User) error {
	// Start a transaction
	tx, err := p.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Ensure rollback is called if the function exits before commit
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute the insert query
	_, err = tx.Exec("INSERT INTO users (email, username, password) VALUES ($1, $2, $3)", user.Email, user.Name, user.HashedPassword)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Login checks if the user exists in the database
func (p *Postgres) Login(data factory.User) (factory.User, error) {
	var user factory.User
	var hashedPassword string

	row := p.dbConn.QueryRow("SELECT email, username, password FROM users WHERE email = $1", data.Email)
	err := row.Scan(&user.Email, &user.Name, &hashedPassword)
	if err != nil {
		return factory.User{}, fmt.Errorf("user not found: %w", err)
	}

	if !users.VerifyPassword(hashedPassword, data.Password) {
		return factory.User{}, fmt.Errorf("invalid password")
	}

	return user, nil
}

// GetUser retrieves a user by email from the database
func (p *Postgres) GetUser(email string) (factory.User, error) {
	// Start a transaction
	tx, err := p.dbConn.Begin()
	if err != nil {
		return factory.User{}, fmt.Errorf("error starting transaction: %w", err)
	}

	// Ensure rollback is called if the function exits before commit
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the query
	row := tx.QueryRow("SELECT email, username,created FROM users WHERE email = $1", email)

	// Scan the result into a User struct
	var user factory.User
	err = row.Scan(&user.Email, &user.Name, &user.Created)
	if err != nil {
		return factory.User{}, fmt.Errorf("error retrieving user: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return factory.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return user, nil
}
