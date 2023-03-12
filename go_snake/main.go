package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

type Snake struct {
	body      []Coord
	bodyShape string
	head      string
}

type Coord struct {
	x int
	y int
}

type Food struct {
	pos    Coord
	symbol string
}

type Game struct {
	snake  Snake
	food   Food
	score  int
	speed  time.Duration
	width  int
	height int
}

func NewGame() *Game {
	g := &Game{
		snake: Snake{
			body:      []Coord{{10, 10}, {11, 10}, {12, 10}},
			bodyShape: "-",
			head:      "<",
		},
		food: Food{
			pos:    Coord{15, 10},
			symbol: " ",
		},
		score:  0,
		speed:  700,
		width:  45,
		height: 20,
	}
	return g
}

func (g *Game) Run() {

	rand.Seed(time.Now().UnixNano())

	dir := "left"
	prevdir := "left"

	var currchar rune
	var currkey keyboard.Key
	var err error

	go func() {
		for {
			currchar, currkey, err = keyboard.GetKey()
			if err != nil {
				panic(err)
			}
		}
	}()

	for {

		if g.score > 2 {
			g.speed = 600
		} else if g.score > 5 {
			g.speed = 555
		} else if g.score > 10 {
			g.speed = 400
		} else if g.score > 15 {
			g.speed = 300
		} else if g.score > 20 {
			g.speed = 200
		} else if g.score > 30 {
			g.speed = 100
		}

		time.Sleep(g.speed * time.Millisecond)

		if currkey == keyboard.KeyEsc || currchar == 'q' || currchar == 'Q' {
			os.Exit(0)
		}

		if currchar == 'w' || currchar == 'W' || currkey == keyboard.KeyArrowUp {
			dir = "up"
			g.snake.head = "^"
			g.snake.bodyShape = "|"
		} else if currchar == 'a' || currchar == 'A' || currkey == keyboard.KeyArrowLeft {
			dir = "left"
			g.snake.head = "<"
			g.snake.bodyShape = "-"
		} else if currchar == 's' || currchar == 'S' || currkey == keyboard.KeyArrowDown {
			dir = "down"
			g.snake.head = "v"
			g.snake.bodyShape = "|"
		} else if currchar == 'd' || currchar == 'D' || currkey == keyboard.KeyArrowRight {
			dir = "right"
			g.snake.head = ">"
			g.snake.bodyShape = "-"
		}

		if currchar == 0 && currkey == 0 {
			dir = prevdir
		}

		prevdir = dir

		g.Update(dir)
		g.Draw()
	}
}

func (g *Game) Update(dir string) {

	head := g.snake.body[0]

	var newHead Coord

	switch dir {
	case "up":
		newHead = Coord{head.x, head.y - 1}
	case "down":
		newHead = Coord{head.x, head.y + 1}
	case "left":
		newHead = Coord{head.x - 1, head.y}
	case "right":
		newHead = Coord{head.x + 1, head.y}
	}

	if newHead.x < 0 || newHead.x >= g.width || newHead.y < 0 || newHead.y >= g.height {
		fmt.Println("You hit the wall!")
		fmt.Printf("Your score: %d\n", g.score)
		os.Exit(0)
	}

	g.snake.body = append([]Coord{newHead}, g.snake.body...)

	if g.food.pos == newHead {
		g.score++
		g.speed -= 5
		g.GenerateFood()
	} else {
		g.snake.body = g.snake.body[:len(g.snake.body)-1]
	}

	if g.IsGameOver() {
		fmt.Println("Game Over!")
		fmt.Printf("Your score: %d\n", g.score)
		os.Exit(0)
	}
}

func (g *Game) BorderContains(p Coord) bool {

	if p.x == 0 || p.y == 0 {
		return true
	}

	return true
}

func (g *Game) Draw() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Printf("Score: %d\n", g.score)

	headX := g.snake.body[0].x
	headY := g.snake.body[0].y

	if headX < 0 {
		headX = g.width - 1
	} else if headX >= g.width {
		headX = 0
	}
	if headY < 0 {
		headY = g.height - 1
	} else if headY >= g.height {
		headY = 0
	}

	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			p := Coord{j, i}

			switch {
			case p == Coord{headX, headY}:
				fmt.Print(g.snake.head)
			case g.snake.Contains(p):
				if (p.x+p.y)%2 == 0 {
					g.snake.bodyShape = "\\"
				} else {
					g.snake.bodyShape = "/"
				}
				fmt.Print(g.snake.bodyShape)
			case g.food.pos == p:
				fmt.Print(g.food.symbol)
			case g.BorderContains(p):
				fmt.Print("#")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Game) GenerateFood() {
	var x, y int
	for {
		x = rand.Intn(g.width)
		y = rand.Intn(g.height)

		if !g.snake.Contains(Coord{x, y}) {
			break
		}
	}
	g.food.pos = Coord{x, y}
}

func (g *Game) IsGameOver() bool {
	head := g.snake.body[0]
	for _, v := range g.snake.body[1:] {
		if v == head {
			return true
		}
	}
	return false
}

func (s Snake) Contains(c Coord) bool {
	for _, v := range s.body {
		if v == c {
			return true
		}
	}
	return false
}

func main() {

	log.SetFlags(3 | 16)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	intro := `  
	Go________              __           
	/   _____/ ____ _____  |  | __ ____  
	\_____  \ /    \\__  \ |  |/ // __ \ 
	/  ___   \  ||  \/ __ \|    <\  ___/ 
       /_______  /__||  (____  /__|_ \\___  >
               \/     \/     \/     \/    \/ 
	     	           by SAPPHIRE_KNIGHT
							  
	Press ENTER to START a New Game 
	Press Q to QUIT the Game 
	
	


                      Developed in GoLang`

	startFlag := false
	MaxLoader := 5

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	go func() {
		for {

			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}

			if key == keyboard.KeyEnter {
				startFlag = true
				return
			} else if char == 'q' || char == 'Q' {
				println()
				os.Exit(0)
			}
		}
	}()

	for !startFlag {
		fmt.Print(intro)
		for i := 0; i < MaxLoader; i++ {
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
			if startFlag {
				break
			}
		}
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	g := NewGame()
	g.Run()
}
