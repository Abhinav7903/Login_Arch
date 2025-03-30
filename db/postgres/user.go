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
