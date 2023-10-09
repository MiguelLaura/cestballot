package geometry

import "math"

type Rectangle struct {
	ptLeftUpCorner   Point2D
	ptRightLowCorner Point2D
}

func NewRectangle(ptLeft Point2D, ptRight Point2D) *Rectangle {
	rect := Rectangle{ptLeft, ptRight}
	return &rect
}

func (rect Rectangle) PointLeftUpCorner() Point2D {
	return rect.ptLeftUpCorner
}

func (rect Rectangle) PointRightLowCorner() Point2D {
	return rect.ptRightLowCorner
}

func (rect *Rectangle) SetPointLeftUpCorner(pt Point2D) {
	rect.ptLeftUpCorner = pt
}

func (rect *Rectangle) SetPointRightLowCorner(pt Point2D) {
	rect.ptRightLowCorner = pt
}

func (rect Rectangle) Height() float64 {
	return math.Abs(rect.ptRightLowCorner.Y() - rect.ptLeftUpCorner.Y())
}

func (rect Rectangle) Width() float64 {
	return math.Abs(rect.ptRightLowCorner.X() - rect.ptLeftUpCorner.X())
}
