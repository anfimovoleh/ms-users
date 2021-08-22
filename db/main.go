package db

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *dbx.DB
}

func New(link string) (*DB, error) {
	db, err := dbx.Open("postgres", link)
	return &DB{db: db}, err
}

//go:generate go-bindata -nometadata -ignore .+\.go$ -pkg db -o bindata.go ./...
//go:generate gofmt -w bindata.go

const (
	MigrationsDir = "migrations"
)

var (
	Migrations *MigrationsLoader
)

type AssetFn func(name string) ([]byte, error)
type AssetDirFn func(name string) ([]string, error)
