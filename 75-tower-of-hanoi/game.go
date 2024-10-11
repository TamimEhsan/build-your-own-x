package main

import (
	"fmt"
	"os"
)

type point struct {
	x, y int
}
type game struct {
	diskCount     int
	base, offset  int
	height, width int
	board         [][]string
	disks         [3][]int
	score         int
}

func (g *game) init(diskCount int) {
	g.diskCount = diskCount
	g.base = diskCount + 1
	g.offset = diskCount*2 + 1
	g.height = diskCount + 1
	g.width = (2*diskCount+1)*3 + 2
	g.board = make([][]string, g.height+2)
	for i := 0; i < g.height+2; i++ {
		g.board[i] = make([]string, g.width+2)
	}

	g.disks[0] = make([]int, diskCount)
	for i := 0; i < diskCount; i++ {
		g.disks[0][i] = diskCount - i
	}
	g.clearBoard()
}

func (g *game) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (g *game) clearBoard() {
	for i := 0; i < g.height+2; i++ {
		for j := 0; j < g.width+2; j++ {
			g.board[i][j] = " "
		}
	}
}

func (g *game) moveDisk(from, to int) {
	if from == to || len(g.disks[from]) == 0 {
		return
	}
	disk := g.disks[from][len(g.disks[from])-1]

	if len(g.disks[to]) > 0 && g.disks[to][len(g.disks[to])-1] < disk {
		return
	}

	g.disks[from] = g.disks[from][:len(g.disks[from])-1]
	g.disks[to] = append(g.disks[to], disk)
	g.score++
}

func (g *game) isGameOver() bool {
	return len(g.disks[2]) == g.diskCount || len(g.disks[1]) == g.diskCount
}

func (g *game) end() {
	fmt.Println()
	fmt.Println("You did it in ", g.score, " moves")
	os.Exit(0)
}

func (g *game) renderBoard() {
	for i := 0; i < 3; i++ {
		g.board[g.height+1][g.base+g.offset*i] = string(i + 49)
		for j := 0; j < len(g.disks[i]); j++ {
			for k := 0; k < g.disks[i][j]; k++ {
				g.board[g.height-j][g.base+g.offset*i-k-1] = "X"
			}
			g.board[g.height-j][g.base+g.offset*i] = string(g.disks[i][j] + 48)
			for k := 0; k < g.disks[i][j]; k++ {
				g.board[g.height-j][g.base+g.offset*i+k+1] = "X"
			}
		}
		for j := len(g.disks[i]); j <= g.diskCount; j++ {
			g.board[g.height-j][g.base+g.offset*i] = "|"
		}

	}
}

func (g *game) printBoard() {
	fmt.Println()
	for i := 0; i < g.height+2; i++ {
		for j := 0; j < g.width+2; j++ {
			fmt.Print(g.board[i][j])
		}
		fmt.Println()
	}
	fmt.Println("Moves:", g.score)
	//os.Stdout.Sync()
}
