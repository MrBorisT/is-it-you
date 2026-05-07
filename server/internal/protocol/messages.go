package protocol

type ClientMessage struct {
	Type string `json:"type"`

	MoveRight bool `json:"move_right,omitempty"`
	Running   bool `json:"running,omitempty"`

	TargetX float64 `json:"target_x,omitempty"`
	TargetY float64 `json:"target_y,omitempty"`
}

type PlayerState struct {
	ID            string  `json:"id"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Running       bool    `json:"running"`
	ReachedFinish bool    `json:"reached_finish"`
	Alive         bool    `json:"alive"`
	HasBullet     bool    `json:"has_bullet"`
}

type ServerMessage struct {
	Type     string        `json:"type"`
	PlayerID string        `json:"player_id,omitempty"`
	Players  []PlayerState `json:"players,omitempty"`
	GameOver bool          `json:"game_over,omitempty"`
	WinnerID string        `json:"winner_id,omitempty"`
}
