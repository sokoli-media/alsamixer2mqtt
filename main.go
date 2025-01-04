package main

import (
	"alsamixer2mqtt/internal"
	"os"
)

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func main() {
	mqttBroker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	mqttUsername := getEnv("MQTT_USERNAME", "")
	mqttPassword := getEnv("MQTT_PASSWORD", "")
	alsaDevice := getEnv("ALSA_DEVICE", "default")
	alsaControl := getEnv("ALSA_CONTROL", "Master")
	sensorName := getEnv("SENSOR_NAME", "wyoming_satellite_sound_level") // used also as part of the topic name

	internal.Run(mqttBroker, mqttUsername, mqttPassword, alsaDevice, alsaControl, sensorName)
}
