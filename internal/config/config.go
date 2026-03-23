package config

import "os"

type Config struct {
	DatabaseURL       string
	FaceitClientID    string
	FaceitAPIKey      string
	FaceitRedirectURI string
	CookieSecret      string
	BaseURL           string
	Port              string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://csinder:csinder@localhost:5432/csinder?sslmode=disable"),
		FaceitClientID:    getEnv("FACEIT_CLIENT_ID", ""),
		FaceitAPIKey:      getEnv("FACEIT_API_KEY", ""),
		FaceitRedirectURI: getEnv("FACEIT_REDIRECT_URI", "http://localhost:8080/auth/callback"),
		CookieSecret:      getEnv("COOKIE_SECRET", "dev-secret-change-me-in-production"),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		Port:              getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
