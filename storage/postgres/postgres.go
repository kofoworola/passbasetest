package postgres

import (
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DbHost       string `default:"localhost"`
	DbPassword   string `required:"true"`
	DbUsername   string `default:"postgres"`
	DbName       string `default:"postgres"`
	DBPort       string `default:"5432"`
	DbMigrations string `default:"migrations"`
}

type Storage struct {
	db  *sqlx.DB
	cfg *Config

	migrate *migrate.Migrate
}

func New(cfg *Config) (*Storage, error) {
	dbString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.DbHost,
		cfg.DBPort,
		cfg.DbUsername,
		cfg.DbName,
		cfg.DbPassword,
	)
	fmt.Printf("db string is %s\n", dbString)

	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.DbMigrations,
		cfg.DbName,
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db:      db,
		cfg:     cfg,
		migrate: m,
	}, nil
}

func (s *Storage) Migrate() error {

	if err := s.migrate.Up(); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
