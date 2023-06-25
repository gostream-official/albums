package main

import (
	"fmt"

	"github.com/gostream-official/albums/impl/funcs/createalbum"
	"github.com/gostream-official/albums/impl/funcs/deletealbum"
	"github.com/gostream-official/albums/impl/funcs/getalbum"
	"github.com/gostream-official/albums/impl/funcs/getalbums"
	"github.com/gostream-official/albums/impl/funcs/updatealbum"
	"github.com/gostream-official/albums/impl/inject"
	"github.com/gostream-official/albums/pkg/env"
	"github.com/gostream-official/albums/pkg/router"
	"github.com/gostream-official/albums/pkg/store"

	"github.com/revx-official/output/log"
)

// Description:
//
//	The package initializer function.
//	Initializes the log level to info.
func init() {
	log.Level = log.LevelInfo
}

// Description:
//
//	The main function.
//	Represents the entry point of the application.
func main() {
	log.Infof("booting service instance ...")

	mongoUsername, err := env.GetEnvironmentVariable("MONGO_USERNAME")
	if err != nil {
		log.Fatalf("Cannot retrieve mongo username via environment variable")
	}

	mongoPassword, err := env.GetEnvironmentVariable("MONGO_PASSWORD")
	if err != nil {
		log.Fatalf("Cannot retrieve mongo password via environment variable")
	}

	mongoHost := env.GetEnvironmentVariableWithFallback("MONGO_HOST", "127.0.0.1:27017")

	connectionURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUsername, mongoPassword, mongoHost)
	instance, err := store.NewMongoInstance(connectionURI)

	log.Infof("establishing database connection ...")
	if err != nil {
		log.Fatalf("failed to connect to mongo instance: %s", err)
	}

	log.Infof("successfully established database connection")

	injector := inject.Injector{
		MongoInstance: instance,
	}

	log.Infof("launching router engine ...")
	engine := router.Default()

	engine.HandleWith("GET", "/albums", getalbums.Handler).Inject(injector)
	engine.HandleWith("GET", "/albums/:id", getalbum.Handler).Inject(injector)
	engine.HandleWith("POST", "/albums", createalbum.Handler).Inject(injector)
	engine.HandleWith("PUT", "/albums/:id", updatealbum.Handler).Inject(injector)
	engine.HandleWith("DELETE", "/albums/:id", deletealbum.Handler).Inject(injector)

	err = engine.Run(9871)
	if err != nil {
		log.Fatalf("failed to launch router engine: %s", err)
	}
}
