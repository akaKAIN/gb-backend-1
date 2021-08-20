package main

import (
	"github.com/akaKAIN/gb-backend-1/lesson_01/chat/models"
	"github.com/akaKAIN/gb-backend-1/lesson_01/chat/server"
)

var chat *models.Chat

func init() {
	chat = models.InitChat()
}

func main() {
	server.StartChat(chat)
}
