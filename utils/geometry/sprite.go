package geometry

type Sprite struct {
	position       Point2D
	hitbox         Rectangle
	factZoom       uint64
	bitmapFileName string
}

func NewSprite(position Point2D, hitbox Rectangle, factzoom uint64, bitmapfile string) *Sprite {
	sprt := Sprite{position, hitbox, factzoom, bitmapfile}
	return &sprt
}

/**
	Getters
**/

func (sprt *Sprite) Position() Point2D {
	return sprt.position
}

func (sprt *Sprite) Hitbox() Rectangle {
	return sprt.hitbox
}

func (sprt *Sprite) FactZoom() uint64 {
	return sprt.factZoom
}

func (sprt *Sprite) BitmapFileName() string {
	return sprt.bitmapFileName
}

/**
	Setters
**/

func (sprt *Sprite) SetPosition(position Point2D) {
	sprt.position = position
}

func (sprt *Sprite) SetHitbox(hitbox Rectangle) {
	sprt.hitbox = hitbox
}

func (sprt *Sprite) SetFactZoom(factzoom uint64) {
	sprt.factZoom = factzoom
}

func (sprt *Sprite) SetBitmapFileName(fileName string) {
	sprt.bitmapFileName = fileName
}

/**
	Move
**/

func (sprt *Sprite) MoveX(delta float64) {
	sprt.position.SetX(sprt.position.X() + delta)

	newHitBoxLeftPoint := sprt.hitbox.PointLeftUpCorner()
	newHitBoxLeftPoint.SetX(newHitBoxLeftPoint.X() + delta)
	sprt.hitbox.SetPointLeftUpCorner(newHitBoxLeftPoint)

	newHitBoxRightPoint := sprt.hitbox.PointRightLowCorner()
	newHitBoxRightPoint.SetX(newHitBoxRightPoint.X() + delta)
	sprt.hitbox.SetPointRightLowCorner(newHitBoxRightPoint)
}

func (sprt *Sprite) MoveY(delta float64) {
	sprt.position.SetY(sprt.position.Y() + delta)

	newHitBoxLeftPoint := sprt.hitbox.PointLeftUpCorner()
	newHitBoxLeftPoint.SetY(newHitBoxLeftPoint.Y() + delta)
	sprt.hitbox.SetPointLeftUpCorner(newHitBoxLeftPoint)

	newHitBoxRightPoint := sprt.hitbox.PointRightLowCorner()
	newHitBoxRightPoint.SetY(newHitBoxRightPoint.Y() + delta)
	sprt.hitbox.SetPointRightLowCorner(newHitBoxRightPoint)
}

/**
	Collisions
**/

func (sprt *Sprite) IsCollision(sprt2 *Sprite) bool {
	hitbox1, hitbox2 := sprt.Hitbox(), sprt2.Hitbox()
	xL1, xL2, yL1, yL2 := hitbox1.PointLeftUpCorner().X(), hitbox2.PointLeftUpCorner().X(), hitbox1.PointLeftUpCorner().Y(), hitbox2.PointLeftUpCorner().Y()
	xR1, xR2, yR1, yR2 := hitbox1.PointRightLowCorner().X(), hitbox2.PointRightLowCorner().X(), hitbox1.PointRightLowCorner().Y(), hitbox2.PointRightLowCorner().Y()

	return !(xL1 >= xR2 ||
		xR1 <= xL2 ||
		yL1 <= yR2 ||
		yR1 >= yL2)
}

func (sprt *Sprite) Collision(sprt2 *Sprite) *Rectangle {
	if !sprt.IsCollision(sprt2) {
		return nil
	}
	hitbox1, hitbox2 := sprt.Hitbox(), sprt2.Hitbox()
	xL1, xL2, yL1, yL2 := hitbox1.PointLeftUpCorner().X(), hitbox2.PointLeftUpCorner().X(), hitbox1.PointLeftUpCorner().Y(), hitbox2.PointLeftUpCorner().Y()
	xR1, xR2, yR1, yR2 := hitbox1.PointRightLowCorner().X(), hitbox2.PointRightLowCorner().X(), hitbox1.PointRightLowCorner().Y(), hitbox2.PointRightLowCorner().Y()

	collRectLeft := NewPoint2D(
		max(xL1, xL2),
		min(yL1, yL2))

	collRectRight := NewPoint2D(
		min(xR1, xR2),
		max(yR1, yR2))

	res := NewRectangle(*collRectLeft, *collRectRight)
	return res
}