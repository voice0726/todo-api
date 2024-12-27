package infra

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	"github.com/voice0726/todo-app-api/config"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type DataBase struct {
	dialector gorm.Dialector
	config    *gorm.Config
	*gorm.DB
}

func (d *DataBase) Open() error {
	var err error
	d.DB, err = gorm.Open(d.dialector, d.config)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func NewDB(c *config.Config, lc fx.Lifecycle) (*DataBase, error) {
	psql := postgres.New(postgres.Config{DSN: c.DSN})

	gormConfig := &gorm.Config{
		Logger: zapgorm2.New(zap.L()),
	}

	db := &DataBase{dialector: psql, config: gormConfig}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.Open()
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}
