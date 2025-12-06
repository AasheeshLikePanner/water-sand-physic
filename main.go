package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Empty = iota
	Sand
	Water
	Wall
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	grid [screenWidth][screenHeight]uint8
}

func isEmptyOrWater(cell uint8) bool {
	return cell == Empty || cell == Water
}

func (g *Game) Update() error {
	drawType := Empty
	if ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		drawType = Wall
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		drawType = Sand
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		drawType = Water
	}

	if drawType != Empty {
		x, y := ebiten.CursorPosition()
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				dx, dy := x+i, y+j
				if dx >= 0 && dx < screenWidth && dy >= 0 && dy < screenHeight {
					if g.grid[dx][dy] != Wall || drawType == Wall {
						g.grid[dx][dy] = uint8(drawType)
					}
				}
			}
		}
	}

	for y := screenHeight - 2; y >= 0; y-- {
		start, end, step := 0, screenWidth, 1
		if rand.Intn(2) == 0 {
			start, end, step = screenWidth-1, -1, -1
		}

		for x := start; x != end; x += step {
			cell := g.grid[x][y]

			if cell == Sand {
				down := g.grid[x][y+1]
				if isEmptyOrWater(down) {
					g.grid[x][y] = down
					g.grid[x][y+1] = Sand
				} else if x > 0 && isEmptyOrWater(g.grid[x-1][y+1]) {
					g.grid[x][y] = g.grid[x-1][y+1]
					g.grid[x-1][y+1] = Sand
				} else if x < screenWidth-1 && isEmptyOrWater(g.grid[x+1][y+1]) {
					g.grid[x][y] = g.grid[x+1][y+1]
					g.grid[x+1][y+1] = Sand
				}

			} else if cell == Water {

				if g.grid[x][y+1] == Empty {
					g.grid[x][y] = Empty
					g.grid[x][y+1] = Water

				} else if g.grid[x][y+1] == Sand {

					splashLeft := rand.Intn(2) == 0

					if splashLeft {
						if x > 0 && y > 0 && g.grid[x-1][y-1] == Empty {

							g.grid[x][y] = Empty
							g.grid[x][y+1] = Water
							g.grid[x-1][y-1] = Sand

						} else if x > 0 && g.grid[x-1][y] == Empty {
							g.grid[x][y] = Empty
							g.grid[x][y+1] = Water
							g.grid[x-1][y] = Sand
						}

					} else { // Splash Right
						if x < screenWidth-1 && y > 0 && g.grid[x+1][y-1] == Empty {

							g.grid[x][y] = Empty
							g.grid[x][y+1] = Water
							g.grid[x+1][y-1] = Sand

						} else if x < screenWidth-1 && g.grid[x+1][y] == Empty {
							g.grid[x][y] = Empty
							g.grid[x][y+1] = Water
							g.grid[x+1][y] = Sand
						}
					}

				} else if x > 0 && g.grid[x-1][y+1] == Empty {
					g.grid[x][y] = Empty
					g.grid[x-1][y+1] = Water
				} else if x < screenWidth-1 && g.grid[x+1][y+1] == Empty {
					g.grid[x][y] = Empty
					g.grid[x+1][y+1] = Water
				} else if x > 0 && g.grid[x-1][y] == Empty {
					g.grid[x][y] = Empty
					g.grid[x-1][y] = Water
				} else if x < screenWidth-1 && g.grid[x+1][y] == Empty {
					g.grid[x][y] = Empty
					g.grid[x+1][y] = Water
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			cell := g.grid[x][y]

			switch cell {
			case Sand:
				screen.Set(x, y, color.RGBA{194, 178, 128, 255})
			case Water:
				screen.Set(x, y, color.RGBA{64, 164, 223, 255})
			case Wall:
				screen.Set(x, y, color.RGBA{100, 100, 100, 255})
			}
		}
	}
	ebitenutil.DebugPrint(screen, "Left: Sand | Right: Water | W+Click: Wall")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{}
	ebiten.SetWindowTitle("Falling Sand Engine")
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
