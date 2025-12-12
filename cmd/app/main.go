package main

import (
	"budget-tracking-app/internal/handler/rest"
	"budget-tracking-app/internal/repository"
	"budget-tracking-app/internal/service"
	"budget-tracking-app/pkg/bcrypt"
	"budget-tracking-app/pkg/config"
	"budget-tracking-app/pkg/database/mariadb"
	"budget-tracking-app/pkg/jwt"
	"budget-tracking-app/pkg/middleware"
	"log"
)

func main() {
	config.LoadEnvironment()

	db, err := mariadb.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	bcrypt := bcrypt.Init()
	jwtAuth := jwt.Init()
	svc := service.NewService(db, repo, bcrypt, jwtAuth)
	middleware := middleware.Init(svc, jwtAuth)
	r := rest.NewRest(svc, middleware)
	r.MountEndPoint()
	r.Run()
}
