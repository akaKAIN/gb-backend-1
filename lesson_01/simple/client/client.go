package client

import (
	"io"
	"log"
	"net"
	"os"
)

func Start() {
	conn, err := net.Dial("tcp", net.JoinHostPort("localhost", "9000"))
	if err != nil {
		log.Fatal(err)
	}

	// Выводим поток из conn в поток вывода (Stdout)
	if _, err = io.Copy(os.Stdout, conn); err != nil {
		log.Println("Copy conn->Stdout error:", err)
	}
}
