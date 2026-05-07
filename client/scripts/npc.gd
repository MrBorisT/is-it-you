extends CharacterBody2D

@export var walk_speed := 55.0
@export var finish_x := 900.0

var target_x := 0.0
var wait_timer := 0.0
var is_moving := false
var is_dead := false
var reached_finish := false

func _ready():
	collision_layer = 0
	collision_mask = 0

	target_x = global_position.x
	pick_next_wait()

func _physics_process(delta):
	if is_dead or reached_finish:
		velocity = Vector2.ZERO
		return

	if global_position.x >= finish_x:
		reached_finish = true
		velocity = Vector2.ZERO
		queue_redraw()
		return

	if not is_moving:
		wait_timer -= delta
		if wait_timer <= 0:
			pick_step()
	else:
		move_towards_target()

	move_and_slide()
	queue_redraw()

func pick_next_wait():
	wait_timer = randf_range(0.3, 1.8)

func pick_step():
	var step_distance := randf_range(15.0, 70.0)
	target_x = min(global_position.x + step_distance, finish_x)
	is_moving = true

func move_towards_target():
	var distance_left := target_x - global_position.x

	if distance_left <= 2.0:
		global_position.x = target_x
		velocity = Vector2.ZERO
		is_moving = false
		pick_next_wait()
		return

	velocity = Vector2.RIGHT * walk_speed

func kill():
	if is_dead:
		return

	is_dead = true
	velocity = Vector2.ZERO
	queue_redraw()

func _draw():
	var color := Color.WHITE

	if is_dead:
		color = Color.DARK_RED
	elif reached_finish:
		color = Color.DIM_GRAY

	draw_circle(Vector2.ZERO, 10, color)
	draw_circle(Vector2.ZERO, 3, Color.BLACK)
