package database

import (
	"context"
	"database/sql"
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

// GetAccount fetch account with matching username.
// Returns the account and boolean whether it's exist or not.
func (db *SQLiteDatabase) GetAccount(ctx context.Context, username string) (model.Account, bool, error) {
	account := model.Account{}
	err := db.GetContext(ctx, &account, `SELECT 
    	id, username, owner FROM account where username = ?`,
		username)
	if err != nil {
		//errors.WithStack(err) 是 Go 语言 errors 包中的一个函数，它的作用是将原始错误（err）包装为一个新的错误，该新错误包含了堆栈跟踪信息。
		return account, false, errors.WithStack(err)
	}

	return account, account.ID != 0, nil
}

// GetAccounts fetch list of account (without its password) based on submitted options
func (db *SQLiteDatabase) GetAccounts(ctx context.Context, opts GetAccountsOptions) ([]model.Account, error) {
	// Create query
	args := []interface{}{}
	query := `SELECT id, username, owner FROM  account where 1`

	if opts.Keyword != "" {
		query += "AND username LIKE ? "
		args = append(args, "%"+opts.Keyword+"%")
	}

	if opts.Owner {
		query += " AND owner = 1"
	}
	query += " ORDER BY username"

	// Fetch list account
	var accounts []model.Account
	err := db.SelectContext(ctx, &accounts, query, args...)
	if err != nil && err != sql.ErrNoRows {
		// WithStack在WithStack被调用时用堆栈跟踪来注释err。
		// 如果err是nil，WithStack返回nil。
		return nil, errors.WithStack(err)
	}

	return accounts, nil
}
