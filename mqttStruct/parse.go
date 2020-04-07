package mqttStruct

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"os"
)
var (
	mqttFilePath =  "sk_mqtt.conf"
)

//parse mqttConfig and json unmarshal
func SetMqttClient(v interface{}) {
	file, err := os.Open(mqttFilePath)
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
	if err := json.Unmarshal(mqttJson, &v); err != nil {
		log.Errorf("Mqtt json unmarshal error %v", err)
	}
}
