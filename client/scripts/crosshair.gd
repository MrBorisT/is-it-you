extends Node2D

func _process(delta):
	global_position = get_global_mouse_position()
	queue_redraw()

func _draw():
	draw_line(Vector2(-8, 0), Vector2(8, 0), Color.RED, 2)
	draw_line(Vector2(0, -8), Vector2(0, 8), Color.RED, 2)
	draw_circle(Vector2.ZERO, 12, Color.RED, false, 2)
