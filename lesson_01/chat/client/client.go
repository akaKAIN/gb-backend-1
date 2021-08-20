package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func StartClient() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	conn, err := net.Dial("tcp", net.JoinHostPort("localhost", "9000"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if _, err = fmt.Fprint(conn, name); err != nil {
		log.Println(err)
	}

	for scanner.Scan() {
		if _, err = fmt.Fprint(conn, scanner.Text()); err != nil {
			log.Println(err)
		}
	}
}
