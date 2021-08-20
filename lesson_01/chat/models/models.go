package models

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
)

type User *net.Conn

type Address struct {
	Host    string
	Port    string
	Network string
}

// ChanList Структура содержащая в себе все каналы отвечающие за сообщения в чате
type ChanList struct {
	entering chan string
	leaving  chan string
	messages chan string
}

// Enter Отправка сообщения о входе пользователя в чат
func (cl *ChanList) Enter(name string) {
	msg := fmt.Sprintf("%s was entering in chat.", name)
	cl.entering <- msg
}

// Leave Отправка сообщения о выходе пользователя из чата
func (cl *ChanList) Leave(name string) {
	msg := fmt.Sprintf("%s was leaving us.", name)
	cl.leaving <- msg
}

// Say Отправка сообщения в общий чат
func (cl *ChanList) Say(name, msg string) {
	message := fmt.Sprintf("%s say: %s was entering in chat.", name, msg)
	cl.messages <- message
}

// Закрытие каналов
func (cl *ChanList) closeChannels() {
	close(cl.leaving)
	close(cl.entering)
	close(cl.messages)
}

// InitChanList инициализация.
func InitChanList() *ChanList {
	return &ChanList{
		entering: make(chan string),
		leaving:  make(chan string),
		messages: make(chan string),
	}
}

// Chat Структура для отправки сообщений пользователей в общий чат
type Chat struct {
	sync.Mutex
	Address  Address
	Channels *ChanList
	Users    map[User]string
}

func (chat *Chat) AddUser(conn *net.Conn, userName string) {
	chat.createOrUpdateUser(conn, userName)
}

func (chat *Chat) createOrUpdateUser(conn *net.Conn, userName string) {
	chat.Lock()
	defer chat.Unlock()
	chat.Users[conn] = userName
}

func (chat *Chat) DeleteUser(conn *net.Conn) {
	chat.Lock()
	defer chat.Unlock()
	delete(chat.Users, conn)
}

// Рассылка сообщения по всем активным пользователям
func (chat *Chat) broadcast(msg string) {
	for user, _ := range chat.Users {
		if _, err := fmt.Fprint(*user, msg); err != nil {
			log.Println("broadcasting error:", err)
		}
	}
}

// Close Завершение процессов чата
func (chat *Chat) Close() {
	chat.broadcast("Attention! Server will shutdown.")
	chat.Channels.closeChannels()
}

// StartListening Прослушивание каналов и рассылка сообщений
func (chat *Chat) StartListening(ctx context.Context) {
	defer chat.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-chat.Channels.entering:
			chat.broadcast(message)
		case message := <-chat.Channels.leaving:
			chat.broadcast(message)
		case message := <-chat.Channels.messages:
			chat.broadcast(message)
		}
	}
}

// InitChat инициализация
func InitChat() *Chat {
	return &Chat{
		Mutex:    sync.Mutex{},
		Address:  Address{Host: "localhost", Port: "9000", Network: "tcp"},
		Channels: InitChanList(),
		Users:    make(map[User]string),
	}
}
