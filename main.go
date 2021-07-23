package main

import (
	"chat/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := models.NewChanRoom()
	r.Run()
	sm := models.NewSaveMessageChan()

	router := mux.NewRouter()

	router.HandleFunc("/ws", models.ChannelChat(r, sm))

	router.HandleFunc("/addrooms", models.NewRoom)
	router.HandleFunc("/rooms", models.RoomsList)

	fmt.Println("Server run 8080")
	http.ListenAndServe(":8080", router)
}
