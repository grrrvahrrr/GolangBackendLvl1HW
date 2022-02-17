package main

import (
	"context"
	"lesson6/lesson6/reguser/internal/infrastructure/api/handler"
	"lesson6/lesson6/reguser/internal/infrastructure/api/routerchi"
	"lesson6/lesson6/reguser/internal/infrastructure/db/files/userfilemanager"
	"lesson6/lesson6/reguser/internal/infrastructure/server"
	"lesson6/lesson6/reguser/internal/usecases/app/repos/userrepo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// ust := usermemstore.NewUsers()
	ust, err := userfilemanager.NewUsers("./data.log", "mem://userRefreshTopic")
	if err != nil {
		log.Fatal(err)
	}

	us := userrepo.NewUsers(ust)
	hs := handler.NewHandlers(us)
	// h := defmux.NewRouter(hs)
	h := routerchi.NewRouterChi(hs)
	srv := server.NewServer(":8000", h)

	srv.Start(us)
	log.Print("Start")

	<-ctx.Done()

	srv.Stop()
	cancel()
	ust.Close()

	log.Print("Exit")
}
