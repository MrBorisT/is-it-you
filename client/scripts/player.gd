extends CharacterBody2D

@export var walk_speed := 80.0
@export var run_speed := 180.0
@export var finish_x := 900.0

var is_running := false
var is_dead := false
var has_won := false

func _ready():
	collision_layer = 0
	collision_mask = 0

func _physics_process(delta):
	if is_dead or has_won:
		velocity = Vector2.ZERO
		return

	var wants_to_move := Input.is_action_pressed("move_right")

	is_running = wants_to_move and Input.is_action_pressed("run")

	if wants_to_move:
		var speed := run_speed if is_running else walk_speed
		velocity = Vector2.RIGHT * speed
	else:
		velocity = Vector2.ZERO

	move_and_slide()

	if global_position.x >= finish_x:
		has_won = true
		print("Player reached finish. You win.")

	queue_redraw()

func kill():
	if is_dead or has_won:
		return

	is_dead = true
	velocity = Vector2.ZERO
	print("Player died.")

	queue_redraw()

func _draw():
	var color := Color.WHITE

	if is_dead:
		color = Color.DARK_RED
	elif has_won:
		color = Color.GREEN
	elif is_running:
		color = Color.RED

	draw_circle(Vector2.ZERO, 10, color)
	draw_circle(Vector2.ZERO, 4, Color.BLACK)
