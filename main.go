package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	// log.Printf("IN: %#v", input)
	var xs []int
	var ys []int
	var difxs []int
	var difys []int
	var difs []int
	var scores []int
	posX := 0
	posY := 0
	wasHit := false
	// score := 0
	dir := ""
	url := input.Links.Self.Href
	state := input.Arena.State
	dimensions := input.Arena.Dimensions
	fmt.Printf("The dimensions are: %v", dimensions)
	// Check current status and location
	for player_url, player := range state {

		if player_url == url {
			posX = player.X
			posY = player.Y
			// score = player.Score
			dir = player.Direction
			wasHit = player.WasHit
			println("My position is", posX, " ", posY, "and I'm looking", dir)
		} else {
			xs = append(xs, player.X)
			ys = append(ys, player.Y)
			scores = append(scores, player.Score)
		}
	}

	if wasHit {
		return canMove(posX, posY, xs, ys, dir)
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

		posabs := Abs(difY) + Abs(difX)
		difxs = append(difxs, difX)
		difys = append(difys, difY)
		difs = append(difs, posabs)
	}

	close := closest(difs)
	println("closest is", difxs[close], difys[close])

	if dir == "N" {
		if difys[close] < 0 {
			println("Forward, difX", difxs[close], "difY", difys[close])
			return "F"
		} else if difxs[close] < 0 {
			println("Left, difX", difxs[close], "difY", difys[close])
			return "L"
		} else if difxs[close] > 0 {
			println("Right, difX", difxs[close], "difY", difys[close])
			return "R"
		} else if difys[close] > 0 {
			println("Behind, difX", difxs[close], "difY", difys[close])
			return "R"
		}
	} else if dir == "W" {
		if difxs[close] < 0 {
			println("Forward, difX", difxs[close], "difY", difys[close])
			return "F"
		} else if difys[close] < 0 {
			println("Right, difX", difxs[close], "difY", difys[close])
			return "R"
		} else if difys[close] > 0 {
			println("Left, difX", difxs[close], "difY", difys[close])
			return "L"
		} else if difxs[close] > 0 {
			println("Behind, difX", difxs[close], "difY", difys[close])
			return "R"
		}
	} else if dir == "S" {
		if difxs[close] < 0 {
			println("Right, difX", difxs[close], "difY", difys[close])
			return "R"
		} else if difys[close] > 0 {
			println("Forward, difX", difxs[close], "difY", difys[close])
			return "F"
		} else if difxs[close] > 0 {
			println("Left, difX", difxs[close], "difY", difys[close])
			return "L"
		} else if difys[close] < 0 {
			println("Behind, difX", difxs[close], "difY", difys[close])
			return "R"
		}
	} else if dir == "E" {
		if difxs[close] > 0 {
			println("Forward, difX", difxs[close], "difY", difys[close])
			return "F"
		} else if difys[close] < 0 {
			println("Left, difX", difxs[close], "difY", difys[close])
			return "L"
		} else if difys[close] > 0 {
			println("Right, difX", difxs[close], "difY", difys[close])
			return "R"
		} else if difxs[close] < 0 {
			println("Behind, difX", difxs[close], "difY", difys[close])
			return "R"
		}
	}

	commands := []string{"F", "R", "L", "T"}
	rand := rand.Intn(4)
	returning := commands[rand]
	println("Rand: ", returning)
	// TODO add your implementation here to replace the random response
	return returning
}

func closest(list []int) int {
	min := list[0]
	response := 0
	for i := 1; i < len(list); i++ {
		if min > list[i] {
			min = list[i]
			response = i
		}
	}
	return response
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func canMove(posx int, posy int, xs []int, ys []int, direction string) string {
	up := true
	down := true
	right := true
	left := true

	for i := 0; i < len(xs); i++ {
		if xs[i] == posx && ys[i] == (posy-1) {
			up = false
		} else if xs[i] == (posx-1) && ys[i] == (posy) {
			left = false
		} else if xs[i] == (posx) && ys[i] == (posy+1) {
			right = false
		} else if xs[i] == (posx+1) && ys[i] == (posy) {
			down = false
		}
	}

	if direction == "N" {
		if up {
			return "F"
		} else if right {
			return "R"
		} else if left {
			return "L"
		} else if down {
			return "R"
		}
	} else if direction == "W" {
		if left {
			return "F"
		} else if up {
			return "R"
		} else if down {
			return "L"
		} else if right {
			return "R"
		}
	} else if direction == "S" {
		if down {
			return "F"
		} else if right {
			return "L"
		} else if left {
			return "R"
		} else if up {
			return "R"
		}
	} else if direction == "E" {
		if right {
			return "F"
		} else if up {
			return "L"
		} else if down {
			return "R"
		} else if left {
			return "R"
		}
	}
	return "F"
}
