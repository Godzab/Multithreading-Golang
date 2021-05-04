package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	Position Vector2D
	Velocity Vector2D
	Id       int
}

func (b Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Microsecond)
	}
}

func (b Boid) CalculateAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
				}
			}
		}
	}
	accel := Vector2D{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}
	return accel
}

func (b Boid) moveOne() {
	b.Velocity = b.Velocity.Add(b.CalculateAcceleration())
	boidMap[int(b.Position.x)][int(b.Position.y)] = -1
	b.Position = b.Position.Add(b.Velocity)
	boidMap[int(b.Position.x)][int(b.Position.y)] = b.Id
	next := b.Position.Add(b.Velocity)
	if next.x >= screenwidth || next.x < 0 {
		b.Velocity = Vector2D{-b.Velocity.x, b.Velocity.y}
	}
	if next.y >= screenheight || next.y < 0 {
		b.Velocity = Vector2D{b.Velocity.x, -b.Velocity.y}
	}
}

func CreateBoid(id int) {
	b := Boid{
		Position: Vector2D{rand.Float64() * screenwidth, rand.Float64() * screenheight},
		Velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		Id:       id,
	}

	boids[id] = &b
	boidMap[int(b.Position.x)][int(b.Position.y)] = b.Id
	go b.start()
}
