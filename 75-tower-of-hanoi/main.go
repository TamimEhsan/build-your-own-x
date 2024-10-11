package main

import (
	"fmt"
	"sync"

	"github.com/eiannone/keyboard"
)

// Shared buffer for produced items
type Buffer struct {
	data  string
	mutex sync.Mutex
}

// Reader goroutine: consumes data every 1 second
func runLoop(msgChan <-chan int) {
	game := game{}
	game.init(5)

	game.clearBoard()
	game.renderBoard()
	game.clearScreen()
	game.printBoard()

	for {
		// Lock the buffer and check if there is data to consume
		from := <-msgChan
		fmt.Print("move disk from ", from)
		to := <-msgChan
		game.moveDisk(from-1, to-1)

		game.clearBoard()
		game.renderBoard()
		game.clearScreen()
		game.printBoard()

		fmt.Println("moved disk from ", from, " to ", to)

		if game.isGameOver() {
			game.end()
		}

	}
}

// Writer goroutine: reads user input and produces data
func inputReader(msgChan chan int) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if char == '1' {
			msgChan <- 1
		} else if char == '2' {
			msgChan <- 2
		} else if char == '3' {
			msgChan <- 3
		}

	}
}

func main() {
	// Create a shared buffer between the reader and writer
	msgChan := make(chan int, 2)
	// Start reader and writer goroutines

	go inputReader(msgChan)
	runLoop(msgChan)

}

/*
▀▁▂▃▄▅▆▇█▉▊▋▌▍▎▏
▐░▒▓▔▕▖ ▗▘▙▚▛▜▝▞▟
*/
