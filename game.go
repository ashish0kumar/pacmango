package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rthornton128/goncurses"
)

func (g *GameState) mainLoop() {
	g.drawWindow()
	g.Win.Refresh()
	g.Status.Refresh()
	delay(1000)

	for g.Food > 0 {
		prevPoints := g.Points
		prevInvincible := g.Invincible
		g.movePacman()
		g.drawWindow()
		g.checkCollision()

		// Play eat sound effect (oto)
		if g.Points > prevPoints {
			PlayEatSoundOto()
		}
		// Play powerup sound effect (oto)
		if g.Invincible && !prevInvincible {
			PlayPowerupSoundOto()
		}

		g.moveGhosts()
		g.drawWindow()
		g.checkCollision()

		if g.Points > g.FreeLife {
			g.Lives++
			g.FreeLife *= 2
		}

		g.gameDelay()
	}

	g.drawWindow()
	delay(1000)
}

func (g *GameState) gameDelay() {
	start := time.Now()
	for time.Since(start).Milliseconds() < int64(g.SpeedOfGame) {
		g.getInput()
		time.Sleep(time.Millisecond)
	}
}

func (g *GameState) getInput() {
	ch := g.Win.GetChar()

	switch ch {
	case goncurses.KEY_UP, 'w', 'W':
		g.PendingDir = Direction{-1, 0}
	case goncurses.KEY_DOWN, 's', 'S':
		g.PendingDir = Direction{1, 0}
	case goncurses.KEY_LEFT, 'a', 'A':
		g.PendingDir = Direction{0, -1}
	case goncurses.KEY_RIGHT, 'd', 'D':
		g.PendingDir = Direction{0, 1}
	case 'p', 'P':
		g.pauseGame()
	case 'q', 'Q':
		exitProgram("Bye!")
	}
}

func (g *GameState) checkCollision() {
	for i := 0; i < MaxGhosts; i++ {
		if g.Ghosts[i].Pos == g.Pacman.Pos {
			if g.Invincible {
				// Ghost dies
				g.Points += g.GhostsInARow * 20
				g.Win.ColorOn(int16(ColorPacman))
				g.Win.MovePrint(g.Pacman.Pos.Y, g.Pacman.Pos.X-1, fmt.Sprintf("%d", g.GhostsInARow*20))
				g.GhostsInARow *= 2
				g.Win.Refresh()

				PlayDeathSoundOto() // Play death sound effect (oto)
				delay(1000)
				g.Ghosts[i].Pos = g.Ghosts[i].StartPos
			} else {
				// Pacman dies
				g.Win.ColorOn(int16(ColorPacman))
				g.Win.MovePrint(g.Pacman.Pos.Y, g.Pacman.Pos.X, "X")
				g.Win.Refresh()

				PlayDeathSoundOto() // Play death sound effect (oto)
				g.Lives--
				g.clearStatus()
				delay(1000)

				if g.Lives < 0 {
					exitProgram("Game Over!")
				}

				g.resetPositions()
				g.drawWindow()
				delay(1000)
			}
		}
	}
}

func (g *GameState) resetPositions() {
	// Reset all entities to starting positions
	g.Pacman.Pos = g.Pacman.StartPos
	g.Pacman.Dir = Direction{0, -1}

	for i := range g.Ghosts {
		g.Ghosts[i].Pos = g.Ghosts[i].StartPos
		switch i {
		case 0:
			g.Ghosts[i].Dir = Direction{1, 0}
		case 1:
			g.Ghosts[i].Dir = Direction{-1, 0}
		case 2:
			g.Ghosts[i].Dir = Direction{0, -1}
		case 3:
			g.Ghosts[i].Dir = Direction{0, 1}
		}
	}
}

func (g *GameState) movePacman() {
	// Check if pending direction is now valid
	if g.PendingDir != (Direction{}) {
		nextPos := Position{
			Y: (g.Pacman.Pos.Y + g.PendingDir.Y + LevelHeight) % LevelHeight,
			X: (g.Pacman.Pos.X + g.PendingDir.X + LevelWidth) % LevelWidth,
		}

		if g.Level[nextPos.Y][nextPos.X] != int(CellWall) &&
			g.Level[nextPos.Y][nextPos.X] != int(CellGhostWall) {
			g.Pacman.Dir = g.PendingDir
			g.PendingDir = Direction{}
		}
	}

	// Handle wrapping
	if g.Pacman.Pos.Y == 0 && g.Pacman.Dir.Y == -1 {
		g.Pacman.Pos.Y = LevelHeight - 1
	} else if g.Pacman.Pos.Y == LevelHeight-1 && g.Pacman.Dir.Y == 1 {
		g.Pacman.Pos.Y = 0
	} else if g.Pacman.Pos.X == 0 && g.Pacman.Dir.X == -1 {
		g.Pacman.Pos.X = LevelWidth - 1
	} else if g.Pacman.Pos.X == LevelWidth-1 && g.Pacman.Dir.X == 1 {
		g.Pacman.Pos.X = 0
	} else {
		// Move Pacman
		newPos := Position{
			Y: g.Pacman.Pos.Y + g.Pacman.Dir.Y,
			X: g.Pacman.Pos.X + g.Pacman.Dir.X,
		}

		if newPos.Y >= 0 && newPos.Y < LevelHeight &&
			newPos.X >= 0 && newPos.X < LevelWidth {
			if g.Level[newPos.Y][newPos.X] != int(CellWall) &&
				g.Level[newPos.Y][newPos.X] != int(CellGhostWall) {
				g.Pacman.Pos = newPos
			}
		}
	}

	// Check what Pacman is eating
	cell := CellType(g.Level[g.Pacman.Pos.Y][g.Pacman.Pos.X])
	switch cell {
	case CellPellet:
		g.Level[g.Pacman.Pos.Y][g.Pacman.Pos.X] = int(CellBlank)
		g.Points++
		g.Food--
	case CellPowerup:
		g.Level[g.Pacman.Pos.Y][g.Pacman.Pos.X] = int(CellBlank)
		g.Invincible = true
		if g.GhostsInARow == 0 {
			g.GhostsInARow = 1
		}
		g.TimeLeft = 16 - g.LevelNumber
	}

	// Handle invincibility timer
	if g.Invincible {
		g.TimeLeft--
		if g.TimeLeft <= 0 {
			g.Invincible = false
			g.GhostsInARow = 0
			g.TimeLeft = 0
		}
	}
}

func (g *GameState) moveGhosts() {
	slowerGhosts := 0

	if g.Invincible {
		slowerGhosts++
		if slowerGhosts > g.HowSlow {
			slowerGhosts = 0
		}
	}

	if !g.Invincible || slowerGhosts < g.HowSlow {
		for i := 0; i < MaxGhosts; i++ {
			ghost := &g.Ghosts[i]

			// Handle wrapping
			if ghost.Pos.Y == 0 && ghost.Dir.Y == -1 {
				ghost.Pos.Y = LevelHeight - 1
			} else if ghost.Pos.Y == LevelHeight-1 && ghost.Dir.Y == 1 {
				ghost.Pos.Y = 0
			} else if ghost.Pos.X == 0 && ghost.Dir.X == -1 {
				ghost.Pos.X = LevelWidth - 1
			} else if ghost.Pos.X == LevelWidth-1 && ghost.Dir.X == 1 {
				ghost.Pos.X = 0
			} else {
				// Determine valid directions
				checkSides := [4]bool{}
				y, x := ghost.Pos.Y, ghost.Pos.X

				if y+1 < LevelHeight && g.Level[y+1][x] != int(CellWall) {
					checkSides[0] = true // down
				}
				if y-1 >= 0 && g.Level[y-1][x] != int(CellWall) {
					checkSides[1] = true // up
				}
				if x+1 < LevelWidth && g.Level[y][x+1] != int(CellWall) {
					checkSides[2] = true // right
				}
				if x-1 >= 0 && g.Level[y][x-1] != int(CellWall) {
					checkSides[3] = true // left
				}

				// Don't do 180 unless we have to
				validMoves := 0
				for _, valid := range checkSides {
					if valid {
						validMoves++
					}
				}

				if validMoves > 1 {
					if ghost.Dir.Y == 1 {
						checkSides[1] = false
					} else if ghost.Dir.Y == -1 {
						checkSides[0] = false
					} else if ghost.Dir.X == 1 {
						checkSides[3] = false
					} else if ghost.Dir.X == -1 {
						checkSides[2] = false
					}
				}

				// Choose direction
				for attempts := 0; attempts < 100; attempts++ {
					direction := rand.Intn(4)

					if checkSides[direction] {
						switch direction {
						case 0: // down
							ghost.Dir = Direction{1, 0}
						case 1: // up
							ghost.Dir = Direction{-1, 0}
						case 2: // right
							ghost.Dir = Direction{0, 1}
						case 3: // left
							ghost.Dir = Direction{0, -1}
						}
						break
					} else {
						// AI behavior
						if !g.Invincible {
							// Chase Pacman
							if g.Pacman.Pos.Y > ghost.Pos.Y && checkSides[0] {
								ghost.Dir = Direction{1, 0}
								break
							} else if g.Pacman.Pos.Y < ghost.Pos.Y && checkSides[1] {
								ghost.Dir = Direction{-1, 0}
								break
							} else if g.Pacman.Pos.X > ghost.Pos.X && checkSides[2] {
								ghost.Dir = Direction{0, 1}
								break
							} else if g.Pacman.Pos.X < ghost.Pos.X && checkSides[3] {
								ghost.Dir = Direction{0, -1}
								break
							}
						} else {
							// Run away from Pacman
							if g.Pacman.Pos.Y > ghost.Pos.Y && checkSides[1] {
								ghost.Dir = Direction{-1, 0}
								break
							} else if g.Pacman.Pos.Y < ghost.Pos.Y && checkSides[0] {
								ghost.Dir = Direction{1, 0}
								break
							} else if g.Pacman.Pos.X > ghost.Pos.X && checkSides[3] {
								ghost.Dir = Direction{0, -1}
								break
							} else if g.Pacman.Pos.X < ghost.Pos.X && checkSides[2] {
								ghost.Dir = Direction{0, 1}
								break
							}
						}
					}
				}

				// Move ghost
				ghost.Pos.Y += ghost.Dir.Y
				ghost.Pos.X += ghost.Dir.X
			}
		}
	}
}

func (g *GameState) introScreen() {
	const asciiTitle = `pacmango`

	g.Win.ColorOn(int16(ColorPacman))

	lines := strings.Split(asciiTitle, "\n")
	startY := (LevelHeight-len(lines))/2 - 2

	for i, line := range lines {
		if line != "" {
			startX := (LevelWidth - len(line)) / 2
			if startX < 0 {
				startX = 0
			}
			g.Win.MovePrint(startY+i, startX, line)
		}
	}

	g.Win.ColorOn(int16(ColorNormal))
	g.Win.MovePrint(LevelHeight-3, (LevelWidth-16)/2, "Press any key...")
	g.Win.Refresh()

	for {
		ch := g.Win.GetChar()
		if ch != 0 {
			break
		}
		delay(50)
	}
}

func (g *GameState) pauseGame() {
	g.Win.ColorOn(int16(ColorPacman))
	g.Win.MovePrint(12, 10, "        ")
	g.Win.MovePrint(13, 10, " PAUSED ")
	g.Win.MovePrint(14, 10, "        ")
	g.Win.Refresh()

	for {
		ch := g.Win.GetChar()
		if ch != 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (g *GameState) showLevelComplete() {
	g.Win.ColorOn(int16(ColorPacman))
	g.Win.MovePrint(12, 4, "                    ")
	g.Win.MovePrint(13, 4, fmt.Sprintf("  LEVEL %d COMPLETE   ", g.CurrentLevel))
	g.Win.MovePrint(14, 4, "                    ")
	g.Win.Refresh()

	delay(2000)
}
