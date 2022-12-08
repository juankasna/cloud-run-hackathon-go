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
	var xs []int
	var ys []int
	var scores []int
	posX := 0
	posY := 0
	// score := 0
	dir := ""
	url := input.Links.Self.Href
	state := input.Arena.State
	dimensions := input.Arena.Dimensions
	println("The dimensions are: %v", dimensions)
	// Check current status and location
	for player_url, player := range state {

		if player_url == url {
			posX = player.X
			posY = player.Y
			// score = player.Score
			dir = player.Direction
			println("My position is", posX, " ", posY, "and I'm looking", dir)
		} else {
			xs = append(xs, player.X)
			ys = append(ys, player.Y)
			scores = append(scores, player.Score)
		}
	}
	// See if can to shoot
	for i := 0; i < len(xs); i++ {
		difX := xs[i] - posX
		difY := ys[i] - posY
		if dir == "N" {
			if (difX == 0) && ((difY >= -3) && (difY < 0)) {
				println("Shoot, difX", difX, "difY", difY, "xs[i]", xs[i], "xy[i]", ys[i])
				return "T"
			}
		} else if dir == "W" {
			if (difY == 0) && ((difX >= -3) && (difX < 0)) {
				println("Shoot, difX", difX, "difY", difY, "xs[i]", xs[i], "xy[i]", ys[i])
				return "T"
			}
		} else if dir == "S" {
			if (difX == 0) && ((difY <= 3) && (difY > 0)) {
				println("Shoot, difX", difX, "difY", difY, "xs[i]", xs[i], "xy[i]", ys[i])
				return "T"
			}
		} else if dir == "E" {
			if (difY == 0) && ((difX <= 3) && (difX > 0)) {
				println("Shoot, difX", difX, "difY", difY, "xs[i]", xs[i], "xy[i]", ys[i])
				return "T"
			}
		}
	}

	commands := []string{"F", "R", "L", "T"}
	rand := rand2.Intn(4)
	returning := commands[rand]
	println("Rand: ", returning)
	// TODO add your implementation here to replace the random response
	return returning
}
