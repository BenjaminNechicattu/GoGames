package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed icon.png
var iconData []byte

const (
	playerX = "X"
	playerO = "O"
)

var winLines = [8][3][2]int{
	{{0, 0}, {0, 1}, {0, 2}}, // rows
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	{{0, 0}, {1, 0}, {2, 0}}, // columns
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	{{0, 0}, {1, 1}, {2, 2}}, // diagonals
	{{0, 2}, {1, 1}, {2, 0}},
}

func main() {
	a := app.New()

	// Set app icon
	icon := fyne.NewStaticResource("icon.png", iconData)
	a.SetIcon(icon)

	w := a.NewWindow("XOX Game")
	w.SetIcon(icon)

	currentPlayer := playerX
	vsComputer := true
	buttons := [3][3]*widget.Button{}
	status := widget.NewLabel("Your turn (X)")

	getBoardState := func() [3][3]string {
		var board [3][3]string
		for i := range buttons {
			for j := range buttons[i] {
				board[i][j] = buttons[i][j].Text
			}
		}
		return board
	}

	checkWinnerForBoard := func(board [3][3]string) string {
		for _, line := range winLines {
			v1 := board[line[0][0]][line[0][1]]
			if v1 != "" && v1 == board[line[1][0]][line[1][1]] && v1 == board[line[2][0]][line[2][1]] {
				return v1
			}
		}
		return ""
	}

	isBoardFull := func(board [3][3]string) bool {
		for i := range board {
			for j := range board[i] {
				if board[i][j] == "" {
					return false
				}
			}
		}
		return true
	}

	var minimax func(board [3][3]string, depth int, isMaximizing bool, alpha, beta int) int
	minimax = func(board [3][3]string, depth int, isMaximizing bool, alpha, beta int) int {
		winner := checkWinnerForBoard(board)
		if winner == playerO {
			return 10 - depth
		}
		if winner == playerX {
			return depth - 10
		}
		if isBoardFull(board) {
			return 0
		}

		if isMaximizing {
			bestScore := -1000
			for i := range board {
				for j := range board[i] {
					if board[i][j] == "" {
						board[i][j] = playerO
						score := minimax(board, depth+1, false, alpha, beta)
						board[i][j] = ""
						if score > bestScore {
							bestScore = score
						}
						if bestScore > alpha {
							alpha = bestScore
						}
						if beta <= alpha {
							return bestScore
						}
					}
				}
			}
			return bestScore
		}
		bestScore := 1000
		for i := range board {
			for j := range board[i] {
				if board[i][j] == "" {
					board[i][j] = playerX
					score := minimax(board, depth+1, true, alpha, beta)
					board[i][j] = ""
					if score < bestScore {
						bestScore = score
					}
					if bestScore < beta {
						beta = bestScore
					}
					if beta <= alpha {
						return bestScore
					}
				}
			}
		}
		return bestScore
	}

	resetBoard := func() {
		currentPlayer = playerX
		statusText := "Player X's turn"
		if vsComputer {
			statusText = "Your turn (X)"
		}
		status.SetText(statusText)
		for i := range buttons {
			for j := range buttons[i] {
				buttons[i][j].SetText("")
				buttons[i][j].Enable()
			}
		}
	}

	disableAllButtons := func() {
		for i := range buttons {
			for j := range buttons[i] {
				buttons[i][j].Disable()
			}
		}
	}

	var checkWinner func()
	checkWinner = func() {
		for _, line := range winLines {
			v1 := buttons[line[0][0]][line[0][1]].Text
			if v1 != "" && v1 == buttons[line[1][0]][line[1][1]].Text && v1 == buttons[line[2][0]][line[2][1]].Text {
				status.SetText("Winner: " + v1)
				disableAllButtons()
				return
			}
		}

		// Check draw
		for i := range buttons {
			for j := range buttons[i] {
				if buttons[i][j].Text == "" {
					return
				}
			}
		}
		status.SetText("It's a draw!")
		disableAllButtons()
	}

	var computerMove func()
	computerMove = func() {
		board := getBoardState()
		bestScore := -1000
		bestMoveI, bestMoveJ := -1, -1

		for i := range board {
			for j := range board[i] {
				if board[i][j] == "" {
					board[i][j] = playerO
					score := minimax(board, 0, false, -1000, 1000)
					board[i][j] = ""
					if score > bestScore {
						bestScore = score
						bestMoveI, bestMoveJ = i, j
					}
				}
			}
		}

		if bestMoveI != -1 {
			buttons[bestMoveI][bestMoveJ].SetText(playerO)
			checkWinner()
			if status.Text != "Winner: O" && status.Text != "It's a draw!" {
				currentPlayer = playerX
				status.SetText("Your turn (X)")
			}
		}
	}

	grid := container.NewGridWithColumns(3)
	for i := range buttons {
		for j := range buttons[i] {
			i, j := i, j
			buttons[i][j] = widget.NewButton("", func() {
				if buttons[i][j].Text != "" || (vsComputer && currentPlayer != playerX) {
					return
				}

				buttons[i][j].SetText(currentPlayer)
				checkWinner()

				// If game is over, don't switch players
				if status.Text == "Winner: "+currentPlayer || status.Text == "It's a draw!" {
					return
				}

				if vsComputer {
					currentPlayer = playerO
					status.SetText("Computer's turn...")
					computerMove()
				} else {
					if currentPlayer == playerX {
						currentPlayer = playerO
					} else {
						currentPlayer = playerX
					}
					status.SetText("Player " + currentPlayer + "'s turn")
				}
			})
			grid.Add(buttons[i][j])
		}
	}

	resetButton := widget.NewButton("Reset", func() {
		resetBoard()
	})

	modeSelect := widget.NewSelect([]string{"Player vs Player", "Player vs Computer"}, func(value string) {
		vsComputer = (value == "Player vs Computer")
		resetBoard()
	})
	modeSelect.SetSelected("Player vs Computer")

	content := container.NewVBox(
		modeSelect,
		status,
		grid,
		resetButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}
