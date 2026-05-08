package game

type Player struct {
	ID string

	X float64
	Y float64

	MoveRight bool
	Running   bool

	AimX float64
	AimY float64

	ReachedFinish bool
	Alive         bool
	HasBullet     bool
}
