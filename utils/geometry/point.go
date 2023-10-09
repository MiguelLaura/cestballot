package geometry

import "math"

type Point2D struct {
	x float64
	y float64
}

func NewPoint2D(x, y float64) *Point2D {
	point := Point2D{x: x, y: y}

	return &point
}

func (pt Point2D) X() float64 {
	return pt.x
}

func (pt Point2D) Y() float64 {
	return pt.y
}

func (pt *Point2D) SetX(x float64) {
	pt.x = x
}

func (pt *Point2D) SetY(y float64) {
	pt.y = y
}

func (pt Point2D) Clone() Point2D {
	return Point2D{pt.X(), pt.Y()}
}

func (pt Point2D) Module() float64 {
	return math.Sqrt(pt.X()*pt.X() + pt.Y()*pt.Y())
}
