package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/akaKAIN/gb-backend-1/lesson_01/chat/models"
)

func StartChat(chat *models.Chat) {
	ctx, cancel := context.WithCancel(context.Background())
	go chat.StartListening(ctx)

	sc := make(chan os.Signal)
	go shutdownListen(cancel, sc)

	listener, err := net.Listen(chat.Address.Network, net.JoinHostPort(chat.Address.Host, chat.Address.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server was start...")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
			}
			go HandleConnection(ctx, cancel, chat, &conn)
		}
	}()
	<-ctx.Done()
	log.Println("Server was stop")
}

func HandleConnection(ctx context.Context, cancel context.CancelFunc, chat *models.Chat, conn *net.Conn) {
	//userName := GetUserName(conn)
	buf := make([]byte, 256)
	if _, err := (*conn).Read(buf); err != nil {
		log.Println(err)
	}
	name := string(buf)
	log.Println("name is ->", name)
	chat.AddUser(conn, name)
	chat.Channels.Enter(name)
	defer chat.Channels.Leave(name)
	for {
		_, err := (*conn).Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
			}
		}
		if _, err := fmt.Fprintf(os.Stdout, "%s: %s\n", name, string(buf)); err != nil {
			log.Println(err)
			break
		}
	}

}

// GetUserName Получение имени пользователя
func GetUserName(conn *net.Conn) string {
	if _, err := fmt.Fprintf(*conn, "Введите свое имя: "); err != nil {
		log.Println(err)
	}
	buf := make([]byte, 4)
	for {
		n, err := (*conn).Read(buf)
		if err != nil {
			log.Println("Err:", err)
		}
		if n != 0 || err == io.EOF {
			log.Println("N", n, err)
			break
		}
	}

	log.Println("DONE", string(buf))
	return string(buf)
}

func shutdownListen(cancel context.CancelFunc, ch chan os.Signal) {
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}
