package generate

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/modeltool/conf"
	"sync"
)

var gdb *gorm.DB
var store *Store
var storeOnce sync.Once

type Store struct {
	db *gorm.DB
}

func SharedStore() *Store {
	storeOnce.Do(func() {
		err := initDb()
		if err != nil {
			panic(err)
		}
		store = NewStore(gdb)
	})
	return store
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func initDb() *gorm.DB {
	cfg := conf.GetConfig()
	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		cfg.DataSource.User, cfg.DataSource.Password, cfg.DataSource.Addr, cfg.DataSource.Database)
	var err error
	gdb, err = gorm.Open(cfg.DataSource.DriverName, url)
	if err != nil {
		panic(err)
		return nil
	}
	gdb.SingularTable(true)
	return gdb
}
