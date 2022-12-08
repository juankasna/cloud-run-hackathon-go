package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)

	posX := 0
	posY := 0
	dir := ""
	url := input.Links.Self.Href
	players := input.Arena.State

	for player_url, player := range players {

		if player_url == url {
			posX = player.X
			posY = player.Y
			dir = player.Direction
			log.Printf("This is me %v, %v, %v", posX, posY, dir)
		}
	}

	commands := []string{"F", "R", "L", "T"}
	rand := rand2.Intn(4)

	// TODO add your implementation here to replace the random response
	return commands[rand]
}
