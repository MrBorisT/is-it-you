package game

type Player struct {
	ID            string
	X             float64
	Y             float64
	MoveRight     bool
	Running       bool
	ReachedFinish bool
	Alive         bool
	HasBullet     bool
}
