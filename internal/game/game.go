package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
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
	NoDirection SnakeDirection = iota
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

type State int

const (
	Idle State = iota
	Started
	Over
)

type Food struct {
	*Object
	Color color.Color
}

type Game struct {
	snake           *Snake
	backgroundColor color.Color
	gridColor       color.Color
	tick            int
	tps             int
	rand            *rand.Rand
	seed            int64
	state           State
	food            *Food
}

func (g *Game) Update() error {
	g.tick++
	if g.tick%g.tps == 0 {
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
	case g.isRightInputPressed():
		g.snake.Direction = Right
	}
	g.checkCollisions()
	if g.state == Over {
		g.reset()
	}
	return nil
}

func (g *Game) isRightInputPressed() bool {
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
	g.drawFood(screen)
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

func (g *Game) checkCollisions() {
	x, y := g.snake.X, g.snake.Y
	// check walls
	if x < 0 || x >= Columns || y < 0 || y >= Rows {
		g.state = Over
		g.reset()
	}
	// todo: check tail
	// check food
	if x == g.food.X && y == g.food.Y {
		g.snake.Length++
		g.spawnFood()
	}
}

func (g *Game) spawnFood() {
	// todo: ensure food does not spawn on occupied cells
	g.food.X = g.rand.Intn(Columns)
	g.food.Y = g.rand.Intn(Rows)
}

func (g *Game) reset() {
	g.tick = 0
	g.seed = rand.Int63()
	g.rand = rand.New(rand.NewSource(g.seed))
	g.state = Idle
	g.snake.X = 10
	g.snake.Y = 10
	g.snake.Direction = NoDirection
}

func (g *Game) drawFood(screen *ebiten.Image) {
	x, y, w, h := float32(g.food.X*CellSize), float32(g.food.Y*CellSize), float32(g.food.W*CellSize), float32(g.food.H*CellSize)
	vector.DrawFilledRect(screen, x, y, w, h, g.food.Color, true)
}

func New() *Game {
	seed := rand.Int63()
	return &Game{
		backgroundColor: color.Gray{Y: 31},
		gridColor:       color.Gray{Y: 63},
		rand:            rand.New(rand.NewSource(seed)),
		seed:            seed,
		snake: &Snake{
			Object: &Object{
				X: 10,
				Y: 10,
				W: 1,
				H: 1,
			},
			Color:     color.Gray{Y: 255},
			Direction: NoDirection,
		},
		tps: 10,
		food: &Food{
			Object: &Object{
				X: 5,
				Y: 5,
				W: 1,
				H: 1,
			},
			Color: color.RGBA{
				R: 255,
				G: 127,
				B: 0,
				A: 255,
			},
		},
	}
}
