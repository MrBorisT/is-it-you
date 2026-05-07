extends Node2D

@export var finish_x := 900.0
@export var min_y := 150.0
@export var max_y := 400.0

func _draw():
	draw_line(
		Vector2(finish_x, min_y),
		Vector2(finish_x, max_y),
		Color.GREEN,
		3
	)
