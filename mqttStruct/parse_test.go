package mqttStruct

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func TestSendMsg(t *testing.T) {
	msg := readFile()
	fmt.Println(msg)

	containerName = "socker"
	opts := mqtt.NewClientOptions().AddBroker("tcp://121.40.101.210:1883")
	opts.SetCleanSession(true)
	opts.SetClientID("wang")
	opts.OnConnect = OnConnect2
	opts.OnConnectionLost = OnConnectLost
	opts.SetWill(GetTopic(SysOnLinePub), OffLine, 1, true)

	//replace {CN} with containerName
	Replace(containerName)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)


	sendMessage(client)
}

func OnConnect2(client mqtt.Client) {
	sendMessage(client)
}