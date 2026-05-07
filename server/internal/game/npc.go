package game

type NPC struct {
	ID string

	X float64
	Y float64

	Alive         bool
	ReachedFinish bool

	TargetX   float64
	WaitTimer float64
	Moving    bool
}
