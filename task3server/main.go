package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

var nickBase = []string{""}
var gameStarted bool
var gameResult int

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	lis, err := cfg.Listen(ctx, "tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	log.Println("I'm started!")
	go broadcaster()

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			wg.Add(1)
			go handleConn(conn, wg)
		}
	}()

	<-ctx.Done()
	log.Println("Done.")

	lis.Close()
	wg.Wait()
	log.Println("Exit.")
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	var nickName string
	var hasName bool

	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	ch <- "Please, enter your name and press 'Enter'"
	reader := bufio.NewReader(conn)
	nickName, _ = reader.ReadString('\n')
	nickName = strings.TrimSpace(nickName)

	for !hasName {
		for i, v := range nickBase {
			if nickName == v {
				ch <- "Nickname already exists. Please, enter another name and press 'Enter'"
				nickName, _ = reader.ReadString('\n')
				nickName = strings.TrimSpace(nickName)
				break
			} else if i == len(nickBase)-1 {
				nickBase = append(nickBase, nickName)
				hasName = true
			}
		}
	}

	if !gameStarted {
		go func() {
			for {
				if !gameStarted {
					gameResult = startGame()
					gameStarted = true
				}
			}
		}()
	}

	ch <- "You are " + nickName
	messages <- nickName + " has arrived"
	entering <- ch

	log.Println(nickName + " (" + who + ")" + " has arrived")

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- nickName + ": " + input.Text()
		answer, _ := strconv.Atoi(input.Text())
		if answer == gameResult {
			messages <- nickName + " is a Winner!"
			gameStarted = false
		} else {
			messages <- nickName + " is a Looser!"
		}
	}

	leaving <- ch
	messages <- nickName + " has left"

	for i, v := range nickBase {
		if nickName == v {
			nickBase[i] = nickBase[len(nickBase)-1]
			nickBase[len(nickBase)-1] = ""
			nickBase = nickBase[:len(nickBase)-1]
		}
	}

	gameStarted = false
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func startGame() int {
	expr, result := genRndExpression()
	messages <- "New expression: " + expr
	return result
}

func genRndNum() int {
	min := 0
	max := 10
	return rand.Intn(max-min+1) + min
}

func genRndSymbol() string {
	symbols := []string{"+", "-", "/", "*"}
	return symbols[rand.Intn(len(symbols))]
}

func genRndExpression() (string, int) {
	num1 := genRndNum()
	num2 := genRndNum()
	symbol := genRndSymbol()

	switch symbol {
	case "+":
		expr := strconv.Itoa(num1) + " + " + strconv.Itoa(num2) + " ="
		result := num1 + num2
		return expr, result
	case "-":
		expr := strconv.Itoa(num1) + " - " + strconv.Itoa(num2) + " ="
		result := num1 - num2
		return expr, result
	case "*":
		expr := strconv.Itoa(num1) + " * " + strconv.Itoa(num2) + " ="
		result := num1 * num2
		return expr, result
	case "/":
		for num2 == 0 {
			num2 = genRndNum()
		}
		expr := strconv.Itoa(num1) + " / " + strconv.Itoa(num2) + " ="
		result := num1 / num2
		return expr, result
	}

	return "", 0
}
