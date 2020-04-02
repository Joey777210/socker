package mqttStruct

import (
	"testing"
)

func TestClient(t *testing.T) {
	SetMqttClient(mqttClient)
	t.Logf("parse success")
}

func TestMqtt(t *testing.T) {
	var m MqttImpl
	err := m.Connect("socker")
	t.Logf("%v", err)
}