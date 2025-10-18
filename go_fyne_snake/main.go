package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	cellSize     = 20
	gridW        = 20
	gridH        = 20
	canvasWidth  = cellSize * gridW
	canvasHeight = cellSize * gridH
	initialSpeed = 200
	minSpeed     = 50
	speedDelta   = 10
	eyeSize      = 12
	pupilSize    = 6
	noseSize     = 6
)

var (
	colorBg        = color.RGBA{40, 40, 40, 255}
	colorGrid      = color.RGBA{50, 50, 50, 255}
	colorHeadDark  = color.RGBA{0, 180, 0, 255}
	colorBodyDark  = color.RGBA{0, 180, 0, 255}
	colorBodyLight = color.RGBA{0, 220, 0, 255}
	colorNose      = color.RGBA{0, 220, 0, 255}
	colorEyeWhite  = color.RGBA{255, 255, 255, 255}
	colorEyePupil  = color.RGBA{0, 0, 0, 255}
	colorOverlay   = color.RGBA{0, 0, 0, 180}
	colorPauseOver = color.RGBA{0, 0, 0, 150}
	colorTextGray  = color.RGBA{200, 200, 200, 255}
)

type Game struct {
	snake         []Point
	dir           Direction
	food          Point
	foodColor     color.Color
	score         int
	speed         int
	gameRunning   bool
	paused        bool
	canvas        *fyne.Container
	scoreLabel    *widget.Label
	gameOverLabel *widget.Label
	restartBtn    *widget.Button
	mu            sync.Mutex
}

func (g *Game) Draw() {
	g.canvas.Objects = nil
	canvasSizeF := fyne.NewSize(float32(canvasWidth), float32(canvasHeight))

	// Draw background
	bg := canvas.NewRectangle(colorBg)
	bg.Resize(canvasSizeF)
	g.canvas.Add(bg)

	// Draw brick grid
	brickSize := fyne.NewSize(cellSize-1, cellSize-1)
	for y := 0; y < gridH; y++ {
		for x := 0; x < gridW; x++ {
			brick := canvas.NewRectangle(colorGrid)
			brick.Move(fyne.NewPos(float32(x*cellSize), float32(y*cellSize)))
			brick.Resize(brickSize)
			g.canvas.Add(brick)
		}
	}

	// Draw food as round circle with random color
	foodCircle := canvas.NewCircle(g.foodColor)
	foodCircle.Move(fyne.NewPos(float32(g.food.X*cellSize)+1, float32(g.food.Y*cellSize)+1))
	foodCircle.Resize(fyne.NewSize(cellSize-2, cellSize-2))
	g.canvas.Add(foodCircle)

	// Draw snake
	for i, s := range g.snake {
		if i == len(g.snake)-1 {
			g.drawSnakeHead(s)
		} else {
			g.drawSnakeBody(s, i)
		}
	}

	// Draw Paused overlay if game is paused
	if g.paused && g.gameRunning {
		overlay := canvas.NewRectangle(colorPauseOver)
		overlay.Resize(canvasSizeF)
		g.canvas.Add(overlay)

		// Paused text
		pausedText := canvas.NewText("PAUSED", color.White)
		pausedText.TextSize = 20
		pausedText.TextStyle = fyne.TextStyle{Bold: true}
		pausedText.Move(fyne.NewPos(float32(canvasWidth/2-70), float32(canvasHeight/2-20)))
		g.canvas.Add(pausedText)

		// Hint text
		hintText := canvas.NewText("Press SPACE to resume", colorTextGray)
		hintText.TextSize = 12
		hintText.Move(fyne.NewPos(float32(canvasWidth/2-90), float32(canvasHeight/2+25)))
		g.canvas.Add(hintText)
	}

	// Draw Game Over overlay if game is over
	if !g.gameRunning {
		overlay := canvas.NewRectangle(colorOverlay)
		overlay.Resize(canvasSizeF)
		g.canvas.Add(overlay)

		// Game Over text
		gameOverText := canvas.NewText("Game Over!", color.White)
		gameOverText.TextSize = 20
		gameOverText.TextStyle = fyne.TextStyle{Bold: true}
		gameOverText.Move(fyne.NewPos(float32(canvasWidth/2-80), float32(canvasHeight/2-40)))
		g.canvas.Add(gameOverText)

		// Score text
		scoreText := canvas.NewText(fmt.Sprintf("Score: %d", g.score), color.White)
		scoreText.TextSize = 12
		scoreText.Move(fyne.NewPos(float32(canvasWidth/2-50), float32(canvasHeight/2+10)))
		g.canvas.Add(scoreText)
	}

	g.canvas.Refresh()
}

func (g *Game) drawSnakeHead(s Point) {
	baseX := float32(s.X * cellSize)
	baseY := float32(s.Y * cellSize)
	center := float32(cellSize / 2)

	// Main head circle
	headCircle := canvas.NewCircle(colorHeadDark)
	headCircle.Move(fyne.NewPos(baseX+2, baseY+2))
	headCircle.Resize(fyne.NewSize(cellSize-4, cellSize-4))
	g.canvas.Add(headCircle)

	// Nose
	nose := canvas.NewCircle(colorNose)
	nose.Resize(fyne.NewSize(noseSize, noseSize))

	// Googly eyes
	eyeWhite1 := canvas.NewCircle(colorEyeWhite)
	eyeWhite2 := canvas.NewCircle(colorEyeWhite)
	pupil1 := canvas.NewCircle(colorEyePupil)
	pupil2 := canvas.NewCircle(colorEyePupil)

	// Position based on direction
	switch g.dir {
	case Up:
		nose.Move(fyne.NewPos(baseX+center-noseSize/2, baseY))
		eyeWhite1.Move(fyne.NewPos(baseX, baseY-2))
		eyeWhite2.Move(fyne.NewPos(baseX+cellSize-eyeSize, baseY-2))
		pupil1.Move(fyne.NewPos(baseX+3, baseY+1))
		pupil2.Move(fyne.NewPos(baseX+cellSize-eyeSize+3, baseY+1))
	case Down:
		nose.Move(fyne.NewPos(baseX+center-noseSize/2, baseY+cellSize-noseSize))
		eyeWhite1.Move(fyne.NewPos(baseX, baseY+cellSize-eyeSize-2))
		eyeWhite2.Move(fyne.NewPos(baseX+cellSize-eyeSize, baseY+cellSize-eyeSize-2))
		pupil1.Move(fyne.NewPos(baseX+3, baseY+cellSize-eyeSize+1))
		pupil2.Move(fyne.NewPos(baseX+cellSize-eyeSize+3, baseY+cellSize-eyeSize+1))
	case Left:
		nose.Move(fyne.NewPos(baseX, baseY+center-noseSize/2))
		eyeWhite1.Move(fyne.NewPos(baseX-2, baseY))
		eyeWhite2.Move(fyne.NewPos(baseX-2, baseY+cellSize-eyeSize))
		pupil1.Move(fyne.NewPos(baseX+1, baseY+3))
		pupil2.Move(fyne.NewPos(baseX+1, baseY+cellSize-eyeSize+3))
	case Right:
		nose.Move(fyne.NewPos(baseX+cellSize-noseSize, baseY+center-noseSize/2))
		eyeWhite1.Move(fyne.NewPos(baseX+cellSize-eyeSize-2, baseY))
		eyeWhite2.Move(fyne.NewPos(baseX+cellSize-eyeSize-2, baseY+cellSize-eyeSize))
		pupil1.Move(fyne.NewPos(baseX+cellSize-eyeSize+1, baseY+3))
		pupil2.Move(fyne.NewPos(baseX+cellSize-eyeSize+1, baseY+cellSize-eyeSize+3))
	}

	g.canvas.Add(nose)

	eyeWhite1.Resize(fyne.NewSize(eyeSize, eyeSize))
	eyeWhite2.Resize(fyne.NewSize(eyeSize, eyeSize))
	pupil1.Resize(fyne.NewSize(pupilSize, pupilSize))
	pupil2.Resize(fyne.NewSize(pupilSize, pupilSize))

	g.canvas.Add(eyeWhite1)
	g.canvas.Add(eyeWhite2)
	g.canvas.Add(pupil1)
	g.canvas.Add(pupil2)
}

func (g *Game) drawSnakeBody(s Point, index int) {
	var bodyColor color.Color
	if index%2 == 0 {
		bodyColor = colorBodyDark
	} else {
		bodyColor = colorBodyLight
	}
	body := canvas.NewCircle(bodyColor)
	body.Move(fyne.NewPos(float32(s.X*cellSize)+1, float32(s.Y*cellSize)+1))
	body.Resize(fyne.NewSize(cellSize-2, cellSize-2))
	g.canvas.Add(body)
}

func (g *Game) spawnFood() {
	for {
		g.food = Point{X: rand.Intn(gridW), Y: rand.Intn(gridH)}
		onSnake := false
		for _, s := range g.snake {
			if s == g.food {
				onSnake = true
				break
			}
		}
		if !onSnake {
			break
		}
	}
	// Generate random color for food
	g.foodColor = color.RGBA{
		uint8(rand.Intn(156) + 100), // R: 100-255
		uint8(rand.Intn(156) + 100), // G: 100-255
		uint8(rand.Intn(156) + 100), // B: 100-255
		255,
	}
}

func (g *Game) MoveSnake() bool {
	head := g.snake[len(g.snake)-1]
	newHead := head

	switch g.dir {
	case Up:
		newHead.Y--
	case Down:
		newHead.Y++
	case Left:
		newHead.X--
	case Right:
		newHead.X++
	}

	// Check collision with walls
	if newHead.X < 0 || newHead.X >= gridW || newHead.Y < 0 || newHead.Y >= gridH {
		return false
	}

	// Check collision with self
	for _, s := range g.snake {
		if s == newHead {
			return false
		}
	}

	g.snake = append(g.snake, newHead)

	// Check if eating food
	if newHead == g.food {
		g.score++
		fyne.Do(func() {
			g.scoreLabel.SetText(fmt.Sprintf("Score: %d", g.score))
		})
		g.spawnFood()
		// Increase speed every 5 points
		if g.score%5 == 0 && g.speed > minSpeed {
			g.speed -= speedDelta
		}
	} else {
		g.snake = g.snake[1:] // remove tail
	}

	return true
}

func (g *Game) Reset() {
	g.snake = []Point{{X: gridW / 2, Y: gridH / 2}}
	g.dir = Right
	g.score = 0
	g.speed = initialSpeed
	g.scoreLabel.SetText("Score: 0")
	g.gameOverLabel.SetText("")
	g.gameRunning = true
	g.paused = false
	g.spawnFood()

	// Start new game loop
	go func() {
		for g.gameRunning {
			time.Sleep(time.Duration(g.speed) * time.Millisecond)

			g.mu.Lock()
			if !g.gameRunning {
				g.mu.Unlock()
				break
			}

			if g.paused {
				g.mu.Unlock()
				continue
			}

			if !g.MoveSnake() {
				g.gameRunning = false
				g.mu.Unlock()
				fyne.Do(g.Draw)
				break
			}
			g.mu.Unlock()

			fyne.Do(g.Draw)
		}
	}()
}

type Point struct {
	X, Y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func showSplashScreen(w fyne.Window, onStart func()) {
	// ASCII Art Title
	title := widget.NewLabel(`    Go________              __           
    /   _____/ ____ _____  |  | __ ____  
    \_____  \ /    \\__  \ |  |/ // __ \ 
    /  ___   \  ||  \/ __ \|    <\  ___/ 
       /_______  /__||  (____  /__|_ \\___  >
               \/     \/     \/     \/    \/ `)
	title.TextStyle = fyne.TextStyle{Monospace: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("by SAPPHIRE_KNIGHT")
	subtitle.Alignment = fyne.TextAlignCenter
	subtitle.TextStyle = fyne.TextStyle{Bold: true}

	rulesTitle := widget.NewLabel("GAME RULES:")
	rulesTitle.Alignment = fyne.TextAlignCenter
	rulesTitle.TextStyle = fyne.TextStyle{Bold: true}

	rule1 := widget.NewLabel("• Use Arrow Keys to move the snake")
	rule1.Alignment = fyne.TextAlignCenter

	rule2 := widget.NewLabel("• Eat red food to grow and score")
	rule2.Alignment = fyne.TextAlignCenter

	rule3 := widget.NewLabel("• Don't hit the walls or yourself")
	rule3.Alignment = fyne.TextAlignCenter

	rule4 := widget.NewLabel("• Press SPACE to pause/resume")
	rule4.Alignment = fyne.TextAlignCenter

	rule5 := widget.NewLabel("• Speed increases every 5 points")
	rule5.Alignment = fyne.TextAlignCenter

	footer := widget.NewLabel("Developed in GoLang")
	footer.Alignment = fyne.TextAlignCenter
	footer.TextStyle = fyne.TextStyle{Italic: true}

	startBtn := widget.NewButton("START GAME", func() {
		onStart()
	})
	startBtn.Importance = widget.HighImportance

	quitBtn := widget.NewButton("QUIT", func() {
		w.Close()
	})

	buttons := container.NewHBox(
		startBtn,
		quitBtn,
	)

	content := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		widget.NewSeparator(),
		container.NewCenter(rulesTitle),
		container.NewCenter(rule1),
		container.NewCenter(rule2),
		container.NewCenter(rule3),
		container.NewCenter(rule4),
		container.NewCenter(rule5),
		widget.NewSeparator(),
		container.NewCenter(footer),
		widget.NewSeparator(),
		container.NewCenter(buttons),
	)

	scrollContent := container.NewVScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(500, 550))

	w.SetContent(container.NewCenter(scrollContent))
	w.Resize(fyne.NewSize(520, 600))
}

func main() {
	a := app.New()
	w := a.NewWindow("Snake Game")

	scoreLabel := widget.NewLabel("Score: 0")
	gameOverLabel := widget.NewLabel("")
	restartBtn := widget.NewButton("Restart", nil)

	gameArea := container.NewWithoutLayout()

	g := &Game{
		canvas:        gameArea,
		scoreLabel:    scoreLabel,
		gameOverLabel: gameOverLabel,
		restartBtn:    restartBtn,
	}

	restartBtn.OnTapped = g.Reset

	// Keyboard handling
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		g.mu.Lock()
		defer g.mu.Unlock()

		switch k.Name {
		case fyne.KeyUp:
			if g.dir != Down && g.gameRunning {
				g.dir = Up
			}
		case fyne.KeyDown:
			if g.dir != Up && g.gameRunning {
				g.dir = Down
			}
		case fyne.KeyLeft:
			if g.dir != Right && g.gameRunning {
				g.dir = Left
			}
		case fyne.KeyRight:
			if g.dir != Left && g.gameRunning {
				g.dir = Right
			}
		case fyne.KeySpace:
			if g.gameRunning {
				g.paused = !g.paused
				fyne.Do(g.Draw)
			}
		}
	})

	// UI layout
	bottomBar := container.NewBorder(
		nil,
		nil,
		scoreLabel,
		restartBtn,
		container.NewCenter(gameOverLabel),
	)

	// Function to start the game
	startGame := func() {
		w.SetContent(container.NewBorder(nil, bottomBar, nil, nil, gameArea))
		w.Resize(fyne.NewSize(float32(canvasWidth+6), float32(canvasHeight+60)))
		w.SetFixedSize(true)
		g.Reset()
	}

	// Show splash screen first
	showSplashScreen(w, startGame)

	w.ShowAndRun()
}
