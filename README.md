<h1 align="center">pacmango</h1> 
<p align="center">
Classic pacman for the terminal, written in Go
</p>

<p align="center">
<img src="https://github.com/user-attachments/assets/7e77bd1e-e825-42be-8dba-e5edc2923787" width="70%" alt="pacmango gameplay demo" />
</p>

---

## Features

- **Authentic Pacman gameplay** with pellets, power-ups, ghosts, and scoring
- **9 challenging levels** with progressive difficulty and unique maze layouts
- **Smooth continuous movement** with pending direction system for responsive controls
- **Intelligent ghost AI** that chases and flees based on power-up state
- **Terminal-optimized UI** using goncurses
- **Cross-platform compatibility** - runs anywhere Go runs

## Installation

### Via `go install`

```bash
go install github.com/ashish0kumar/pacmango@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/ashish0kumar/pacmango.git
cd pacmango

# Build the application
go build

# Move to a directory in your PATH
sudo mv pacmango /usr/local/bin/
```

## Usage

```bash
# Start the game
pacmango

# Start from a specific level
pacmango --level 5
```

### Controls

- **Arrow Keys** or **WASD**: Move Pacman
- **P**: Pause/unpause the game
- **Q**: Quit the game

### Gameplay

- Eat all pellets (`.`) to complete the level
- Power-ups (`*`) make you invincible and allow you to eat ghosts
- Avoid ghosts unless you're invincible
- Score points by eating pellets and ghosts

## Audio Features

- Background music plays automatically at game start (requires `assets/music/background.wav`).
- Sound effects play for pellet, powerup, and death events (requires `assets/sounds/eat.wav`, `assets/sounds/powerup.wav`, and `assets/sounds/death.wav`).
- Audio playback uses the [oto](https://github.com/hajimehoshi/oto) and [go-audio/wav](https://github.com/go-audio/wav) Go libraries for reliable cross-platform sound.

### Adding Sound Effects

Place your sound files in:
- `assets/music/background.wav` (background music)
- `assets/sounds/eat.wav` (pellet collection)
- `assets/sounds/powerup.wav` (powerup collection)
- `assets/sounds/death.wav` (death event)

You can use free sound effects from sites like [freesound.org](https://freesound.org/) or create your own.

### Dependencies

Install required Go packages:
```bash
go get github.com/hajimehoshi/oto github.com/go-audio/wav
```

If you encounter issues with audio playback, ensure your system supports audio output from Go programs.

## Acknowledgements

- Original Pacman by Namco
- [goncurses](https://github.com/rthornton128/goncurses) for terminal UI capabilities

## Contributing

Contributions are welcome! If you have ideas for improvements or bug reports, please feel free to open an issue or submit a pull request.

<br>

<p align="center"> 
<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" alt="catppuccin" />
</p>
<p align="center">
    <i><code>&copy 2025-present <a href="https://github.com/ashish0kumar">Ashish Kumar</a></code></i>
</p>
<div align="center">
<a href="https://github.com/ashish0kumar/pacmango/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ashish0kumar/pacmango?style=for-the-badge&color=CBA6F7&logoColor=cdd6f4&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</div>
