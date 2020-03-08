package network

import (
	"Socker/container"
	"fmt"
)

func Connect(networkName string, cinfo *container.ContainerInfo) error {
	network, ok := networks[networkName]

	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}

	ip, err := ipAllocator.Allocate(network.IPRange)
	if err != nil {
		return err
	}

	//create nw endpoint
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", cinfo.Id, networkName),
		IPAddress:      ip,
		Network:  network,
		PortMapping: cinfo.PortMapping,
	}

	//connect ep
	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}

	//set container IP and route
	//if err = configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}

	//set container portmapping
	//return configPortMapping(ep, cinfos)
	
