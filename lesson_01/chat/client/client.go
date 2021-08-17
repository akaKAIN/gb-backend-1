package client

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

func StartClient() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Bytes()

	conn, err := net.Dial("tcp", net.JoinHostPort("localhost", "9000"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if _, err = conn.Write(name); err != nil {
		log.Println(err)
	}

	// Выводим поток из conn в поток вывода (Stdout)
	go func() {
		if _, err = io.Copy(os.Stdout, conn); err != nil {
			log.Println("Copy conn->Stdout error:", err)
		}
	}()

	if _, err = io.Copy(conn, os.Stdin); err != nil {
		log.Println("Copy Stdin => conn error:", err)
	}
	log.Println("sssss")

}
