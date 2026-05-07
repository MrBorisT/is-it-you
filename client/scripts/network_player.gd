extends Node2D

var player_id := ""
var is_me := false
var is_running := false
var reached_finish := false
var alive := true
var has_bullet := true

func set_player_data(
	id: String,
	me: bool,
	running: bool,
	finish: bool,
	is_alive: bool,
	bullet: bool
):
	player_id = id
	is_me = me
	is_running = running
	reached_finish = finish
	alive = is_alive
	has_bullet = bullet
	queue_redraw()

func _draw():
	var color := Color.WHITE

	if is_me:
		color = Color.GREEN

	if is_running:
		color = Color.RED

	if reached_finish:
		color = Color.GOLD

	if not alive:
		color = Color.DARK_RED

	draw_circle(Vector2.ZERO, 12, color)
	draw_circle(Vector2.ZERO, 4, Color.BLACK)

	if not has_bullet and alive:
		draw_circle(Vector2(12, -12), 3, Color.GRAY)
	elif has_bullet and alive:
		draw_circle(Vector2(12, -12), 3, Color.YELLOW)

	draw_string(
		ThemeDB.fallback_font,
		Vector2(-12, -18),
		player_id,
		HORIZONTAL_ALIGNMENT_LEFT,
		-1,
		12,
		Color.WHITE
	)
