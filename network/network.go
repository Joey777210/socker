package network

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"
)

var (
	defaultNetworkPath = "/var/run/socker/network/network/"
	drivers = map[string]NetworkDriver{}
	networks = map[string]*Network{}
)

//struct of network contains muti containers
type Network struct {
	Name string			//net name
	IPRange *net.IPNet	//address
	Driver string		//net driver
}

//veth. link container and network
type Endpoint struct {
	ID string 			`json:"id"`
	Device netlink.Veth `json:"dev"`
	IPAddress net.IP	`json:"ip"`
	MacAddress net.HardwareAddr		`json:"mac"`
	PortMapping []string `json:"protmapping"`
	Network	*Network
}

type NetworkDriver interface{
	//driver name
	Name() string

	Create(subnet string, name string) (*Network, error)

	Delete(network Network) error

	Connect(network *Network, endpoint *Endpoint) error

	Disconnect(network Network, endpoint *Endpoint) error
}

/**
	//create a network
 */
func CreateNetwork(driver string, subnet string, name string) error {
	//parse subnet string to net.IPNet object
	_, cidr, _ := net.ParseCIDR(subnet)
	//IPAM allocate gateway IP. the first IP address as gateway IP
	gatewayIp, err := ipAllocator.Allocate(cidr)

	if err != nil {
		return err
	}
	cidr.IP = gatewayIp

	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		log.Errorf("Create driver %s error %v", driver, err)
		return err
	}

	//save nw
	return nw.dump(defaultNetworkPath)
}


func (nw *Network) dump(dumpPath string) error {
	if _,err := os.Stat (dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		}else {
			return err
		}
	}

	nwPath := path.Join(dumpPath, nw.Name)
	//clear all content | write only | create if not exit
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("create nwFile %s error %v",nwPath, err)
		return err
	}
	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		log.Errorf("json marshal nw %v error: %s",nw, err)
	}

	_, err = nwFile.Write(nwJson)
	if err != nil {
		log.Errorf("write nwjson into nwfile %v error %v", nwFile, err)
		return err
	}
	return nil
}


func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		log.Errorf("open nwConfigFIle %s error %v", dumpPath, err)
		return err
	}

	//read config json
	nwJson := make([]byte, 1024*1024)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		log.Errorf("Unmarshal nwJson %s error %v", nwJson, err)
		return err
	}
	return nil
}

func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}else {
		os.Remove(path.Join(dumpPath, nw.Name))
	}
	return nil
}
//init driver and network, get all networks created
func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNetworkPath, 0644)
		} else {
			return err
		}
	}

	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir(){
			return nil
		}

		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name:    nwName,
		}

		if err := nw.load(nwPath); err != nil {
			log.Errorf("error load network %s", err)
		}

		//add network config into networks
		networks[nwName] = nw
		return nil
	})
	return nil
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIPrange\tDriver\n")

	//range network
	for _, nw := range networks{
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			nw.Name,
			nw.IPRange.String(),
			nw.Driver,
			)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No such Network: %s", networkName)
	}

	//release netwok gateway IP
	if err := ipAllocator.Release(nw.IPRange, &nw.IPRange.IP); err != nil {
		return fmt.Errorf("Remove Network gateway ip %s error %v", nw.Name, err)
	}

	//call driver delete device and config
	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("Remove Network gateway ip %s error %v", nw.Name, err)
	}

	return nw.remove(defaultNetworkPath)
}