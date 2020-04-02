package mqttStruct

import (
	"Socker/container"
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

const (
	OnLine = "online"
	OffLine = "offline"
)

type Imqtt interface {
	//connect mqtt
	connect() error
}

type MqttImpl struct {
}

var containerName string

func (m *MqttImpl) Connect(cn string) error {
	containerName = cn
	SetMqttClient(&mqttClient)
	fmt.Println(mqttClient.Server)
	opts := mqtt.NewClientOptions().AddBroker(mqttClient.Server)
	opts.SetCleanSession(true)
	opts.SetClientID(mqttClient.ClientID)
	opts.OnConnect = OnConnect
	opts.OnConnectionLost = OnConnectLost
	//replace {CN} with containerName
	Replace(cn)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)

	//send container online
	if token := client.Publish(GetTopic(SysOnLinePub), 0, false, OnLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}

	var flag = 0
	for {
		if flag == 0 {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				flag = 1
			} else {
				break
			}
		} else if flag == 1 {
			if token := client.Connect(); token.Wait() && token.Error() != nil {

			} else {
				flag = 0
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func OnConnect(client mqtt.Client) {
	fmt.Println("onconnect")

	if token := client.Publish(GetTopic(SysOnLinePub), 0, false, OnLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}

	if token := client.Subscribe(GetTopic(SysDataSub), 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Errorf("client subscribe message Error %v", token.Error())
	}

		//for loop send MSG




}

func OnConnectLost(client mqtt.Client, err error) {
	log.Error("mqtt client lost!")
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	log.Infof("Received message on topic: %s \t Message: %s\n", message.Topic(), message.Payload())
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName);

	filePath := dirURL + "mqttSub"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Create file %s error %v \n", filePath, err)
	}
	file.Write(message.Payload())
}
