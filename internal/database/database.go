package database

import (
	"context"
	"embed"
	"github.com/jmoiron/sqlx"
	"github.com/new-aspect/shiori-practice/internal/model"
)

// OrderMethod is the order method for getting bookmarks
type OrderMethod int

type GetBookmarksOptions struct {
	IDs          []int
	Tags         []string
	ExcludedTags []string
	Keyword      string
	WithContent  bool
	OrderMethod  OrderMethod
	Limit        int
	Offset       int
}

// DB is interface for accessing and manipulating data in database .
type DB interface {
	// Migrate runs migrations for this database
	Migrate() error

	// SaveBookmarks saves bookmarks data to database
	SaveBookmarks(ctx context.Context, create bool, bookmarks ...model.Bookmark) ([]model.Bookmark, error)

	GetBookMarks(ctx context.Context, opts GetBookmarksOptions) ([]model.Bookmark, error)
}

type dbbase struct {
	sqlx.DB
}

//go:embed migrations/*
var migrations embed.FS
