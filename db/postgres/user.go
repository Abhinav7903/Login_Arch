package postgres

import (
	"LoginArch/factory"
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
