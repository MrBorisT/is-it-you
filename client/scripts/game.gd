extends Node2D

@export var npc_scene: PackedScene
@export var npc_count := 25

@export var start_x := 80.0
@export var finish_x := 900.0
@export var min_y := 120.0
@export var max_y := 500.0

@export var hit_radius := 14.0

@onready var player := $Player
@onready var npcs := $NPCs
@onready var status_label := $CanvasLayer/UI/StatusLabel
@onready var bullet_label := $CanvasLayer/UI/BulletLabel

var has_bullet := true
var game_over := false

func _ready():
	randomize()
	setup_player()
	spawn_npcs()
	update_ui("Reach the finish. One bullet.")

func _process(delta):
	if Input.is_action_just_pressed("restart"):
		get_tree().reload_current_scene()
		return
	if game_over:
		return

	if Input.is_action_just_pressed("shoot"):
		shoot()

	check_game_state()
	queue_redraw()

func setup_player():
	player.global_position = Vector2(start_x, 300)

	if "finish_x" in player:
		player.finish_x = finish_x

func spawn_npcs():
	for i in npc_count:
		var npc = npc_scene.instantiate()

		var y := randf_range(min_y, max_y)
		npc.global_position = Vector2(start_x, y)

		if "finish_x" in npc:
			npc.finish_x = finish_x

		npcs.add_child(npc)

func shoot():
	if not has_bullet:
		update_ui("No bullet left.")
		return

	has_bullet = false
	update_ui("Shot fired.")

	var mouse_pos := get_global_mouse_position()
	var closest_target = find_target_under_crosshair(mouse_pos)

	if closest_target == null:
		print("Miss")
		update_ui("Miss. No bullet left.")
		return

	if closest_target.has_method("kill"):
		closest_target.kill()
		print("Hit target")
		update_ui("Hit. No bullet left.")

func find_target_under_crosshair(mouse_pos: Vector2):
	var closest = null
	var closest_distance := hit_radius

	for npc in npcs.get_children():
		if npc.is_dead:
			continue

		var distance: float = npc.global_position.distance_to(mouse_pos)
		if distance <= closest_distance:
			closest = npc
			closest_distance = distance

	return closest

func check_game_state():
	if player.is_dead:
		game_over = true
		update_ui("You died. Defeat.")
		return

	if player.has_won:
		game_over = true
		update_ui("You reached the finish. Victory.")
		return

func update_ui(status_text: String):
	status_label.text = status_text + " Press R to restart."
	bullet_label.text = "Bullet: " + ("1" if has_bullet else "0")

func _draw():
	draw_line(
		Vector2(start_x, min_y - 40),
		Vector2(start_x, max_y + 40),
		Color.BLUE,
		3
	)

	draw_line(
		Vector2(finish_x, min_y - 40),
		Vector2(finish_x, max_y + 40),
		Color.GREEN,
		3
	)
