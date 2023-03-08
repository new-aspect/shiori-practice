package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/new-aspect/shiori-practice/internal/model"
	"github.com/pkg/errors"
)

type SQLiteDatabase struct {
	dbbase
}

// OpenSQLiteDatabase creates and open connection to new SQLite3 database.
func OpenSQLiteDatabase(ctx context.Context, databasePath string) (sqliteDB *SQLiteDatabase, err error) {
	// open database
	db, err := sqlx.ConnectContext(ctx, "sqlite", databasePath)
	if err != nil {
		return nil, err
	}
	sqliteDB = &SQLiteDatabase{
		dbbase: dbbase{*db},
	}
	return sqliteDB, nil
}

// Migrate runs migrations for this database engine
func (db *SQLiteDatabase) Migrate() error {
	sourceDrive, err := iofs.New(migrations, "migrations/sqlite")
	if err != nil {
		return errors.WithStack(err)
	}

	dbDrive, err := sqlite.WithInstance(db.DB.DB, &sqlite.Config{})
	if err != nil {
		return errors.WithStack(err)
	}

	migration, err := migrate.NewWithInstance(
		"iofs",
		sourceDrive,
		"sqlite",
		dbDrive,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.WithStack(err)
	}

	return nil
}

func (*SQLiteDatabase) SaveBookmarks(ctx context.Context, create bool, bookmarks ...model.Bookmark) ([]model.Bookmark, error) {
	return nil, nil
}

func (*SQLiteDatabase) GetBookMarks(ctx context.Context, opts GetBookmarksOptions) ([]model.Bookmark, error) {
	return nil, nil
}
