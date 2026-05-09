extends Node

@export var server_url := "ws://localhost:8080/ws"
@export var player_scene: PackedScene
@export var npc_scene: PackedScene

@onready var players_root := $"../Players"
@onready var status_label := $"../CanvasLayer/StatusLabel"
@onready var npcs_root: Node2D = $"../NPCs"

@export var crosshair_scene: PackedScene
@export var crosshairs_root: Node2D
var crosshair_nodes := {}

var socket := WebSocketPeer.new()
var connected := false
var my_player_id := ""

var player_nodes := {}
var npc_nodes := {}

@export var finish_x := 900.0

var game_over := false
var winner_id := ""

var my_has_bullet := true
var my_alive := true

var phase := "waiting"

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
		
		if Input.is_action_just_pressed("restart"):
			send_restart()

		read_messages()

	elif state == WebSocketPeer.STATE_CLOSED:
		if connected:
			print("Disconnected from server")

		connected = false
		status_label.text = "Disconnected"

func send_restart():
	if phase != "finished":
		return

	var msg := {
		"type": "restart"
	}

	socket.send_text(JSON.stringify(msg))

func send_shoot():
	if phase != "running":
		return

	if not my_alive:
		return

	if not my_has_bullet:
		status_label.text = "No bullet left."
		return

	var mouse_pos: Vector2 = get_parent().get_global_mouse_position()

	var msg := {
		"type": "shoot",
		"target_x": mouse_pos.x,
		"target_y": mouse_pos.y
	}

	socket.send_text(JSON.stringify(msg))

func send_input():
	if phase != "running":
		return

	var mouse_pos: Vector2 = get_parent().get_global_mouse_position()

	var msg := {
		"type": "input",
		"move_right": Input.is_action_pressed("move_right"),
		"running": Input.is_action_pressed("run"),
		"aim_x": mouse_pos.x,
		"aim_y": mouse_pos.y
	}

	socket.send_text(JSON.stringify(msg))

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

	phase = str(msg.get("phase", "waiting"))
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
		var aim_x := float(player_data.get("aim_x", x))
		var aim_y := float(player_data.get("aim_y", y))

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
		if not crosshair_nodes.has(id):
			spawn_crosshair_node(id)

		var crosshair = crosshair_nodes[id]
		crosshair.global_position = Vector2(aim_x, aim_y)

		if crosshair.has_method("set_crosshair_data"):
			crosshair.set_crosshair_data(
				id,
				id == my_player_id,
				alive,
				has_bullet
			)
	remove_missing_players(seen_ids)
	remove_missing_crosshairs(seen_ids)
	
	var npcs: Array = msg.get("npcs", [])
	handle_npcs_state(npcs)
	
	update_status_text()

func handle_npcs_state(npcs: Array):
	var seen_ids := {}

	for npc_data in npcs:
		var id := str(npc_data.get("id", ""))
		var x := float(npc_data.get("x", 0.0))
		var y := float(npc_data.get("y", 0.0))
		var alive := bool(npc_data.get("alive", true))
		var reached_finish := bool(npc_data.get("reached_finish", false))

		seen_ids[id] = true

		if not npc_nodes.has(id):
			spawn_npc_node(id)

		var node = npc_nodes[id]
		node.global_position = Vector2(x, y)

		if node.has_method("set_npc_data"):
			node.set_npc_data(id, alive, reached_finish)

	remove_missing_npcs(seen_ids)

func spawn_npc_node(id: String):
	var node = npc_scene.instantiate()
	npcs_root.add_child(node)
	npc_nodes[id] = node

func remove_missing_npcs(seen_ids: Dictionary):
	var ids_to_remove := []

	for id in npc_nodes.keys():
		if not seen_ids.has(id):
			ids_to_remove.append(id)

	for id in ids_to_remove:
		npc_nodes[id].queue_free()
		npc_nodes.erase(id)

func update_status_text():
	if phase == "waiting":
		status_label.text = "Waiting for second player..."
		return

	if phase == "finished":
		if winner_id == my_player_id:
			status_label.text = "Victory! Press R to restart."
		else:
			status_label.text = "Defeat. Winner: " + winner_id + ". Press R to restart."
		return

	if phase == "running":
		if not my_alive:
			status_label.text = "You are dead. Wait for opponent to finish."
			return

		status_label.text = "Player: " + my_player_id + " | Bullet: " + ("1" if my_has_bullet else "0")
		return

	status_label.text = "Unknown phase: " + phase

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

func spawn_crosshair_node(id: String):
	var node = crosshair_scene.instantiate()
	crosshairs_root.add_child(node)
	crosshair_nodes[id] = node

func remove_missing_crosshairs(seen_ids: Dictionary):
	var ids_to_remove := []

	for id in crosshair_nodes.keys():
		if not seen_ids.has(id):
			ids_to_remove.append(id)

	for id in ids_to_remove:
		crosshair_nodes[id].queue_free()
		crosshair_nodes.erase(id)
