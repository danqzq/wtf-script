package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type MinMax[T int64 | uint64 | float64] struct {
	Min T `json:"min"`
	Max T `json:"max"`
}

type TypeDefaultRanges struct {
	Int      MinMax[int64]   `json:"int"`
	Uint     MinMax[uint64]  `json:"uint"`
	Float    MinMax[float64] `json:"float"`
	Unofloat MinMax[float64] `json:"unofloat"`
}

type StringDefaults struct {
	Charset string         `json:"charset"`
	Length  MinMax[uint64] `json:"length"`
}

type Config struct {
	TypeDefaultRanges
	StringDefaults
}

var DefaultConfig = Config{
	TypeDefaultRanges: TypeDefaultRanges{
		Int: MinMax[int64]{
			Min: int64(-1000),
			Max: int64(1000),
		},
		Uint: MinMax[uint64]{
			Min: uint64(0),
			Max: uint64(2000),
		},
		Float: MinMax[float64]{
			Min: -1000.0,
			Max: 1000.0,
		},
		Unofloat: MinMax[float64]{
			Min: 0.0,
			Max: 1.0,
		},
	},
	StringDefaults: StringDefaults{
		Charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		Length: MinMax[uint64]{
			Min: 10,
			Max: 10,
		},
	},
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filename string) (*Config, error) {
	cfg := DefaultConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateConfig checks that config values are valid
func validateConfig(cfg *Config) error {
	if cfg.Int.Min >= cfg.Int.Max {
		return fmt.Errorf("int.min (%v) must be less than int.max (%v)", cfg.Int.Min, cfg.Int.Max)
	}

	if cfg.Uint.Min >= cfg.Uint.Max {
		return fmt.Errorf("uint.min (%v) must be less than uint.max (%v)", cfg.Uint.Min, cfg.Uint.Max)
	}

	if cfg.Float.Min >= cfg.Float.Max {
		return fmt.Errorf("float.min (%v) must be less than float.max (%v)", cfg.Float.Min, cfg.Float.Max)
	}

	if cfg.Unofloat.Min < 0.0 || cfg.Unofloat.Min > 1.0 {
		return fmt.Errorf("unofloat.min (%v) must be between 0.0 and 1.0", cfg.Unofloat.Min)
	}
	if cfg.Unofloat.Max < 0.0 || cfg.Unofloat.Max > 1.0 {
		return fmt.Errorf("unofloat.max (%v) must be between 0.0 and 1.0", cfg.Unofloat.Max)
	}
	if cfg.Unofloat.Min >= cfg.Unofloat.Max {
		return fmt.Errorf("unofloat.min (%v) must be less than unofloat.max (%v)", cfg.Unofloat.Min, cfg.Unofloat.Max)
	}

	if cfg.Charset == "" {
		return fmt.Errorf("charset cannot be empty")
	}

	return nil
}
