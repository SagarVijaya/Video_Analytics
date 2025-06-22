package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"videoanalytics/apps/models"
)

var configData *models.Config

// LoadGlobalConfig loads the config from a .env file if not already loaded
func LoadGlobalConfig(envPath string) {
	loadDotEnv(envPath) // silently ignored if file missing
	configData = &models.Config{}

	// fill struct from env
	configData.Database.Host = getenv("DB_HOST")
	configData.Database.User = getenv("DB_USER")
	configData.Database.Port = IntValue(getenv("DB_PORT"))
	configData.Database.Pass = getenv("DB_PASS")
	configData.Database.Name = getenv("DB_NAME")

	configData.Server.Port = IntValue(getenv("SERVER_PORT"))

	configData.Redis.Port = getenv("REDIS_PORT")

	configData.Metrics.Port = IntValue(getenv("METRICS_PORT"))

	configData.Mail.Host = getenv("SMT_HOST")
	configData.Mail.Port = getenv("SMT_PORT")
	configData.Mail.From = getenv("FROM_MAIL")
	configData.Mail.To = getenv("TO_MAIL")
	configData.Mail.Pass = getenv("FROM_PASS")

	log.Println("config loaded from environment")
}

// GetConfig returns the already loaded config
func GetConfig() *models.Config {
	if configData == nil {
		log.Fatal("Env file not loaded! Call LoadGlobalConfig first.")
	}
	return configData
}

func IntValue(key string) int {
	lValue, lErr := strconv.Atoi(key)
	if lErr != nil {
		log.Fatalf("%s must be an integer: %v", key, lErr)
	}
	return lValue
}

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return val
}

// Reads KEY=VAL lines and sets them into the process env.
func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Env file not loaded! Call LoadGlobalConfig first.")
		return
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if kv := strings.SplitN(line, "=", 2); len(kv) == 2 {
			os.Setenv(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
		}
	}
}
