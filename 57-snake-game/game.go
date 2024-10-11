package main

import (
	"fmt"
	"math/rand"
	"os"
)

type point struct {
	x, y int
}
type game struct {
	height int
	width  int
	board  [][]string
	snake  []point
	food   point
	dir    map[string]point
	curDir point
	score  int
}

func (g *game) init(height, width int) {
	g.height = height
	g.width = width
	g.board = make([][]string, g.height+2)
	for i := range g.board {
		g.board[i] = make([]string, g.width+2)
	}
	g.snake = []point{{5, 5}, {5, 4}}
	g.dir = map[string]point{
		"UP":    {-1, 0},
		"DOWN":  {1, 0},
		"LEFT":  {0, -1},
		"RIGHT": {0, 1},
	}
	g.curDir = g.dir["RIGHT"]
	g.spawnFood()

	for i := 0; i < g.height+2; i++ {
		g.board[i][0] = "██"
		g.board[i][g.width+1] = "██"
	}

	for i := 0; i < g.width+2; i++ {
		g.board[0][i] = "██"
		g.board[g.height+1][i] = "██"
	}

}

func (g *game) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (g *game) clearBoard() {
	for i := 1; i <= g.height; i++ {
		for j := 1; j <= g.width; j++ {
			g.board[i][j] = "  "
		}
	}
}

func (g *game) findNextHead() point {
	nextHead := point{g.snake[0].x + g.curDir.x, g.snake[0].y + g.curDir.y}
	nextHead.x = (nextHead.x + g.height) % g.height
	nextHead.y = (nextHead.y + g.width) % g.width
	return nextHead
}

func (g *game) moveSnake() {
	nextHead := g.findNextHead()
	g.snake = append([]point{nextHead}, g.snake...)
	if g.isFood(nextHead) {
		g.score++
		g.spawnFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

}

func (g *game) setDirection(dir string) {
	prevDir := g.curDir
	g.curDir = g.dir[dir]
	if g.findNextHead() == g.snake[1] {
		g.curDir = prevDir
	}
}

func (g *game) isGameOver() bool {
	for i := 1; i < len(g.snake); i++ {
		if g.snake[0] == g.snake[i] {
			return true
		}
	}
	// if g.isWall(g.snake[0]) {
	// 	return true
	// }

	return false
}

func (g *game) end() {
	fmt.Println()
	fmt.Println("Game Over!")
	os.Exit(0)
}

func (g *game) IsSnake(x, y int) bool {
	for i := 0; i < len(g.snake); i++ {
		if g.snake[i].x == x && g.snake[i].y == y {
			return true
		}
	}
	return false
}

func (g *game) isWall(p point) bool {
	if p.x < 0 || p.x >= g.height || p.y < 0 || p.y >= g.width {
		return true
	}
	return false
}

func (g *game) isFood(p point) bool {
	if p == g.food {
		return true
	}
	return false
}

func (g *game) spawnFood() {
	for {
		x := rand.Intn(g.height)
		y := rand.Intn(g.width)
		if !g.IsSnake(x, y) {
			g.food = point{x, y}
			break
		}
	}
}

func (g *game) renderBoard() {

	for i := 0; i < len(g.snake); i++ {
		g.board[g.snake[i].x+1][g.snake[i].y+1] = "░░"
	}

	g.board[g.food.x+1][g.food.y+1] = "▓▓"
}

func (g *game) printBoard() {
	fmt.Println()
	for i := 0; i < g.height+2; i++ {
		for j := 0; j < g.width+2; j++ {
			fmt.Print(g.board[i][j])
		}
		fmt.Println()
	}
	fmt.Print("Score:", g.score)
	//os.Stdout.Sync()
}
