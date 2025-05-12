package main

import (
	"github.com/francisbrito/snake/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowTitle("Snake")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	//ebiten.MaximizeWindow()
	g := game.New()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
