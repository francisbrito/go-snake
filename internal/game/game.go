package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (d SnakeDirection) String() string {
	switch d {
	case NoDirection:
		return "No Direction"
	case Up:
		return "Up"
	case Left:
		return "Left"
	case Down:
		return "Down"
	case Right:
		return "Right"
	default:
		return fmt.Sprintf("Unknown: %d", d)
	}
}

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

func (s State) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Over:
		return "Game Over"
	default:
		return fmt.Sprintf("Unknown: %d", s)
	}
}

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
	grid            [Columns][Rows]bool
}

func (g *Game) Update() error {
	g.tick++
	if g.tick%g.tps == 0 {
		g.tick = 0
		g.grid[g.snake.X][g.snake.Y] = false
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
		if g.snake.X >= 0 && g.snake.X < Columns && g.snake.Y >= 0 && g.snake.Y < Rows {
			g.grid[g.snake.X][g.snake.Y] = true
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
	g.printDebugInfo(screen)
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
	g.grid[g.food.X][g.food.Y] = false
	g.food.X = g.rand.Intn(Columns)
	g.food.Y = g.rand.Intn(Rows)
	g.grid[g.food.X][g.food.Y] = true
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

func (g *Game) printDebugInfo(screen *ebiten.Image) {
	msg := fmt.Sprintf("FPS: %2.2f\nTPS: %2.2f (%d)\nState: %s\nS. Speed: %2d\nS. Direction: %s\nS. Tail Length: %2d",
		ebiten.ActualFPS(), ebiten.ActualTPS(), ebiten.TPS(), g.state, g.tps, g.snake.Direction, g.snake.Length)
	ebitenutil.DebugPrintAt(screen, msg, 16, 16)
	//sb := strings.Builder{}
	//for i := 0; i < Rows; i++ {
	//	for j := 0; j < Columns; j++ {
	//		if g.grid[j][i] && g.food.X == j && g.food.Y == i {
	//			sb.WriteRune('*')
	//		} else if g.grid[j][i] {
	//			sb.WriteRune('S')
	//		} else {
	//			sb.WriteRune('_')
	//		}
	//	}
	//	sb.WriteRune('\n')
	//}
	//ebitenutil.DebugPrintAt(screen, sb.String(), ScreenWidth/2, 16)
}

func New() *Game {
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))
	fx, fy, sx, sy := r.Intn(Columns), r.Intn(Rows), r.Intn(Columns), r.Intn(Rows)
	food := &Food{
		Object: &Object{
			X: fx,
			Y: fy,
			W: 1,
			H: 1,
		},
		Color: color.RGBA{
			R: 255,
			G: 127,
			B: 0,
			A: 255,
		},
	}
	g := &Game{
		backgroundColor: color.Gray{Y: 31},
		gridColor:       color.Gray{Y: 63},
		rand:            r,
		seed:            seed,
		snake: &Snake{
			Object: &Object{
				X: sx,
				Y: sy,
				W: 1,
				H: 1,
			},
			Color:     color.Gray{Y: 255},
			Direction: NoDirection,
		},
		tps:  10,
		food: food,
	}
	g.grid[food.X][food.Y] = true
	return g
}
