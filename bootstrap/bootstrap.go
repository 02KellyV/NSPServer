package bootstrap

import (
	"github.com/juliotorresmoreno/neosmarthpen/config"
	"github.com/juliotorresmoreno/neosmarthpen/db"
	"github.com/juliotorresmoreno/neosmarthpen/models"
)

func Sync(config config.Config) error {
	conn, err := db.NewConn(config)
	if err != nil {
		return err
	}
	defer conn.Close()
	models := []interface{}{
		models.User{},
	}
	for _, model := range models {
		if err = conn.Sync2(model); err != nil {
			return err
		}
	}
	return nil
}
