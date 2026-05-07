extends Node

@export var server_url := "ws://localhost:8080/ws"
@export var player_scene: PackedScene

@onready var players_root := $"../Players"
@onready var status_label := $"../CanvasLayer/StatusLabel"

var socket := WebSocketPeer.new()
var connected := false
var my_player_id := ""

var player_nodes := {}

@export var finish_x := 900.0

var game_over := false
var winner_id := ""

var my_has_bullet := true
var my_alive := true

func _ready():
	var err := socket.connect_to_url(server_url)

	if err != OK:
		status_label.text = "WebSocket error: " + str(err)
		print("WebSocket connect error:", err)
		return

	status_label.text = "Connecting..."
	print("Connecting to server:", server_url)

func _process(delta):
	socket.poll()

	var state := socket.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:
		if not connected:
			connected = true
			status_label.text = "Connected"
			print("Connected to server")

		send_input()

		if Input.is_action_just_pressed("shoot"):
			send_shoot()

		read_messages()

	elif state == WebSocketPeer.STATE_CLOSED:
		if connected:
			print("Disconnected from server")

		connected = false
		status_label.text = "Disconnected"

func send_shoot():
	if game_over:
		return

	if not my_alive:
		return

	if not my_has_bullet:
		status_label.text = "No bullet left."
		return

	var mouse_pos := get_viewport().get_mouse_position()

	var msg := {
		"type": "shoot",
		"target_x": mouse_pos.x,
		"target_y": mouse_pos.y
	}

	socket.send_text(JSON.stringify(msg))

func send_input():
	if game_over:
		return

	var msg := {
		"type": "input",
		"move_right": Input.is_action_pressed("move_right"),
		"running": Input.is_action_pressed("run")
	}

	var json := JSON.stringify(msg)
	socket.send_text(json)

func read_messages():
	while socket.get_available_packet_count() > 0:
		var packet := socket.get_packet().get_string_from_utf8()
		handle_message(packet)

func handle_message(packet: String):
	var parsed = JSON.parse_string(packet)

	if typeof(parsed) != TYPE_DICTIONARY:
		print("Invalid message:", packet)
		return

	var msg: Dictionary = parsed

	match msg.get("type", ""):
		"welcome":
			handle_welcome(msg)

		"state":
			handle_state(msg)

		_:
			print("Unknown message:", msg)

func handle_welcome(msg: Dictionary):
	my_player_id = msg.get("player_id", "")
	status_label.text = "Connected as " + my_player_id
	print("My player id:", my_player_id)

func handle_state(msg: Dictionary):
	var players: Array = msg.get("players", [])

	game_over = bool(msg.get("game_over", false))
	winner_id = str(msg.get("winner_id", ""))

	var seen_ids := {}

	for player_data in players:
		var id := str(player_data.get("id", ""))
		var x := float(player_data.get("x", 0.0))
		var y := float(player_data.get("y", 0.0))
		var running := bool(player_data.get("running", false))
		var reached_finish := bool(player_data.get("reached_finish", false))
		var alive := bool(player_data.get("alive", true))
		var has_bullet := bool(player_data.get("has_bullet", true))

		if id == my_player_id:
			my_alive = alive
			my_has_bullet = has_bullet

		seen_ids[id] = true

		if not player_nodes.has(id):
			spawn_player_node(id)

		var node = player_nodes[id]
		node.global_position = Vector2(x, y)

		if node.has_method("set_player_data"):
			node.set_player_data(
				id,
				id == my_player_id,
				running,
				reached_finish,
				alive,
				has_bullet
			)

	remove_missing_players(seen_ids)
	update_status_text()

func update_status_text():
	if game_over:
		if winner_id == my_player_id:
			status_label.text = "Victory!"
		else:
			status_label.text = "Defeat. Winner: " + winner_id
		return

	if not my_alive:
		status_label.text = "You are dead."
		return

	status_label.text = "Player: " + my_player_id + " | Bullet: " + ("1" if my_has_bullet else "0")

func spawn_player_node(id: String):
	var node = player_scene.instantiate()
	players_root.add_child(node)
	player_nodes[id] = node

func remove_missing_players(seen_ids: Dictionary):
	var ids_to_remove := []

	for id in player_nodes.keys():
		if not seen_ids.has(id):
			ids_to_remove.append(id)

	for id in ids_to_remove:
		player_nodes[id].queue_free()
		player_nodes.erase(id)
