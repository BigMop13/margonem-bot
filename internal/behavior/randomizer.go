package behavior

import (
	"math"
	"math/rand"
)

// Point represents a 2D coordinate
type Point struct {
	X, Y float64
}

// AddJitter adds random offset to a point
func AddJitter(p Point, maxOffset float64) Point {
	return Point{
		X: p.X + (rand.Float64()*2-1)*maxOffset,
		Y: p.Y + (rand.Float64()*2-1)*maxOffset,
	}
}

// Distance calculates Euclidean distance between two points
func Distance(a, b Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// GeneratePath creates intermediate points for smoother movement
func GeneratePath(from, to Point, stepSize float64) []Point {
	dist := Distance(from, to)
	if dist < stepSize {
		return []Point{to}
	}

	steps := int(math.Ceil(dist / stepSize))
	path := make([]Point, steps)

	for i := 0; i < steps; i++ {
		t := float64(i+1) / float64(steps)
		
		// Add slight curve using bezier
		controlX := (from.X+to.X)/2 + (rand.Float64()*2-1)*20
		controlY := (from.Y+to.Y)/2 + (rand.Float64()*2-1)*20
		
		// Quadratic bezier
		x := (1-t)*(1-t)*from.X + 2*(1-t)*t*controlX + t*t*to.X
		y := (1-t)*(1-t)*from.Y + 2*(1-t)*t*controlY + t*t*to.Y
		
		path[i] = Point{X: x, Y: y}
	}

	return path
}

// RandomOffset returns a random offset within a radius
func RandomOffset(radius float64) Point {
	angle := rand.Float64() * 2 * math.Pi
	r := rand.Float64() * radius
	return Point{
		X: r * math.Cos(angle),
		Y: r * math.Sin(angle),
	}
}

// ShouldTakeBreak determines if a break should be taken based on action count
func ShouldTakeBreak(actionCount, breakEvery int) bool {
	if breakEvery <= 0 {
		return false
	}
	return actionCount%breakEvery == 0
}
