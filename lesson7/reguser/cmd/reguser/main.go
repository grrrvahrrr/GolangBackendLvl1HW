package main

import (
	"context"
	"lesson7/lesson7/reguser/internal/infrastructure/api/handler"
	"lesson7/lesson7/reguser/internal/infrastructure/api/routeropenapi"
	"lesson7/lesson7/reguser/internal/infrastructure/db/pgstore"
	"lesson7/lesson7/reguser/internal/infrastructure/server"
	"lesson7/lesson7/reguser/internal/usecases/app/repos/userrepo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// ust := usermemstore.NewUsers()
	// ust, err := userfilemanager.NewUsers("./data.json", "mem://userRefreshTopic")
	ust, err := pgstore.NewUsers("postgres://deus:123@localhost:5432/deus?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	us := userrepo.NewUsers(ust)
	hs := handler.NewHandlers(us)
	// h := defmux.NewRouter(hs)
	// h := routerchi.NewRouterChi(hs)
	h := routeropenapi.NewRouterOpenAPI(hs)
	srv := server.NewServer(":8000", h)

	srv.Start(us)
	log.Print("Start")

	<-ctx.Done()

	srv.Stop()
	cancel()
	ust.Close()

	log.Print("Exit")
}
