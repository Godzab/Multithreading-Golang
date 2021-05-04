package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenwidth, screenheight = 640, 360
	boidCount                 = 500
	viewRadius                = 13
	adjrate                   = 0.015
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid
	boidMap [screenwidth + 1][screenheight + 1]int
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		screen.Set(int(boid.Position.x+1), int(boid.Position.y), green)
		screen.Set(int(boid.Position.x-1), int(boid.Position.y), green)
		screen.Set(int(boid.Position.x), int(boid.Position.y+1), green)
		screen.Set(int(boid.Position.x), int(boid.Position.y-1), green)
	}
}

func (g Game) Layout(_, _ int) (w, h int) {
	return screenwidth, screenheight
}

func main() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCount; i++ {
		CreateBoid(i)
	}
	ebiten.SetWindowSize(screenwidth*2, screenheight*2)
	ebiten.SetWindowTitle("Boids In A Box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
