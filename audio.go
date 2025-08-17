package main

import (
	"encoding/binary"
	"fmt"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto"
	"io"
	"os"
)

var otoCtx *oto.Context

func InitializeOto() {
	if otoCtx != nil {
		return
	}
	otoCtx, _ = oto.NewContext(44100, 2, 2, 8192)
}

func PlayWAVOto(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening WAV file:", path, err)
		return
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		fmt.Println("Invalid WAV file:", path)
		return
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		fmt.Println("Error decoding WAV file:", path, err)
		return
	}

	byteBuf := make([]byte, len(buf.Data)*2)
	for i, v := range buf.Data {
		binary.LittleEndian.PutUint16(byteBuf[i*2:], uint16(v))
	}

	player := otoCtx.NewPlayer()
	if _, err := player.Write(byteBuf); err != nil && err != io.EOF {
		fmt.Println("Error writing audio data:", err)
	}
	player.Close()
}

func PlayBackgroundMusicOto() {
	go func() {
		for {
			PlayWAVOto("assets/music/background.wav")
		}
	}()
}

func PlayEatSoundOto() {
	go PlayWAVOto("assets/sounds/eat.wav")
}
func PlayDeathSoundOto() {
	go PlayWAVOto("assets/sounds/death.wav")
}
func PlayPowerupSoundOto() {
	go PlayWAVOto("assets/sounds/powerup.wav")
}
