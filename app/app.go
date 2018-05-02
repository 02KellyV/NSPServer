package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/juliotorresmoreno/neosmarthpen/bootstrap"
	"github.com/juliotorresmoreno/neosmarthpen/config"
	"github.com/juliotorresmoreno/neosmarthpen/router"
)

type App struct {
	http.Handler
	config.Config
}

func NewApp() App {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	router := router.NewRouter(conf)
	err = bootstrap.Sync(conf)
	if err != nil {
		log.Fatal(err)
	}
	return App{
		Config:  conf,
		Handler: router,
	}
}

func (that App) Start() {
	c := make(chan int)
	go that.listenHTTP()
	<-c
}

func (that App) listenHTTP() {
	fmt.Println("Listen an serve on :4000")
	http.ListenAndServe(":4000", that)
}
