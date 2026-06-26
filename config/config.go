package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	Port             string
	AdminUsername    string
	AdminPassword    string
	DatabaseURL      string
	BaseURL          string
	SupabaseURL      string
	SupabaseKey      string
	SupabaseBucket   string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "faxman2024"),
		DatabaseURL:    buildDSN(),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		SupabaseURL:    getEnv("SUPABASE_URL", ""),
		SupabaseKey:    getEnv("SUPABASE_KEY", ""),
		SupabaseBucket: getEnv("SUPABASE_BUCKET", "faxman-travels"),
	}
}

func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "")
	name := getEnv("DB_NAME", "postgres")
	ssl  := getEnv("DB_SSLMODE", "disable")

	// URL format lets us pass search_path cleanly via query param so every
	// connection in the pool automatically targets the al_haram schema.
	q := url.Values{}
	q.Set("sslmode", ssl)
	q.Set("search_path", "al_haram")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		url.PathEscape(user), url.PathEscape(pass), host, port, name, q.Encode())
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
