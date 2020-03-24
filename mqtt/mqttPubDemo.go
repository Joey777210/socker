package main

import (
"crypto/tls"
log "github.com/Sirupsen/logrus"
"github.com/eclipse/paho.mqtt.golang"
)

const (
	////server = "127.0.0.1:1883"
	//server = "tcp://121.40.101.210:1883"
	//clientID = "zhang"
	//username = "zhang"
	//password = "123"
	//topic = "Hello"
	message = "World!!"
)

func main(){
	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetCleanSession(true)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	tlsConfig := &tls.Config{InsecureSkipVerify:true, ClientAuth:tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)

	if token :=client.Connect();token.Wait() &&token.Error() != nil {
		log.Errorf("client connect error %v\n", token.Error())
		return
	}

	log.Infof("Connect to server %s", server)

	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}

}

