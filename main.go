package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID("alsamixer2mqtt")
	if mqttUsername != "" && mqttPassword != "" {
		opts.SetUsername(mqttUsername)
		opts.SetPassword(mqttPassword)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	go func() {
		for {
			volume := getAlsaVolume(alsaDevice, alsaControl)
			if volume >= 0 {
				topic := "homeassistant/sensor/sound_level/state"
				token := client.Publish(topic, 0, true, strconv.Itoa(volume))
				token.Wait()
				log.Printf("Published volume: %d%%", volume)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	topic := "homeassistant/sensor/sound_level/set"
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		newVolume, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			log.Printf("Invalid volume value: %s", msg.Payload())
			return
		}
		if err = setAlsaVolume(alsaDevice, alsaControl, newVolume); err != nil {
			log.Printf("Set volume to: %d%%", newVolume)
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	select {}
}

func getAlsaVolume(device, control string) int {
	out, err := exec.Command("amixer", "-D", device, "get", control).Output()
	if err != nil {
		log.Printf("Failed to get volume: %v", err)
		return -1
	}

	// example output: "[42%]" -> extract "42"
	var volume int
	_, err = fmt.Sscanf(string(out), "%*[^[][%d%%]", &volume)
	if err != nil {
		log.Printf("Failed to parse volume: %v", err)
		return -1
	}

	return volume
}

func setAlsaVolume(device, control string, volume int) error {
	cmd := exec.Command("amixer", "-D", device, "set", control, fmt.Sprintf("%d%%", volume))
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to set volume: %v", err)
		return err
	}
	return nil
}
