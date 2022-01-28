package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}

	conn, err := d.DialContext(ctx, "tcp", "[::1]:9000")
	if err != nil {
		log.Fatal(err)
	}

	go func(conn net.Conn) {
		if <-ctx.Done() == struct{}{} {
			log.Printf("context canceled\nPress 'Enter' to exit.")
			conn.Close()
		}
	}(conn)

	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin)
	fmt.Printf("%s: exit\n", conn.LocalAddr())
}
