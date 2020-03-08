package network

type BridgeNetworkDriver struct {
}

func (d *BridgeNetworkDriver)Name() string {
	return ""
}

func (d *BridgeNetworkDriver)Create(subnet string, name string) (*Network, error){
	return nil, nil
}

func (d *BridgeNetworkDriver)Delete(network Network) error {
	return nil
}

func (d *BridgeNetworkDriver)Connect(network *Network, endpoint *Endpoint) error {
	return nil
}

func (d *BridgeNetworkDriver)Disconnect(network Network, endpoint *Endpoint) error{
	return nil
}