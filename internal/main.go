package internal

import (
	"encoding/json"
	"fmt"
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

func subscribeToUpdates(config Config, client mqtt.Client) error {
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
		return token.Error()
	}

	return nil
}

func connectToMqtt(config Config) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttBroker)
	opts.SetClientID(config.MqttClientId)
	if config.MqttUsername != "" && config.MqttPassword != "" {
		opts.SetUsername(config.MqttUsername)
		opts.SetPassword(config.MqttPassword)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	return client, nil
}

func connectAndConsume(config Config) {
	client, err := connectToMqtt(config)
	if err != nil {
		log.Printf("couldn't connect to mqtt server: %v", err)
		time.Sleep(5 * time.Second)
		return
	}
	defer client.Disconnect(250)

	err = subscribeToUpdates(config, client)
	if err != nil {
		log.Printf("couldn't start consuming from mqtt: %v", err)
		time.Sleep(5 * time.Second)
		return
	}

	timer := time.NewTicker(30 * time.Second)
	for range timer.C {
		publishCurrentState(config, client)
	}
}

func Run(config Config) {
	for {
		connectAndConsume(config)
	}
}
