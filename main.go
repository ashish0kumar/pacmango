package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rthornton128/goncurses"
)

func main() {
	var levelPath string
	flag.StringVar(&levelPath, "level", "data/level.dat", "Path to level data file")
	flag.Parse()

	game := NewGameState()

	if err := initCurses(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing curses: %v\n", err)
		os.Exit(1)
	}
	defer goncurses.End()

	checkScreenSize()
	game.createWindows()

	// Show intro screen
	game.introScreen()

	if err := game.loadLevel(levelPath); err != nil {
		exitProgram(fmt.Sprintf("Cannot find level file: %s", levelPath))
	}

	game.Invincible = false
	game.mainLoop()

	exitProgram("Level cleared!")
}

// NewGameState creates a new game state with default values
func NewGameState() *GameState {
	return &GameState{
		Lives:       3,
		FreeLife:    1000,
		HowSlow:     3,
		SpeedOfGame: 160,
	}
}

func initCurses() error {
	stdscr, err := goncurses.Init()
	if err != nil {
		return err
	}

	if err := goncurses.StartColor(); err != nil {
		return err
	}

	goncurses.Cursor(0)
	stdscr.Keypad(true)
	stdscr.Timeout(0)
	goncurses.Raw(true)
	goncurses.Echo(false)

	goncurses.UseDefaultColors()

	if goncurses.CanChangeColor() {
		goncurses.InitColor(goncurses.C_BLUE, 0, 0, 500)
	}

	// Initialize color pairs
	goncurses.InitPair(ColorNormal, goncurses.C_WHITE, -1)
	goncurses.InitPair(ColorWall, goncurses.C_BLUE, goncurses.C_BLUE)
	goncurses.InitPair(ColorPellet, goncurses.C_YELLOW, -1)
	goncurses.InitPair(ColorPowerup, goncurses.C_RED, -1)
	goncurses.InitPair(ColorGhostWall, goncurses.C_CYAN, -1)
	goncurses.InitPair(ColorGhost1, goncurses.C_RED, -1)
	goncurses.InitPair(ColorGhost2, goncurses.C_CYAN, -1)
	goncurses.InitPair(ColorGhost3, goncurses.C_MAGENTA, -1)
	goncurses.InitPair(ColorGhost4, goncurses.C_YELLOW, -1)
	goncurses.InitPair(ColorBlueGhost, goncurses.C_WHITE, goncurses.C_BLUE)
	goncurses.InitPair(ColorPacman, goncurses.C_YELLOW, -1)
	goncurses.InitPair(ColorCursor, goncurses.C_BLUE, goncurses.C_YELLOW)

	return nil
}

func checkScreenSize() {
	h, w := goncurses.StdScr().MaxYX()
	if h < 32 || w < 40 {
		goncurses.End()
		fmt.Fprintf(os.Stderr, "\nTerminal too small. Minimum size: 32x40\n")
		os.Exit(1)
	}
}

func (g *GameState) createWindows() {
	maxY, maxX := goncurses.StdScr().MaxYX()

	gameStartY := (maxY - LevelHeight) / 2
	gameStartX := (maxX - LevelWidth) / 2
	statusStartY := gameStartY + LevelHeight + 1
	statusStartX := (maxX - 27) / 2

	var err error
	g.Win, err = goncurses.NewWindow(LevelHeight, LevelWidth, gameStartY, gameStartX)
	if err != nil {
		exitProgram("Error creating game window")
	}
	g.Win.Keypad(true)
	g.Win.Timeout(0)

	g.Status, err = goncurses.NewWindow(3, 27, statusStartY, statusStartX)
	if err != nil {
		exitProgram("Error creating status window")
	}
}

func exitProgram(message string) {
	goncurses.End()
	fmt.Println(message)
	os.Exit(0)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
