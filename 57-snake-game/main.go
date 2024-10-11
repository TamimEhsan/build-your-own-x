package main

import (
	"os"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

// Shared buffer for produced items
type Buffer struct {
	data  string
	mutex sync.Mutex
}

// Reader goroutine: consumes data every 1 second
func runLoop(buffer *Buffer) {
	game := game{}
	game.init(12, 24)

	for {
		// Lock the buffer and check if there is data to consume

		buffer.mutex.Lock()
		dir := buffer.data
		buffer.data = ""
		buffer.mutex.Unlock()
		if dir != "" {
			game.setDirection(dir)
		}
		game.clearBoard()
		game.moveSnake()
		if game.isGameOver() {
			game.end()
		}
		game.renderBoard()
		game.clearScreen()
		game.printBoard()

		time.Sleep(300 * time.Millisecond)
	}
}

// Writer goroutine: reads user input and produces data
func inputReader(buffer *Buffer) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// Lock the buffer and produce data
		buffer.mutex.Lock()
		if key == keyboard.KeyEsc {
			os.Exit(0)
		} else if key == keyboard.KeyArrowUp {
			buffer.data = "UP"
		} else if key == keyboard.KeyArrowDown {
			buffer.data = "DOWN"
		} else if key == keyboard.KeyArrowLeft {
			buffer.data = "LEFT"
		} else if key == keyboard.KeyArrowRight {
			buffer.data = "RIGHT"
		} else {
			//buffer.data = string(char)
		}
		buffer.mutex.Unlock()
	}
}

func main() {
	// Create a shared buffer between the reader and writer
	buffer := &Buffer{}
	// Start reader and writer goroutines

	go inputReader(buffer)
	runLoop(buffer)

}

/*
▀▁▂▃▄▅▆▇█▉▊▋▌▍▎▏
▐░▒▓▔▕▖ ▗▘▙▚▛▜▝▞▟
*/
