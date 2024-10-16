package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"xes.software/rpg/lib/utils"
	"xes.software/rpg/lib/ws"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv has failed to log a .env file: %s", err.Error())
	}
	envVars := utils.GetEnvVars()
	logger := utils.NewLogger(envVars.Environment)
	raidsMap := &ws.RaidsMap{}

	http.HandleFunc("/ws", ws.HandleWs(envVars, logger, raidsMap))

	logger.Log(os.Stdout, "Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe error:", err.Error())
	}
}
