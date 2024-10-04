package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/upper/db/v4"
	mysqladp "github.com/upper/db/v4/adapter/mysql"
	postgresqladp "github.com/upper/db/v4/adapter/postgresql"
)

func main() {
	var (
		dbtype string
		dsn    string
	)

	command := &cobra.Command{
		Use:   "db-data-generator",
		Short: "CLI to generate fake/test data and insert it into the database",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			session, err := createSession(dbtype, dsn)
			if err != nil {
				return err
			}
			rows, err := session.SQL().Query("SELECT 2")
			for rows.Next() {
				var name string
				err = rows.Scan(&name)
				if err != nil {
					break
				}
				fmt.Printf("Result: %s\n", name)
			}
			fmt.Printf("Rows: %v, Error: %v\n", rows, err)
			return
		},
	}
	command.Flags().StringVarP(&dbtype, "type", "t", "postgresql", "Database type (mysql or postgresql)")
	command.Flags().StringVarP(&dsn, "dsn", "d", "postgres://postgres@localhost:5432/postgres", "DSN connection string")
	command.Execute()
}

func createSession(dbtype, dsn string) (db.Session, error) {
	if dbtype == "postgresql" {
		url, err := postgresqladp.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
		return postgresqladp.Open(url)
	} else {
		url, err := mysqladp.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
		return mysqladp.Open(url)
	}
}
