package infra

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	"github.com/voice0726/todo-app-api/config"
	"github.com/voice0726/todo-app-api/models"
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

func NewDB(c *config.Config, lg *zap.Logger, lc fx.Lifecycle) *DataBase {
	psql := postgres.New(postgres.Config{DSN: c.DSN})
	gormConfig := &gorm.Config{Logger: zapgorm2.New(lg)}
	db := &DataBase{dialector: psql, config: gormConfig}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := db.Open()
			if err != nil {
				lg.Error("failed to open database connection", zap.Error(err))
				return err
			}
			lg.Info("database connection established")
			if !c.IsProd {
				lg.Info("running in development mode")
				lg.Info("migrating database")
				err := db.AutoMigrate(models.Todo{}, models.Address{})
				if err != nil {
					lg.Error("failed to migrate", zap.Error(err))
					return err
				}
				// todo: add caller to logger
				lg.Info("migration completed")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			lg.Info("closing database connection")
			if err := db.Close(); err != nil {
				lg.Error("failed to close database connection", zap.Error(err))
				return err
			}
			lg.Info("database connection closed")
			return nil
		},
	})

	return db
}
