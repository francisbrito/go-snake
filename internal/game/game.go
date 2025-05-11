package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	CellSize     = 20
	Columns      = ScreenWidth / CellSize
	Rows         = ScreenHeight / CellSize
)

type Object struct {
	X, Y, W, H int
}

type Snake struct {
	*Object
	Color  color.Color
	Length int
}

type Game struct {
	snake           *Snake
	backgroundColor color.Color
	gridColor       color.Color
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 31})
	g.drawGrid(screen)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for i := 0; i < Columns; i++ {
		x := float32(i * CellSize)
		vector.StrokeLine(screen, x, 0.0, x, ScreenHeight, 1.0, g.gridColor, true)
	}
	for j := 0; j < Rows; j++ {
		y := float32(j * CellSize)
		vector.StrokeLine(screen, 0.0, y, ScreenWidth, y, 1.0, g.gridColor, true)
	}
}

func New() *Game {
	return &Game{
		backgroundColor: color.Gray{Y: 31},
		gridColor:       color.Gray{Y: 127},
	}
}
