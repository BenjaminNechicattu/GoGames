# Fyne Snake Game

A classic Snake game built with Go and Fyne GUI framework, featuring a beautiful interface with googly-eyed snake, smooth animations, and a high score system.

## Features

- **Classic Snake Gameplay**: Control the snake with arrow keys to eat food and grow
- **Smooth Graphics**: Beautiful circular snake body with animated googly eyes
- **High Score System**: Automatically saves your best scores with player names
- **Pause/Resume**: Press SPACE to pause the game anytime
- **Progressive Difficulty**: Speed increases every 5 points
- **Keyboard Shortcuts**:
  - Arrow Keys: Control snake direction
  - SPACE: Pause/Resume game
  - ENTER: Start game / Restart after game over

## Score History System

The game automatically tracks ALL your games with complete history and saves it to:
```
~/.local/share/fyne-snake/score-history.json
```

Features:
- **Track All Plays**: Every game is recorded with player name, score, date, and time
- **High Score Tracking**: The highest score is prominently displayed
- **View Score History**: Click "View Scores" button to see your last 20 games
- **Persistent Storage**: All data is saved between sessions

When you beat the current high score:
1. A dialog will appear asking for your name
2. Your score and name will be saved automatically
3. The high score is displayed on the splash screen and game over screen

The history stores:
- Player name
- Score achieved
- Date and time of the game

## How to Run

```bash
go run main.go
```

## How to Build

```bash
go build -o snake
```

## Controls

- **↑ ↓ ← →**: Move the snake
- **SPACE**: Pause/Resume
- **ENTER**: Start game or restart after game over

## Game Rules

- Use Arrow Keys to move the snake
- Eat food to grow and score
- Don't hit the walls or yourself
- Press SPACE to pause/resume
- Speed increases every 5 points

## Author

by SAPPHIRE_KNIGHT

Developed in GoLang, powered by Fyne GUI
