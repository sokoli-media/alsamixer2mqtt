package internal

import (
	"encoding/json"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

func publishCurrentState(config Config, client mqtt.Client) {
	volume, err := getAlsaVolume(config.AlsaDevice, config.AlsaControl)
	if err != nil {
		log.Printf("error while getting alsa volume: %v", err)
		return
	}

	dumpedPayload, err := json.Marshal(volume)
	if err != nil {
		log.Printf("error while sending updated value: %v", err)
		return
	}
	token := client.Publish(config.StateTopic, 0, true, dumpedPayload)
	token.Wait()
	log.Printf("Published volume to topic=%s: %s", config.StateTopic, dumpedPayload)
}

func subscribeToUpdates(config Config, client mqtt.Client) {
	log.Printf("subscribing to topic: %s", config.SetTopic)
	if token := client.Subscribe(config.SetTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		newVolume, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			log.Printf("Invalid volume value: %s", msg.Payload())
			return
		}

		if err = setAlsaVolume(config.AlsaDevice, config.AlsaControl, newVolume); err != nil {
			log.Printf("Failed to set volume: %v", err)
		} else {
			log.Printf("Set volume to: %d%%", newVolume)
			publishCurrentState(config, client)
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}
}

func connectToMqtt(config Config) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttBroker)
	opts.SetClientID(config.MqttClientId)
	if config.MqttUsername != "" && config.MqttPassword != "" {
		opts.SetUsername(config.MqttUsername)
		opts.SetPassword(config.MqttPassword)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	return client
}

func Run(config Config) {
	client := connectToMqtt(config)
	defer client.Disconnect(250)

	subscribeToUpdates(config, client)

	timer := time.NewTicker(5 * time.Second)
	for range timer.C {
		publishCurrentState(config, client)
	}
}
