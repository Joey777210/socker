package mqttStruct

type Client struct {
	Server string
	ClientID string
	ContainerName string
}

var mqttClient Client