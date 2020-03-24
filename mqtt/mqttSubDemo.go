
package main

import (
"crypto/tls"
"fmt"
log "github.com/Sirupsen/logrus"
"github.com/eclipse/paho.mqtt.golang"
"os"
)

const (
	server = "tcp://121.40.101.210:1883"
	//server = "127.0.0.1:1883"
	clientID = "qi"
	username = "zhang"
	password = "123"
	topic = "Hello"
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}

func main(){
	c := make(chan os.Signal, 1)

	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetUsername(username)
	opts.SetPassword(password)

	tlsConfig := &tls.Config{InsecureSkipVerify:true, ClientAuth:tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)



	opts.OnConnect = func(c mqtt.Client){
		if token := c.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			log.Errorf("client subscribe message Error %v", token.Error())
		}
	}

	client := mqtt.NewClient(opts)

	if token :=client.Connect();token.Wait() &&token.Error() != nil {
		log.Errorf("client connect error %v", token.Error())
		return
	}

	log.Infof("Connect to server %s", server)

	<- c

}