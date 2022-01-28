package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type client chan<- string

var (
	messages = make(chan string)
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	clients := make(map[client]bool)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	lis, err := cfg.Listen(ctx, "tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	log.Println("I'm started!")

	go inputScanner()
	go broabcaster(clients)

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Println(err)
			} else {
				wg.Add(1)
				go handleConn(ctx, conn, clients, wg)
			}
		}
	}()

	<-ctx.Done()

	log.Println("Done.")
	lis.Close()
	wg.Wait()
	log.Println("Exit.")
}

func handleConn(ctx context.Context, conn net.Conn, clients map[client]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	ch := make(chan string)
	clients[ch] = true
	go clientWriter(conn, ch)

	tck := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			delete(clients, ch)
			close(ch)
			return
		case t := <-tck.C:
			fmt.Fprintf(conn, "now:%s\n", t)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broabcaster(clients map[client]bool) {
	for msg := range messages {
		for cli := range clients {
			cli <- "Server message: " + msg
		}
	}
}

func inputScanner() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		messages <- input.Text()
	}
}
