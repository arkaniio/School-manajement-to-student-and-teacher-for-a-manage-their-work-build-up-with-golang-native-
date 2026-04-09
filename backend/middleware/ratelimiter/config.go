package ratelimiter

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Rate        int
	Window      time.Duration
	Burst       int
	IPWhiteList []string
}

func DefaultConfig() Config {
	return Config{
		Rate:        60,
		Window:      1 * time.Minute,
		Burst:       10,
		IPWhiteList: []string{"127.0.0.1", "::1"},
	}
}

func StrictConfig() Config {
	return Config{
		Rate:        10,
		Window:      1 * time.Minute,
		Burst:       3,
		IPWhiteList: []string{"127.0.0.1", "::1"},
	}
}

func LoadFromEnvDefault() Config {

	cfg := DefaultConfig()

	if val := os.Getenv("RATE_LIMIT_RATE"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Rate = convert
		}
	}

	if val := os.Getenv("RATE_LIMIT_WINDOW"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Window = time.Duration(convert)
		}
	}

	if val := os.Getenv("RATE_LIMIT_BURST"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Burst = convert
		}
	}

	if val := os.Getenv("RATE_LIMIT_WHITELIST"); val != "" {
		ips := strings.Split(val, ",")
		trimned := make([]string, 0, len(ips))
		for _, ip := range ips {
			if trimned_id := strings.TrimSpace(ip); trimned_id != "" {
				trimned = append(trimned, trimned_id)
			}
		}
		cfg.IPWhiteList = trimned
	}

	return cfg
}

func LoadFromEnvStrict() Config {

	cfg := StrictConfig()

	if val := os.Getenv("RATE_LIMIT_STRICT_RATE"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Rate = convert
		}
	}

	if val := os.Getenv("RATE_LIMIT_STRICT_WINDOW"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Window = time.Duration(convert)
		}
	}

	if val := os.Getenv("RATE_LIMIT_STRICT_BURST"); val != "" {
		if convert, err := strconv.Atoi(val); err == nil {
			cfg.Burst = convert
		}
	}

	if val := os.Getenv("RATE_LIMIT_STRICT_WHITELIST"); val != "" {
		ips := strings.Split(val, ",")
		trimned := make([]string, 0, len(ips))
		for _, ip := range ips {
			if trimned_id := strings.TrimSpace(ip); trimned_id != "" {
				trimned = append(trimned, trimned_id)
			}
		}
		cfg.IPWhiteList = trimned
	}

	return cfg

}
