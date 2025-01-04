package internal

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

func Run(mqttBroker string, mqttUsername string, mqttPassword string, alsaDevice string, alsaControl string, sensorName string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID(sensorName)
	if mqttUsername != "" && mqttPassword != "" {
		opts.SetUsername(mqttUsername)
		opts.SetPassword(mqttPassword)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	// Publish the discovery message
	sensorConfig := fmt.Sprintf(`{
		"name": "%s",
		"object_id": "%s",
		"state_topic": "homeassistant/sensor/%s/state",
		"command_topic": "homeassistant/sensor/%s/set",
		"unit_of_measurement": "dB",
		"device_class": "sound_pressure",
		"value_template": "{{ value_json.value }}",
		"device": {
			"name": "AlsaMixer for %s",
			"identifiers": ["%s"],
			"model": "AlsaMixer",
			"manufacturer": "AlsaMixer"
		}
	}`, sensorName, sensorName, sensorName, sensorName, sensorName)
	topic := fmt.Sprintf("homeassistant/sensor/%s/config", sensorName)
	token := client.Publish(topic, 0, true, sensorConfig)
	log.Printf("Sent autodiscovery message to topic=%s payload=%s", topic, sensorConfig)
	token.Wait()

	go func() {
		for {
			volume, err := getAlsaVolume(alsaDevice, alsaControl)
			if err != nil {
				log.Printf("error while getting alsa volume: %v", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			topic = fmt.Sprintf("homeassistant/sensor/%s/state", sensorName)
			payload := map[string]float64{
				"value": volume,
			}
			dumpedPayload, err := json.Marshal(payload)
			if err != nil {
				log.Printf("error while sending updated value: %v", err)
				continue
			}
			token = client.Publish(topic, 0, true, dumpedPayload)
			token.Wait()
			log.Printf("Published volume to topic=%s: %s", topic, dumpedPayload)

			time.Sleep(500 * time.Millisecond)
		}
	}()

	topic = fmt.Sprintf("homeassistant/sensor/%s/set", sensorName)
	log.Printf("subscribing to topic: %s", topic)
	if token = client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		newVolume, err := strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			log.Printf("Invalid volume value: %s", msg.Payload())
			return
		}
		if err = setAlsaVolume(alsaDevice, alsaControl, newVolume); err != nil {
			log.Printf("Failed to set volume: %v", err)
		} else {
			log.Printf("Set volume to: %f%%", newVolume)
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	select {}
}
