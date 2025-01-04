package main

import (
	"alsamixer2mqtt/internal"
	"log"
	"os"
)

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("couldn't load environmental variable: %s", key)
	}
	return value
}

func main() {
	config := internal.Config{
		MqttBroker:   getEnv("MQTT_BROKER"),
		MqttClientId: getEnv("MQTT_CLIENT_ID"),
		MqttUsername: getEnv("MQTT_USERNAME"),
		MqttPassword: getEnv("MQTT_PASSWORD"),
		AlsaDevice:   getEnv("ALSA_DEVICE"),
		AlsaControl:  getEnv("ALSA_CONTROL"),
		StateTopic:   getEnv("STATE_TOPIC"),
		SetTopic:     getEnv("SET_TOPIC"),
	}

	internal.Run(config)
}
