package parse

import (
	"Socker/mqttStruct"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"os"
)

func setMqttClient() {
	fileName := "sk_mqtt.conf"
	fileURL := "../" + fileName
	file, err := os.Open(fileURL)
	if err != nil {
		log.Errorf("Open config file error %v", err)
		return
	}
	defer file.Close()
	buf := make([]byte, 1024)
	var mqttJson []byte
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		mqttJson = append(mqttJson, buf[:n]...)
	}
	var client mqttStruct.Client
	if err := json.Unmarshal(mqttJson, client); err != nil {
		log.Errorf("Mqtt json unmarshal error %v", err)
	}
}
