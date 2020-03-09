package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
	"strings"
)

type BridgeNetworkDriver struct {
}

func (d *BridgeNetworkDriver)Name() string {
	return "bridge"
}

func (d *BridgeNetworkDriver)Create(subnet string, name string) (*Network, error){
	ip, ipRange, err := net.ParseCIDR(subnet)
	ipRange.IP = ip
	if err != nil {
		log.Errorf("Parse subnet error %v", err)
	}
	net := &Network{
		Name:name,
		IPRange:ipRange,
		Driver:d.Name(),
	}

	err = d.initBridge(net)
	if err != nil {
		log.Errorf("Init bridge error %v")
	}

	return net, err

}

func (d *BridgeNetworkDriver)Delete(network Network) error {
	bridgeName := network.Name
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	return netlink.LinkDel(br)
}

func (d *BridgeNetworkDriver)Connect(network *Network, endpoint *Endpoint) error {
	return nil
}

func (d *BridgeNetworkDriver)Disconnect(network Network, endpoint *Endpoint) error{
	return nil
}

func (d *BridgeNetworkDriver)initBridge(net *Network) error {
	bridgeName := net.Name
	//1, create bridge interface
	if err := createBridgeInterface(bridgeName); err != nil {
		log.Errorf("Create Bridge Interface %s Error %v",bridgeName, err)
		return err
	}
	//2, set bridge IP and route
	gatewayIP := net.IPRange
	gatewayIP.IP = net.IPRange.IP
	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		log.Errorf("Assign address on bridge %s Error %v", bridgeName, err)
		return err
	}

	//3,start bridge
	if err := setInterfaceUP(bridgeName); err != nil {
		log.Errorf("Set bridge %s up error %v", bridgeName, err)
		return err
	}

	//4, set iptable SNAT
	if err := setupIPTables(bridgeName, net.IPRange); err != nil {
		log.Errorf("Setting iptables for %s error %v", bridgeName, err)
		return err
	}
	return nil
}

//create Linux bridge
//ip link add xxxx
func createBridgeInterface(bridgeName string) error {
	//1, check if the name exists
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface"){
		return err
	}

	//2, init netlink
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName
	
	//3, use netlink create bridge
	br := &netlink.Bridge{
		LinkAttrs:         la,
	}
	if err := netlink.LinkAdd(br); err != nil {
		log.Errorf("Bridge %s create fail error %v", bridgeName, err)
	}
	return nil
}

//set IP for bridge
//param name: bridge name
//param ip: gateway IP and Mask   192.168.0.1/24
func setInterfaceIP(name string, ip string) error {
	bridgeInterface, err := netlink.LinkByName(name)
	if err != nil {
		log.Errorf("Get interface %s error %v", name, err)
		return err
	}
	ipNet, err := netlink.ParseIPNet(ip)
	if err != nil {
		return err
	}
	addr := &netlink.Addr{
		IPNet:       ipNet,
		Label:       "",
		Flags:       0,
		Scope:       0,
	}
	return netlink.AddrAdd(bridgeInterface, addr)

}

//set interface status UP
//ip link set xxxx up
func setInterfaceUP(interfaceName string) error {
	bridgeInterface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		log.Errorf("Get bridge %s error %v", interfaceName, err)
		return err
	}

	if err := netlink.LinkSetUp(bridgeInterface); err != nil {
		log.Errorf("Set bridge interface %s up error %v", interfaceName, err)
		return err
	}
	return nil
}

//iptable SNAT
func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	//myIptables, err := iptables.New()
	//if err != nil {
	//	log.Errorf("New iptables of %s  error %v", bridgeName, err)
	//	return err
	//}
	//table := "nat"
	//cmd := "POSTROUTING -s "+ subnet.String() + " ! -o "+bridgeName+" -j MASQUERADE"
	////err = myIptables.Append(table, bridgeName, rulespec)
	//err = myIptables.Append(table, cmd)
	//if err != nil {
	//	log.Errorf("run Iptable cmd of  %s error %v", bridgeName, err)
	//	return err
	//}
	//return nil



	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}

	return nil
}

