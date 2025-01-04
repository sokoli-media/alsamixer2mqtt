package internal

import (
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

func Run(mqttBroker string, mqttUsername string, mqttPassword string, alsaDevice string, alsaControl string) {
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
