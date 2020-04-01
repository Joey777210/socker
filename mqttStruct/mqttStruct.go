package mqttStruct

import (
	"crypto/tls"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

const (
	//server = "tcp://121.40.101.210:1883"
	server   = "127.0.0.1:1883"
	clientID = "qi"
	username = "zhang"
	password = "123"
	topic    = "Hello"
	message  = "World!!"
)

//Topics
//shifou shangxian ?
//shangchuanxiaoxi
//jieshouxiaoxi
//

type Imqtt interface {
	//connect mqtt
	connect() error
}

type MqttImpl struct {
}

func (m *MqttImpl) Connect() error {

	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetCleanSession(true)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnect = OnConnect
	opts.OnConnectionLost = OnConnectLost

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)

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
	for true {
		if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
			log.Errorf("client publish error %v\n", token.Error())
		}


		if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			log.Errorf("client subscribe message Error %v", token.Error())
		}
	}


}

func OnConnectLost(client mqtt.Client, err error) {
	log.Error("mqtt client lost!")
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	log.Infof("Received message on topic: %s \t Message: %s\n", message.Topic(), message.Payload())
	//filePath := fmt.Sprintf(container.DefaultInfoLocation, c ;
}
