package protocol

type ClientMessage struct {
	Type string `json:"type"`

	MoveRight bool `json:"move_right,omitempty"`
	Running   bool `json:"running,omitempty"`

	// Current crosshair/aim position.
	AimX float64 `json:"aim_x,omitempty"`
	AimY float64 `json:"aim_y,omitempty"`

	// Shot target.
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

	AimX float64 `json:"aim_x"`
	AimY float64 `json:"aim_y"`
}

type NPCState struct {
	ID            string  `json:"id"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Alive         bool    `json:"alive"`
	ReachedFinish bool    `json:"reached_finish"`
}

type ServerMessage struct {
	Type     string        `json:"type"`
	PlayerID string        `json:"player_id,omitempty"`
	Players  []PlayerState `json:"players,omitempty"`
	NPCs     []NPCState    `json:"npcs,omitempty"`
	GameOver bool          `json:"game_over,omitempty"`
	WinnerID string        `json:"winner_id,omitempty"`
	Message  string        `json:"message,omitempty"`
}
