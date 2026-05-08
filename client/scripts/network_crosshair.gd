extends Node2D

var owner_id := ""
var is_me := false
var alive := true
var has_bullet := true

func set_crosshair_data(id: String, me: bool, is_alive: bool, bullet: bool):
	owner_id = id
	is_me = me
	alive = is_alive
	has_bullet = bullet
	queue_redraw()

func _draw():
	if not alive:
		return

	var color := Color.RED

	if is_me:
		color = Color.GREEN

	if not has_bullet:
		color = Color.GRAY

	draw_line(Vector2(-8, 0), Vector2(8, 0), color, 2)
	draw_line(Vector2(0, -8), Vector2(0, 8), color, 2)
	draw_circle(Vector2.ZERO, 12, color, false, 2)

	draw_string(
		ThemeDB.fallback_font,
		Vector2(14, -10),
		owner_id,
		HORIZONTAL_ALIGNMENT_LEFT,
		-1,
		12,
		color
	)
