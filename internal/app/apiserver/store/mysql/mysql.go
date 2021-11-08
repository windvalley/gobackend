package mysql

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"gobackend/pkg/db"
	"gobackend/pkg/errors"
	"gobackend/pkg/log"

	"gobackend/internal/app/apiserver/store"
	v1 "gobackend/internal/pkg/entity/apiserver/v1"
	"gobackend/internal/pkg/gormlog"
	genericoptions "gobackend/internal/pkg/options"
)

type datastore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMysqlFactory ...
func GetMysqlFactory() store.Factory {
	return store.Client()
}

// InitMySQLFactory create mysql factory with the given config.
func InitMySQLFactory(opts *genericoptions.MySQLOptions) error {
	var err error
	var dbIns *gorm.DB

	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifetime: opts.MaxConnectionLifetime,
			LogLevel:              opts.LogLevel,
			Logger:                gormlog.New(opts.LogLevel),
		}

		dbIns, err = db.New(options)
		if err != nil {
			return
		}

		log.Info("start auto migrate mysql database ...")
		// not suggested in production environment.
		err = migrateDatabase(dbIns)
		if err != nil {
			err = fmt.Errorf("migrate database failed: %w", err)

			return
		}

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return fmt.Errorf(
			"failed to init mysql factory, mysqlFactory: %+v, error: %w",
			mysqlFactory,
			err,
		)
	}

	store.SetClient(mysqlFactory)

	log.Infof("init mysql factory instance: %+v", mysqlFactory)

	return nil
}

// cleanDatabase tear downs the database tables.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func cleanDatabase(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
		return errors.Wrap(err, "drop user table failed")
	}

	return nil
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&v1.User{},
	)
}

// resetDatabase resets the database tables.
// nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
func resetDatabase(db *gorm.DB) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}

	return nil
}
