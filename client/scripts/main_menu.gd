extends Control

@export var game_scene_path := "res://scenes/network_test.tscn"

@export var server_executable_path := "res://external/server/is-it-you-server.exe"
@export var server_config_path := "res://external/server/config.cfg"

@onready var ip_line_edit: LineEdit = $VBoxContainer/IPLineEdit
@onready var port_line_edit: LineEdit = $VBoxContainer/PortLineEdit
@onready var host_button: Button = $VBoxContainer/HostButton
@onready var join_button: Button = $VBoxContainer/JoinButton
@onready var status_label: Label = $VBoxContainer/StatusLabel

func _ready():
	host_button.pressed.connect(_on_host_pressed)
	join_button.pressed.connect(_on_join_pressed)

	if port_line_edit.text.strip_edges() == "":
		port_line_edit.text = "8080"

	show_local_ips()

func _on_host_pressed():
	var port := get_port()
	if port == "":
		status_label.text = "Invalid port."
		return

	var ok := start_local_server()
	if not ok:
		return

	AppState.server_url = "ws://127.0.0.1:%s/ws" % port

	status_label.text = "Hosting on " + AppState.server_url
	get_tree().change_scene_to_file(game_scene_path)

func _on_join_pressed():
	var ip := ip_line_edit.text.strip_edges()
	var port := get_port()

	if ip == "":
		status_label.text = "Enter host IP."
		return

	if port == "":
		status_label.text = "Invalid port."
		return

	AppState.is_host = false
	AppState.server_url = "ws://%s:%s/ws" % [ip, port]

	status_label.text = "Joining " + AppState.server_url
	get_tree().change_scene_to_file(game_scene_path)

func get_port() -> String:
	var port := port_line_edit.text.strip_edges()

	if port == "":
		return ""

	if not port.is_valid_int():
		return ""

	var port_number := int(port)
	if port_number <= 0 or port_number > 65535:
		return ""

	return port

func start_local_server() -> bool:
	var executable_path := ProjectSettings.globalize_path(server_executable_path)
	var config_path := ProjectSettings.globalize_path(server_config_path)

	var args := PackedStringArray([
		"-config",
		config_path
	])

	var pid := OS.create_process(executable_path, args)

	if pid == -1:
		status_label.text = "Failed to start server."
		return false

	AppState.server_process_id = pid
	AppState.is_host = true

	status_label.text = "Server started. PID: " + str(pid)
	print("Server started. PID:", pid)

	return true

func show_local_ips():
	var addresses := IP.get_local_addresses()
	var useful := []

	for address in addresses:
		if address.begins_with("192.168.") or address.begins_with("10.") or address.begins_with("172."):
			useful.append(address)

	if useful.size() > 0:
		status_label.text = "Your LAN IP may be: " + ", ".join(useful)
	else:
		status_label.text = "Could not detect LAN IP."
