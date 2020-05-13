package network

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ContainerInfo struct {
	Pid				string
	Id				string
	PortMapping		[]string
}
var c ContainerInfo
func GetContainerInfo(Id string, Pid string, PortMapping []string) {
	c.Id = Id
	c.Pid = Pid
	c.PortMapping = PortMapping
}



//连接容器和网桥
func Connect(networkName string) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}

	//allocate IP for container
	ip, err := ipAllocator.Allocate(network.IPRange)
	if err != nil {
		return err
	}

	//set endpoint
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", c.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: c.PortMapping,
	}

	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}

	//set container Net Namespace
	//network device and route
	if err = configEndpointIpAddressAndRoute(ep); err != nil {
		return err
	}

	//set portmapping  e.g. run -p 8080:80
	return configPortMapping(ep)
}

//设置网络Namespace
//网络设备配置和路由
func configEndpointIpAddressAndRoute(ep *Endpoint) error {
	//拿到Veth的另一端
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("fail config endpoint: %v", err)
	}

	//进入 container Net Namespace
	defer enterContainerNetns(&peerLink)()

	//container IPRange
	interfaceIP := *ep.Network.IPRange
	//container IPAddress
	interfaceIP.IP = ep.IPAddress

	//设置容器中Veth的IP
	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v, %s", ep.Network, err)
	}

	//把Veth端点放在容器中
	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}

	//set up "lo"
	if err = setInterfaceUP("lo"); err != nil {
		return err
	}

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	//set route
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Dst:       cidr,
		Gw:        ep.Network.IPRange.IP,
	}

	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}
	return nil
}

//进入container Net Namespace
func enterContainerNetns(enLink *netlink.Link) func() {
	file, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", c.Pid), os.O_RDONLY, 0)
	if err != nil {
		log.Errorf("error get container net namespace, %v", err)
	}

	//file descriptor
	nsFD := file.Fd()
	//lock os thread, so goroutine would not be scheduled to other thread
	//or you would not insure you are in Net Namespace
	runtime.LockOSThread()

	// move Veth to container Net Namespace
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		log.Errorf("error set link netns , %v", err)
	}

	// get current Net Namespace
	origns, err := netns.Get()
	if err != nil {
		log.Errorf("error get current netns, %v", err)
	}

	//enter Net Namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Errorf("error set netns, %v", err)
	}
	return func() {
		//back to orgin Net Namespace
		netns.Set(origns)
		origns.Close()
		runtime.UnlockOSThread()
		file.Close()
	}
}

//set Port Mapping
func configPortMapping(ep *Endpoint) error {
	//range container port mapping list
	//e.g. 8080:80
	for _, pm := range ep.PortMapping {
		log.Infof("port mapping is %s", pm)
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			log.Errorf("port mapping format error, %v", pm)
			continue
		}

		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])

		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		output, err := cmd.Output()
		if err != nil {
			log.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}
