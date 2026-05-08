package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server ServerConfig
	Game   GameConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type GameConfig struct {
	MaxPlayers int
	TickRate   int

	StartX  float64
	StartY  float64
	RowGap  float64
	FinishX float64

	WalkSpeed float64
	RunSpeed  float64
	HitRadius float64

	NPCCount     int
	NPCWalkSpeed float64
	NPCMinStep   float64
	NPCMaxStep   float64
	NPCMinWait   float64
	NPCMaxWait   float64
	NPCMinY      float64
	NPCMaxY      float64
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Game: GameConfig{
			MaxPlayers: 2,
			TickRate:   20,

			StartX:  80,
			StartY:  220,
			RowGap:  50,
			FinishX: 900,

			WalkSpeed: 80,
			RunSpeed:  180,
			HitRadius: 18,

			NPCCount:     25,
			NPCWalkSpeed: 55,
			NPCMinStep:   15,
			NPCMaxStep:   70,
			NPCMinWait:   0.3,
			NPCMaxWait:   1.8,
			NPCMinY:      140,
			NPCMaxY:      420,
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	values, err := readKeyValues(path)
	if err != nil {
		return cfg, err
	}

	cfg.Server.Host = getString(values, "host", cfg.Server.Host)
	cfg.Server.Port = getString(values, "port", cfg.Server.Port)

	cfg.Game.MaxPlayers = getInt(values, "max_players", cfg.Game.MaxPlayers)
	cfg.Game.TickRate = getInt(values, "tick_rate", cfg.Game.TickRate)

	cfg.Game.StartX = getFloat(values, "start_x", cfg.Game.StartX)
	cfg.Game.StartY = getFloat(values, "start_y", cfg.Game.StartY)
	cfg.Game.RowGap = getFloat(values, "row_gap", cfg.Game.RowGap)
	cfg.Game.FinishX = getFloat(values, "finish_x", cfg.Game.FinishX)

	cfg.Game.WalkSpeed = getFloat(values, "walk_speed", cfg.Game.WalkSpeed)
	cfg.Game.RunSpeed = getFloat(values, "run_speed", cfg.Game.RunSpeed)
	cfg.Game.HitRadius = getFloat(values, "hit_radius", cfg.Game.HitRadius)

	cfg.Game.NPCCount = getInt(values, "npc_count", cfg.Game.NPCCount)
	cfg.Game.NPCWalkSpeed = getFloat(values, "npc_walk_speed", cfg.Game.NPCWalkSpeed)
	cfg.Game.NPCMinStep = getFloat(values, "npc_min_step", cfg.Game.NPCMinStep)
	cfg.Game.NPCMaxStep = getFloat(values, "npc_max_step", cfg.Game.NPCMaxStep)
	cfg.Game.NPCMinWait = getFloat(values, "npc_min_wait", cfg.Game.NPCMinWait)
	cfg.Game.NPCMaxWait = getFloat(values, "npc_max_wait", cfg.Game.NPCMaxWait)
	cfg.Game.NPCMinY = getFloat(values, "npc_min_y", cfg.Game.NPCMinY)
	cfg.Game.NPCMaxY = getFloat(values, "npc_max_y", cfg.Game.NPCMaxY)

	return cfg, nil
}

func (c Config) Addr() string {
	return c.Server.Host + ":" + c.Server.Port
}

func readKeyValues(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	values := make(map[string]string)

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line %d: %q", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key on line %d", lineNumber)
		}

		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan config: %w", err)
	}

	return values, nil
}

func getString(values map[string]string, key string, fallback string) string {
	value, ok := values[key]
	if !ok || value == "" {
		return fallback
	}

	return value
}

func getInt(values map[string]string, key string, fallback int) int {
	raw, ok := values[key]
	if !ok || raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func getFloat(values map[string]string, key string, fallback float64) float64 {
	raw, ok := values[key]
	if !ok || raw == "" {
		return fallback
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}

	return value
}
