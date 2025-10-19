package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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

type GameRecord struct {
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	Date       time.Time `json:"date"`
}

type ScoreHistory struct {
	HighScore *GameRecord  `json:"high_score"`
	AllPlays  []GameRecord `json:"all_plays"`
}

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
	viewScoresBtn *widget.Button
	homeBtn       *widget.Button
	bottomBar     *fyne.Container
	scoreHistory  *ScoreHistory
	window        fyne.Window
	showSplash    func()
	mu            sync.Mutex
}

func getScoreHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "score-history.json"
	}
	dataDir := filepath.Join(homeDir, ".local", "share", "fyne-snake")
	os.MkdirAll(dataDir, 0755)
	return filepath.Join(dataDir, "score-history.json")
}

func loadScoreHistory() *ScoreHistory {
	path := getScoreHistoryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return &ScoreHistory{
			HighScore: &GameRecord{PlayerName: "Player", Score: 0, Date: time.Now()},
			AllPlays:  []GameRecord{},
		}
	}

	var history ScoreHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return &ScoreHistory{
			HighScore: &GameRecord{PlayerName: "Player", Score: 0, Date: time.Now()},
			AllPlays:  []GameRecord{},
		}
	}
	return &history
}

func saveScoreHistory(history *ScoreHistory) error {
	path := getScoreHistoryPath()
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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
		gameOverText.Move(fyne.NewPos(float32(canvasWidth/2-80), float32(canvasHeight/2-60)))
		g.canvas.Add(gameOverText)

		// Score text
		scoreText := canvas.NewText(fmt.Sprintf("Score: %d", g.score), color.White)
		scoreText.TextSize = 12
		scoreText.Move(fyne.NewPos(float32(canvasWidth/2-50), float32(canvasHeight/2-20)))
		g.canvas.Add(scoreText)

		// High Score text
		if g.scoreHistory != nil && g.scoreHistory.HighScore != nil {
			highScoreText := canvas.NewText(fmt.Sprintf("High Score: %d by %s", g.scoreHistory.HighScore.Score, g.scoreHistory.HighScore.PlayerName), colorTextGray)
			highScoreText.TextSize = 11
			highScoreText.Move(fyne.NewPos(float32(canvasWidth/2-90), float32(canvasHeight/2+5)))
			g.canvas.Add(highScoreText)
		}

		// Restart hint text
		restartHint := canvas.NewText("Press ENTER to restart", colorTextGray)
		restartHint.TextSize = 12
		restartHint.Move(fyne.NewPos(float32(canvasWidth/2-90), float32(canvasHeight/2+35)))
		g.canvas.Add(restartHint)

		// Show View Scores and Home buttons in bottom bar
		if g.viewScoresBtn != nil {
			g.viewScoresBtn.Show()
		}
		if g.homeBtn != nil {
			g.homeBtn.Show()
		}
	} else {
		// Hide View Scores and Home buttons during gameplay
		if g.viewScoresBtn != nil {
			g.viewScoresBtn.Hide()
		}
		if g.homeBtn != nil {
			g.homeBtn.Hide()
		}
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

				// Always prompt for name on game over
				fyne.Do(func() {
					g.promptForName()
				})

				fyne.Do(g.Draw)
				break
			}
			g.mu.Unlock()

			fyne.Do(g.Draw)
		}
	}()
}

func (g *Game) promptForName() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name")
	nameEntry.SetText("Player")

	// Check if it's a new high score
	isHighScore := g.score > g.scoreHistory.HighScore.Score

	var title string
	var message string
	if isHighScore {
		title = "🎉 New High Score!"
		message = fmt.Sprintf("Congratulations! New high score: %d", g.score)
	} else {
		title = "Game Over"
		message = fmt.Sprintf("Your score: %d", g.score)
	}

	dialog.ShowCustomConfirm(title, "Save", "Skip",
		container.NewVBox(
			widget.NewLabel(message),
			nameEntry,
		),
		func(save bool) {
			playerName := "Player"
			if save && nameEntry.Text != "" {
				playerName = nameEntry.Text
			}
			g.saveGameRecord(playerName)
			g.Draw()
		},
		g.window,
	)
}

func (g *Game) saveGameRecord(playerName string) {
	record := GameRecord{
		PlayerName: playerName,
		Score:      g.score,
		Date:       time.Now(),
	}

	// Add to all plays
	g.scoreHistory.AllPlays = append(g.scoreHistory.AllPlays, record)

	// Update high score if necessary
	if g.score > g.scoreHistory.HighScore.Score {
		g.scoreHistory.HighScore = &record
	}

	saveScoreHistory(g.scoreHistory)
}

func (g *Game) showScoreHistory() {
	if len(g.scoreHistory.AllPlays) == 0 {
		dialog.ShowInformation("Score History", "No games played yet!", g.window)
		return
	}

	// Create a list of all plays sorted by date (newest first)
	plays := make([]GameRecord, len(g.scoreHistory.AllPlays))
	copy(plays, g.scoreHistory.AllPlays)

	// Sort by date descending
	for i := 0; i < len(plays)-1; i++ {
		for j := i + 1; j < len(plays); j++ {
			if plays[j].Date.After(plays[i].Date) {
				plays[i], plays[j] = plays[j], plays[i]
			}
		}
	}

	// Create content
	var content string
	content = fmt.Sprintf("High Score: %d by %s\n\n", g.scoreHistory.HighScore.Score, g.scoreHistory.HighScore.PlayerName)
	content += "Recent Games:\n"
	content += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

	// Show last 20 plays
	max := 20
	if len(plays) < max {
		max = len(plays)
	}

	for i := 0; i < max; i++ {
		play := plays[i]
		dateStr := play.Date.Format("2006-01-02 15:04")
		content += fmt.Sprintf("%s - %s: %d\n", dateStr, play.PlayerName, play.Score)
	}

	if len(plays) > max {
		content += fmt.Sprintf("\n... and %d more games", len(plays)-max)
	}

	// Create scrollable label
	label := widget.NewLabel(content)
	label.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(label)
	scroll.SetMinSize(fyne.NewSize(400, 500))

	dialog.ShowCustom("Score History", "Close", scroll, g.window)
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

func showSplashScreen(w fyne.Window, onStart func(), history *ScoreHistory, onViewScores func()) {
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

	highScoreLabel := widget.NewLabel(fmt.Sprintf("High Score: %d by %s", history.HighScore.Score, history.HighScore.PlayerName))
	highScoreLabel.Alignment = fyne.TextAlignCenter
	highScoreLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}

	rulesTitle := widget.NewLabel("GAME RULES:")
	rulesTitle.Alignment = fyne.TextAlignCenter
	rulesTitle.TextStyle = fyne.TextStyle{Bold: true}

	rules := widget.NewLabel("• Use Arrow Keys to move the snake\n• Eat food to grow and score\n• Don't hit the walls or yourself\n• Press SPACE to pause/resume\n• Speed increases every 5 points")
	rules.Alignment = fyne.TextAlignCenter

	footer := widget.NewLabel("Developed in GoLang")
	footer.Alignment = fyne.TextAlignCenter
	footer.TextStyle = fyne.TextStyle{Italic: true}

	footer2 := widget.NewLabel("powered by fyne GUI")
	footer2.Alignment = fyne.TextAlignCenter
	footer2.TextStyle = fyne.TextStyle{Italic: true}

	startBtn := widget.NewButton("START GAME", func() {
		onStart()
	})
	startBtn.Importance = widget.HighImportance

	scoresBtn := widget.NewButton("View Scores", func() {
		onViewScores()
	})

	quitBtn := widget.NewButton("QUIT", func() {
		w.Close()
	})

	buttons := container.NewHBox(
		startBtn,
		scoresBtn,
		quitBtn,
	)

	enterHint := widget.NewLabel("Press ENTER to start")
	enterHint.Alignment = fyne.TextAlignCenter
	enterHint.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(
		container.NewCenter(title),
		container.NewCenter(subtitle),
		container.NewCenter(highScoreLabel),
		widget.NewSeparator(),
		container.NewCenter(rulesTitle),
		container.NewCenter(rules),
		widget.NewSeparator(),
		container.NewCenter(footer),
		container.NewCenter(footer2),
		widget.NewSeparator(),
		container.NewCenter(buttons),
		container.NewCenter(enterHint),
	)

	scrollContent := container.NewVScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(500, 550))

	// Add keyboard handler for Enter key
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyReturn || k.Name == fyne.KeyEnter {
			onStart()
		}
	})

	w.SetContent(container.NewCenter(scrollContent))
	w.Resize(fyne.NewSize(520, 600))
}

func main() {
	a := app.New()
	w := a.NewWindow("Snake Game")

	scoreLabel := widget.NewLabel("Score: 0")
	gameOverLabel := widget.NewLabel("")
	restartBtn := widget.NewButton("Restart", nil)
	viewScoresBtn := widget.NewButton("View Scores", nil)
	homeBtn := widget.NewButton("Home", nil)

	gameArea := container.NewWithoutLayout()

	// Load score history
	scoreHistory := loadScoreHistory()

	g := &Game{
		canvas:        gameArea,
		scoreLabel:    scoreLabel,
		gameOverLabel: gameOverLabel,
		restartBtn:    restartBtn,
		viewScoresBtn: viewScoresBtn,
		homeBtn:       homeBtn,
		scoreHistory:  scoreHistory,
		window:        w,
	}

	restartBtn.OnTapped = g.Reset
	viewScoresBtn.OnTapped = func() {
		g.showScoreHistory()
	}
	homeBtn.OnTapped = func() {
		g.gameRunning = false
		if g.showSplash != nil {
			g.showSplash()
		}
	}

	// Initially hide the View Scores and Home buttons
	viewScoresBtn.Hide()
	homeBtn.Hide()

	// UI layout
	rightButtons := container.NewHBox(homeBtn, viewScoresBtn, restartBtn)
	bottomBar := container.NewBorder(
		nil,
		nil,
		scoreLabel,
		rightButtons,
		container.NewCenter(gameOverLabel),
	)
	g.bottomBar = bottomBar

	// Function to start the game
	startGame := func() {
		// Set keyboard handling for game
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
			case fyne.KeyReturn, fyne.KeyEnter:
				if !g.gameRunning {
					go g.Reset()
				}
			}
		})

		w.SetContent(container.NewBorder(nil, bottomBar, nil, nil, gameArea))
		w.Resize(fyne.NewSize(float32(canvasWidth+6), float32(canvasHeight+60)))
		w.SetFixedSize(true)
		g.Reset()
	}

	// Show splash screen function
	var showSplash func()
	showSplash = func() {
		showSplashScreen(w, startGame, scoreHistory, func() {
			g.showScoreHistory()
		})
	}
	g.showSplash = showSplash

	// Show splash screen first
	showSplash()

	w.ShowAndRun()
}
