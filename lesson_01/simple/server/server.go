package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akaKAIN/gb-backend-1/lesson_01/simple/models"
)

func Start() {
	sl := startSignalListener()
	server := models.InitServer()
	go InitServerInputScan(server)

	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", "9000"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started ...")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept connection error: %s\n", err)
			}
			server.AddClient(&conn)

			go func() {
				ticker := time.NewTicker(time.Second)
				for {
					select {
					case <-ticker.C:
						if _, err := fmt.Fprintf(conn, "Ticker step ...\n"); err != nil {
							log.Println(err)
						}
					case <-sl:
						break
					}
				}
			}()

		}
	}()

	<-sl
	close(sl)
	byeMsg := "Server was shutdown."
	log.Println(byeMsg)
	server.Say(byeMsg)
	time.Sleep(time.Second)
}

func startSignalListener() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func InitServerInputScan(server *models.Server) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		server.Say(scanner.Text())
	}
}

func isClosedChan(ch chan os.Signal) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
