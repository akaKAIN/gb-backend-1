package models

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	msgCh   chan string
	clients map[*net.Conn]bool
	errCh   chan error
}

// HandleError Пишет переданную строку и ошибку в канал для ошибок
func (s *Server) HandleError(prefixError string, err error) {
	newErr := fmt.Sprintln(prefixError, err)
	s.errCh <- errors.New(newErr)
}

// Say Отправка переданной строки всем пользователям
func (s *Server) Say(msg string) {
	message := "server say: " + msg
	for client := range s.clients {
		if _, err := fmt.Fprintln(*client, message); err != nil {
			s.HandleError("Say err:", err)
		}
	}
}

func (s *Server) AddClient(conn *net.Conn) {
	s.clients[conn] = true
	log.Printf("client: %s, was added", (*conn).LocalAddr())
}

// InitServer Инициализация сущности Server
func InitServer() *Server {
	return &Server{
		clients: make(map[*net.Conn]bool),
		errCh:   make(chan error),
		msgCh:   make(chan string),
	}
}
