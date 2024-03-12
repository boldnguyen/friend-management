package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/pkg/errors"
)

// ConnectDB establishes a connection to the database using the provided DB URL.
// It returns a database connection and any error encountered during the process.
func ConnectDB(dbURL string) (*sql.DB, error) {

	// Open a database connection using the "postgres" driver
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)

		// Wrap the error using errors.WithStack and add a descriptive error message
		return nil, errors.WithStack(fmt.Errorf("connect DB failed. Err: %w", err))
	}

	// Ping the database to ensure the connection is active
	if err = conn.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Initializing DB connection")

	return conn, nil
}
