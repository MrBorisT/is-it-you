extends Node

var server_url := "ws://localhost:8080/ws"
var is_host := false
var server_process_id := -1

var debug_mode := false

func _ready():
	get_tree().auto_accept_quit = false

func toggle_debug_mode():
	debug_mode = not debug_mode
	print("Debug mode:", debug_mode)

func stop_host_server():
	if not is_host:
		return

	if server_process_id == -1:
		is_host = false
		return

	print("Stopping hosted server. PID:", server_process_id)

	var err := OS.kill(server_process_id)
	if err != OK:
		print("Failed to kill server process:", err)
	else:
		print("Hosted server stopped.")

	server_process_id = -1
	is_host = false

func _notification(what):
	if what == NOTIFICATION_WM_CLOSE_REQUEST:
		stop_host_server()
		get_tree().quit()
