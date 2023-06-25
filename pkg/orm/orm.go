// Package orm is described reusable package for database orm
package orm

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Provider struct {
	*gorm.DB
}

type QueryBuilder struct {
	Search      string
	Page        int
	ItemPerPage int
}

type ConfigConnProvider struct {
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func NewMySQL(ctx context.Context, connString string, cfg *ConfigConnProvider, ormConfig *gorm.Config) (*Provider, error) {
	orm, err := gorm.Open(mysql.Open(connString), ormConfig)
	if err != nil {
		return nil, err
	}

	db, err := orm.WithContext(ctx).DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &Provider{orm}, nil
}
