package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const boardSize = 8

type Move struct {
	fromRow, fromCol, toRow, toCol int
	capturedPiece                  string
}

type ChessGame struct {
	board       [boardSize][boardSize]string
	moveHistory []Move
	redoStack   []Move
}

func NewChessGame() *ChessGame {
	game := &ChessGame{}
	game.setupBoard()
	return game
}

func (c *ChessGame) setupBoard() {
	c.board = [boardSize][boardSize]string{
		{"R", "N", "B", "Q", "K", "B", "N", "R"},
		{"P", "P", "P", "P", "P", "P", "P", "P"},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"p", "p", "p", "p", "p", "p", "p", "p"},
		{"r", "n", "b", "q", "k", "b", "n", "r"},
	}
}

var boardColors = [2]string{"⬜", "⬛"} // White and black squares
var unicodePieces = map[string]string{
	"R": "♜", "N": "♞", "B": "♝", "Q": "♛", "K": "♚", "P": "♟",
	"r": "♖", "n": "♘", "b": "♗", "q": "♕", "k": "♔", "p": "♙",
}

func (c *ChessGame) printBoard() {
	fmt.Println("x   a   b   c   d   e   f   g   h  x \n")
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%d ", 8-i)
		for j := 0; j < boardSize; j++ {
			piece := c.board[i][j]
			color := boardColors[(i+j)%2]
			if piece == "" {
				fmt.Printf("%s  ", color)
			} else {
				fmt.Printf("%s%s ", color, unicodePieces[piece])
			}
		}
		fmt.Printf(" %d\n\n", 8-i)
	}
	fmt.Println("x   a   b   c   d   e   f   g   h  x")
}

func (c *ChessGame) movePiece(from, to string) bool {
	fromRow, fromCol := parsePosition(from)
	toRow, toCol := parsePosition(to)

	if fromRow == -1 || toRow == -1 {
		fmt.Println("Invalid move!")
		return false
	}

	if !c.isValidMove(fromRow, fromCol, toRow, toCol) {
		fmt.Println("Illegal move!")
		return false
	}

	move := Move{fromRow, fromCol, toRow, toCol, c.board[toRow][toCol]}
	c.moveHistory = append(c.moveHistory, move)
	c.redoStack = nil // Clear redo stack on new move

	c.board[toRow][toCol] = c.board[fromRow][fromCol]
	c.board[fromRow][fromCol] = ""
	return true
}

func (c *ChessGame) isValidMove(fromRow, fromCol, toRow, toCol int) bool {
	piece := c.board[fromRow][fromCol]
	if piece == "" {
		return false
	}
	if c.board[toRow][toCol] != "" && sameColor(piece, c.board[toRow][toCol]) {
		return false // Can't capture own piece
	}
	switch strings.ToUpper(piece) {
	case "P":
		return isValidPawnMove(fromRow, fromCol, toRow, toCol, piece, c.board)
	case "R":
		return isValidRookMove(fromRow, fromCol, toRow, toCol, c.board)
	case "N":
		return isValidKnightMove(fromRow, fromCol, toRow, toCol)
	case "B":
		return isValidBishopMove(fromRow, fromCol, toRow, toCol, c.board)
	case "Q":
		return isValidQueenMove(fromRow, fromCol, toRow, toCol, c.board)
	case "K":
		return isValidKingMove(fromRow, fromCol, toRow, toCol)
	default:
		return false
	}
}

func isValidKingMove(fromRow, fromCol, toRow, toCol int) bool {
	dr, dc := abs(fromRow-toRow), abs(fromCol-toCol)
	return dr <= 1 && dc <= 1
}

func isValidQueenMove(fromRow, fromCol, toRow, toCol int, board [boardSize][boardSize]string) bool {
	return isValidRookMove(fromRow, fromCol, toRow, toCol, board) || isValidBishopMove(fromRow, fromCol, toRow, toCol, board)
}

func isValidPawnMove(fromRow, fromCol, toRow, toCol int, piece string, board [boardSize][boardSize]string) bool {
	direction := -1
	if piece == "p" {
		direction = 1
	}
	if fromCol == toCol && board[toRow][toCol] == "" {
		return toRow == fromRow+direction || (fromRow == 1 && toRow == fromRow+2 && piece == "P") || (fromRow == 6 && toRow == fromRow-2 && piece == "p")
	}
	if abs(fromCol-toCol) == 1 && toRow == fromRow+direction && board[toRow][toCol] != "" {
		return true // Capturing diagonally
	}
	return false
}

func isValidKnightMove(fromRow, fromCol, toRow, toCol int) bool {
	dr, dc := abs(fromRow-toRow), abs(fromCol-toCol)
	return (dr == 2 && dc == 1) || (dr == 1 && dc == 2)
}

func isValidRookMove(fromRow, fromCol, toRow, toCol int, board [boardSize][boardSize]string) bool {
	if fromRow != toRow && fromCol != toCol {
		return false
	}
	return isPathClear(fromRow, fromCol, toRow, toCol, board)
}

func isValidBishopMove(fromRow, fromCol, toRow, toCol int, board [boardSize][boardSize]string) bool {
	if abs(fromRow-toRow) != abs(fromCol-toCol) {
		return false
	}
	return isPathClear(fromRow, fromCol, toRow, toCol, board)
}

func isPathClear(fromRow, fromCol, toRow, toCol int, board [boardSize][boardSize]string) bool {
	dr, dc := sign(toRow-fromRow), sign(toCol-fromCol)
	for r, c := fromRow+dr, fromCol+dc; r != toRow || c != toCol; r, c = r+dr, c+dc {
		if board[r][c] != "" {
			return false
		}
	}
	return true
}

func sign(n int) int {
	if n > 0 {
		return 1
	} else if n < 0 {
		return -1
	}
	return 0
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sameColor(a, b string) bool {
	return (a >= "A" && a <= "Z") == (b >= "A" && b <= "Z")
}

func parsePosition(pos string) (int, int) {
	if len(pos) != 2 {
		return -1, -1
	}
	col := int(pos[0] - 'a')
	row := 8 - int(pos[1]-'0')
	if col < 0 || col >= boardSize || row < 0 || row >= boardSize {
		return -1, -1
	}
	return row, col
}

func (c *ChessGame) undoMove() {
	if len(c.moveHistory) == 0 {
		fmt.Println("No moves to undo!")
		return
	}

	lastMove := c.moveHistory[len(c.moveHistory)-1]
	c.moveHistory = c.moveHistory[:len(c.moveHistory)-1]

	c.redoStack = append(c.redoStack, lastMove)
	c.board[lastMove.fromRow][lastMove.fromCol] = c.board[lastMove.toRow][lastMove.toCol]
	c.board[lastMove.toRow][lastMove.toCol] = lastMove.capturedPiece
}

func (c *ChessGame) redoMove() {
	if len(c.redoStack) == 0 {
		fmt.Println("No moves to redo!")
		return
	}

	move := c.redoStack[len(c.redoStack)-1]
	c.redoStack = c.redoStack[:len(c.redoStack)-1]
	c.moveHistory = append(c.moveHistory, move)

	c.board[move.toRow][move.toCol] = c.board[move.fromRow][move.fromCol]
	c.board[move.fromRow][move.fromCol] = ""
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	game := NewChessGame()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		clearTerminal()
		game.printBoard()
		fmt.Print("Enter move (e.g., e2-e4, or 'undo', or 'quit'): ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "quit" {
			break
		} else if strings.ToLower(input) == "undo" {
			game.undoMove()
		} else if strings.ToLower(input) == "redo" {
			game.redoMove()
		} else {
			moveParts := strings.Split(input, " ")
			if len(moveParts) != 2 {
				fmt.Println("Invalid move format. Use 'from to' (e.g., e2 e4).")
				continue
			}
			from, to := moveParts[0], moveParts[1]
			if !game.movePiece(from, to) {
				continue
			}
		}
	}
}
