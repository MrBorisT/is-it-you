extends Node2D

var npc_id := ""
var alive := true
var reached_finish := false

func set_npc_data(id: String, is_alive: bool, finish: bool):
	npc_id = id
	alive = is_alive
	reached_finish = finish
	queue_redraw()

func _draw():
	var color := Color.WHITE

	if reached_finish:
		color = Color.DIM_GRAY

	if not alive:
		color = Color.DARK_RED

	draw_circle(Vector2.ZERO, 12, color)
	draw_circle(Vector2.ZERO, 3, Color.BLACK)
