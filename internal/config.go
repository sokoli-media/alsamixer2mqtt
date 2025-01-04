package internal

type Config struct {
	MqttBroker   string
	MqttClientId string
	MqttUsername string
	MqttPassword string
	AlsaDevice   string
	AlsaControl  string
	StateTopic   string
	SetTopic     string
}
