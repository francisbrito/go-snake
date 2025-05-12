package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type SnakeDirection int

const (
	Idle SnakeDirection = iota
	Up
	Left
	Down
	Right
)

type Snake struct {
	*Object
	Color     color.Color
	Length    int
	Direction SnakeDirection
}

type Game struct {
	snake           *Snake
	backgroundColor color.Color
	gridColor       color.Color
	tick            int
	tps             int
}

func (g *Game) Update() error {
	g.tick++
	if g.tick%g.tps < 1 {
		g.tick = 0
		switch g.snake.Direction {
		case Up:
			g.snake.Y -= 1
		case Left:
			g.snake.X -= 1
		case Down:
			g.snake.Y += 1
		case Right:
			g.snake.X += 1
		default:
		}
	}
	switch {
	case g.isUpInputPressed():
		g.snake.Direction = Up
	case g.isLeftInputPressed():
		g.snake.Direction = Left
	case g.isDownInputPressed():
		g.snake.Direction = Down
	case g.isRightPressed():
		g.snake.Direction = Right
	}
	return nil
}

func (g *Game) isRightPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight)
}

func (g *Game) isDownInputPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown)
}

func (g *Game) isLeftInputPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft)
}

func (g *Game) isUpInputPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp)
}

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
	g.drawGrid(screen)
	g.drawSnake(screen)
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

func (g *Game) drawSnake(screen *ebiten.Image) {
	x, y, w, h := float32(g.snake.X*CellSize), float32(g.snake.Y*CellSize), float32(g.snake.W*CellSize), float32(g.snake.H*CellSize)
	vector.DrawFilledRect(screen, x, y, w, h, g.snake.Color, true)
}

func New() *Game {
	return &Game{
		backgroundColor: color.Gray{Y: 31},
		gridColor:       color.Gray{Y: 127},
		snake: &Snake{
			Object: &Object{
				X: 10,
				Y: 10,
				W: 1,
				H: 1,
			},
			Color: color.Gray{Y: 255},
		},
		tps: ebiten.TPS() / 16,
	}
}
