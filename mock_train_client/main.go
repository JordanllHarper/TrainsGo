package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JordanllHarper/trainsgo/shared"
)

type state int

const addr = "http://localhost:8080"

const (
	moving state = iota
	ended
)

type (
	trainInfo struct {
		key,
		ref string
	}

	pos struct {
		posX, posY int
	}
)

func parsePos(x, y string) (pos, error) {
	xInt, err := strconv.Atoi(x)
	if err != nil {
		return pos{}, errors.Join(errors.New("invalid x"), err)
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		return pos{}, errors.Join(errors.New("invalid y"), err)
	}
	return pos{xInt, yInt}, nil
}

const errorMessage = `Invalid number of arguments - need:
startX startY
endX endY
ref and key`
const requiredNumArgs = 6

func main() {

	updateTime := flag.Int("update", 250, "Milliseconds before an update")
	movementInc := flag.Int("inc", 1, "how many coordinates to move every update")
	flag.Parse()

	args := flag.Args()
	if len(args) != requiredNumArgs {
		log.Fatalln(errorMessage)
	}
	startX, startY := args[0], args[1]
	endX, endY := args[2], args[3]
	key, ref := args[4], args[5]
	startPos, err := parsePos(startX, startY)
	if err != nil {
		log.Fatalln("Error parsing start position:", err)
	}
	endPos, err := parsePos(endX, endY)
	if err != nil {
		log.Fatalln("Error parsing end position:", err)
	}

	tr := trainInfo{key: key, ref: ref}

	dur := time.Millisecond * (time.Duration(*updateTime))

	posCh := make(chan pos)
	go mockTravel(startPos, endPos, *movementInc, dur, posCh)
	for p := range posCh {
		log.Printf("Received pos - X: %v Y: %v", p.posX, p.posY)
		if err := sendHttpUpdate(p, addr, tr); err != nil {
			log.Fatalln("Error sending http update:", err)
		}
	}
	log.Println("Finished")
}

func getUpdateBody(p pos) (*bytes.Reader, error) {
	req := shared.PatchTrainRequest{
		PosX: &p.posX, PosY: &p.posY,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return &bytes.Reader{}, err
	}
	return bytes.NewReader(b), nil
}

func mockTravel(
	startPos, endPos pos,
	inc int,
	updDur time.Duration,
	posCh chan pos,
) {
	currentPos := startPos
	currentState := moving

	for {
		posCh <- currentPos
		switch currentState {
		case moving:
			currentPos = move(currentPos, endPos, inc)
			if currentPos == endPos {
				currentState = ended
			}
			time.Sleep(updDur)
		case ended:
			close(posCh)
			return
		default:
			panic(fmt.Sprintf("unexpected main.state: %#v", currentState))
		}
	}
}

func move(currentPos, endPos pos, inc int) pos {
	move := func(src, dest int) int {
		if src < dest && dest-src > inc {
			src += inc
		} else if src > dest && src-dest > inc {
			src -= inc
		} else {
			src = dest
		}
		return src
	}
	if currentPos.posX != endPos.posX {
		currentPos.posX = move(currentPos.posX, endPos.posX)
	}

	if currentPos.posY != endPos.posY {
		currentPos.posY = move(currentPos.posY, endPos.posY)
	}
	return currentPos
}

func sendHttpUpdate(p pos, addr string, t trainInfo) error {
	body, err := getUpdateBody(p)
	if err != nil {
		return fmt.Errorf("error getting update body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/trains/%s", addr, t.ref),
		body,
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", t.key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received %v status code", res.StatusCode)
	}
	var tr shared.Train
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}
