package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/juliotorresmoreno/neosmarthpen/config"
)

func NewConn(config config.Config) (*xorm.Engine, error) {
	driver := config.Database.Driver
	host := config.Database.Host
	user := config.Database.User
	pass := config.Database.Pass
	port := config.Database.Port
	name := config.Database.Name
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, name)
	conn, err := xorm.NewEngine(driver, dsn)

	return conn, err
}
