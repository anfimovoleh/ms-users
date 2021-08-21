package db

import (
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
	migrationsDir = "migrations"
)

var (
	Migrations *MigrationsLoader
)

var log = logrus.New()

type AssetFn func(name string) ([]byte, error)
type AssetDirFn func(name string) ([]string, error)

func init() {
	Migrations = NewMigrationsLoader()
	if err := Migrations.loadDir(migrationsDir); err != nil {
		log.WithField("service", "load-migrations").WithError(err).Fatal("failed to load migrations")
		return
	}
}
