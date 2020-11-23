package generate

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

func initDb() error {
	cfg := conf.GetConfig()
	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Addr, cfg.Database)
	var err error
	gdb, err = gorm.Open(cfg.DriverName, url)
	if err != nil {
		panic(err)
		return nil
	}
	gdb.SingularTable(true)
	return nil
}
