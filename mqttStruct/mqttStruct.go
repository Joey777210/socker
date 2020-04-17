package mqttStruct

import (
	"Socker/container"
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
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

//mqtt connect
func (m *MqttImpl) Connect(cn string) error {
	containerName = cn
	SetMqttClient(&mqttClient)
	fmt.Println(mqttClient.Server)
	opts := mqtt.NewClientOptions().AddBroker(mqttClient.Server)
	opts.SetCleanSession(true)
	opts.SetClientID(mqttClient.ClientID)
	opts.OnConnect = OnConnect
	opts.OnConnectionLost = OnConnectLost
	opts.SetWill(GetTopic(SysOnLinePub), OffLine, 1, true)

	//replace {CN} with containerName
	Replace(cn)

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
	log.Infoln("onconnect  + " + GetTopic(SysDataSub))

	if token := client.Publish(GetTopic(SysOnLinePub), 0, false, OnLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}
	if token := client.Subscribe(GetTopic(SysDataSub), 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Errorf("client subscribe message Error %v", token.Error())
	}

	//watch file change and send message
	sendMessage(client)
}

func OnConnectLost(client mqtt.Client, err error) {
	log.Error("mqtt client lost!")
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	log.Infof("Received message on topic: %s \t Message: %s\n", message.Topic(), message.Payload())
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)

	fileName := dirURL + "/mqttSub"
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf("Create file %s error %v \n", fileName, err)
	}
	jsonStr := string(message.Payload())
	file.WriteString(jsonStr)
}

func sendMessage(client mqtt.Client) {
	dirURL := fmt.Sprintf("/root/mergeDir/%s", containerName)
	fileName := dirURL + "/mqttPub"
	isExist, _ := PathExists(fileName)
	if !isExist {
		if _, err := os.Create(fileName); err != nil {
			log.Errorf("Create file %s error %v", fileName, err)
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("New watcher error %v", err)
	}
	defer watcher.Close()

	//done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Infoln("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					message := readFile(fileName)
					msgPub(client, message)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infof("Watch file error1 %v", err)
			}
		}
	}()

	err = watcher.Add(fileName)
	if err != nil {
		log.Errorf("Watch file error2 %v", err)
	}
	//<-done
	//循环
	select {}
}

func readFile(fileName string) string {

	var message []byte
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Open file %s error %v \n", fileName, err)
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		message = append(message, buf[:n]...)
	}
	////clear file
	//_ = os.Truncate(fileName, 0)
	//_, _ = file.Seek(0, 0)
	//after test, don't need clear!!
	return string(message)
}

func msgPub(client mqtt.Client, message string) {
	if token := client.Publish(GetTopic(SysDataPub), 0, false, message); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
