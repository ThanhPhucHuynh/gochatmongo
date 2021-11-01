package main

import (
	"chat/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	wsServer := models.NewWebsocketServer()
	go wsServer.Run()

	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Da")
		models.ServeWs(wsServer, w, r)
	})

	// router.HandleFunc("/addrooms", models.NewRoom)
	router.HandleFunc("/rooms", models.RoomsList)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
	})

	fmt.Println("Server run 8080")
	http.ListenAndServe(":8080", c.Handler(router))
}
