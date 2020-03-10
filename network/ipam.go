package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
	"strings"
)

const ipamDefaultAllocatorPath = "/var/run/socker/network/ipam/subnet.json"

type IPAM struct {
	//locate of allocate infomation file
	SubnetAllocatorPath string
	//for bitmap. key is network segment, value is array of allocated segment
	//string as bitmap array
	Subnets *map[string]string
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM)Allocate(subnet *net.IPNet) (ip net.IP, err error){
	ipam.Subnets = &map[string]string{}
	err = ipam.load()
	if err != nil {
		log.Errorf("load ipam subnetConfigFile error %v", err)
		return nil, err
	}

	//one is size of '1's in subnet mask
	//size is 24
	one, size := subnet.Mask.Size()

	//if the array has not been init, init it.
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist{
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1 << uint8(size - one))
	}

	for c:= range (*ipam.Subnets)[subnet.String()] {
		//find "0"
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			//set '0' as '1' means allocate this IP
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			ip = subnet.IP

			ip = ip.To4()
			for t := uint(4); t > 0; t -=1 {
				[]byte(ip)[4-t] += uint8(c >> ((t-1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	ipam.dump()

	return
}

func (ipam *IPAM)Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())
	err := ipam.load()
	if err != nil{
		log.Errorf("Error dump allocation info, %v", err)
		return err
	}

	//figure IP index in the bitmap array
	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1] - subnet.IP[t-1]) << ((4-t) * 8)
	}

	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.dump()
	return nil
}

//load subnet allocate info
func (ipam *IPAM) load() error {
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		log.Errorf("open subnetConfigFile %s error %v", ipam.SubnetAllocatorPath, err)
		return err
	}

	subnetJson := make([]byte, 1024*1024)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		log.Errorf("read subnetConfigFile %s error %v", ipam.SubnetAllocatorPath, err)
	}

	err = json.Unmarshal([]byte(subnetJson[:n]), ipam)
	if err != nil {
		log.Errorf("Unmarshal subnetJson error %s", err)
		return err
	}
	return nil
}

//save subnet allocate info
func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0644)
		} else {
			return err
		}
	}

	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}

	subnetJson,err := json.Marshal(ipam)
	if err != nil {
		log.Errorf("json marshal IPAM err %v", err)
		return err
	}

	_, err = subnetConfigFile.Write(subnetJson)
	if err != nil {
		log.Errorf("dump subnetJson error %v", err)
		return err
	}

	return nil
}