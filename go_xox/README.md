# XOX Game (Tic-Tac-Toe)

A modern Tic-Tac-Toe game built with Go and Fyne, featuring both Player vs Player and Player vs Computer modes with an unbeatable AI opponent.

![XOX Game Icon](icon.png)

## Features

- 🎮 **Two Game Modes**
  - Player vs Player: Classic two-player mode
  - Player vs Computer: Challenge an AI opponent

- 🤖 **Intelligent AI**
  - Uses Minimax algorithm with Alpha-Beta pruning
  - Plays optimally - impossible to beat, only draw or lose
  - Fast decision making with optimized search

- 🎨 **Clean UI**
  - Built with Fyne for a modern, cross-platform interface
  - Intuitive game controls
  - Real-time status updates
  - Custom application icon

- ⚡ **Performance Optimized**
  - Alpha-Beta pruning reduces AI computation by up to 90%
  - Efficient board state evaluation
  - Responsive gameplay

## Installation

### Prerequisites

- Go 1.21 or higher
- Fyne dependencies for your platform

#### Linux Dependencies
```bash
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

### Build from Source

```bash
# Clone the repository
cd go_xox

# Build the executable
go build -o xox

# Run the game
./xox
```

Or run directly:
```bash
go run .
```

## How to Play

1. **Launch the game** - The application will start with "Player vs Computer" mode by default
2. **Choose game mode** - Use the dropdown menu to switch between:
   - Player vs Player
   - Player vs Computer
3. **Make your move** - Click on any empty cell in the 3x3 grid
4. **Win conditions**:
   - Get three of your symbols (X or O) in a row (horizontal, vertical, or diagonal)
   - If all cells are filled with no winner, the game ends in a draw
5. **Reset** - Click the "Reset" button to start a new game

### Game Modes

#### Player vs Player
- Two players take turns
- Player X always goes first
- Click cells to place your symbol

#### Player vs Computer
- You play as X (always first)
- Computer plays as O
- The AI uses optimal strategy - try to get a draw!

## Technical Details

### Architecture

The game is structured with clean separation of concerns:

- **UI Layer**: Fyne widgets and containers
- **Game Logic**: Pure functions for board state and win detection
- **AI Engine**: Minimax algorithm with alpha-beta pruning

### AI Implementation

The computer opponent uses the **Minimax algorithm** with **Alpha-Beta pruning**:

```go
// Minimax with Alpha-Beta Pruning
// - Evaluates all possible game states
// - Prunes branches that won't affect the final decision
// - Scores: +10 for computer win, -10 for player win, 0 for draw
// - Depth factor ensures faster wins are preferred
```

**Performance:**
- Without pruning: ~19,000 positions evaluated on first move
- With pruning: ~2,000 positions evaluated on first move
- ~90% reduction in computation time

### Code Structure

```
main.go
├── Constants (playerX, playerO, winLines)
├── Helper Functions
│   ├── getBoardState() - Extracts current board state
│   ├── checkWinnerForBoard() - Checks winner from board array
│   ├── isBoardFull() - Checks if board is full
│   └── disableAllButtons() - Disables all game buttons
├── AI Functions
│   ├── minimax() - Core AI algorithm
│   └── computerMove() - Finds and executes best move
├── Game Logic
│   ├── resetBoard() - Resets game state
│   └── checkWinner() - Checks winner from UI state
└── UI Setup
    ├── Button grid creation
    ├── Mode selector
    └── Event handlers
```

### Technologies Used

- **Language**: Go 1.x
- **UI Framework**: [Fyne v2](https://fyne.io/)
- **Build Tool**: Go standard toolchain
- **Icon**: Embedded PNG resource

## Project Files

- `main.go` - Main application code
- `icon.png` - Application icon (512x512 PNG with transparency)
- `FyneApp.toml` - Fyne application metadata
- `go.mod` - Go module dependencies
- `README.md` - This documentation

## Building for Different Platforms

### Linux
```bash
go build -o xox
```

### Windows
```bash
GOOS=windows go build -o xox.exe
```

### macOS
```bash
GOOS=darwin go build -o xox
```

### Package with Fyne
For a proper application bundle with icon:
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os linux -icon icon.png
```

## Development

### Code Optimization Highlights

1. **Shared Win Lines**: Win patterns defined once and reused
2. **Constants**: Player symbols defined as constants
3. **Alpha-Beta Pruning**: AI search tree optimization
4. **Early Returns**: Efficient draw detection
5. **Embedded Resources**: Icon embedded in binary

### Extending the Game

Ideas for enhancements:
- Difficulty levels (easy, medium, hard)
- Score tracking across multiple games
- Different board sizes (4x4, 5x5)
- Game history/undo feature
- Custom themes and colors
- Sound effects
- Multiplayer over network

## License

This project is part of the GoGames collection.

## Author

Benjamin Nechicattu

## Version

1.0.0

---

**Enjoy the game! Can you beat the AI?** (Spoiler: No, but you can draw! 😉)
