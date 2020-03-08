package network

import (
	"net"
	"testing"
)

func TestAllocate(t *testing.T){
	//allocate 192.168.0.0/24
	_, ipnet, _ := net.ParseCIDR("192.168.0.0/24")
	ip, _ := ipAllocator.Allocate(ipnet)
	t.Logf("alloc ip : %v", ip)
}

func TestRelease(t *testing.T){
	ip, ipnet, _ := net.ParseCIDR("192.168.0.0/24")
	ipAllocator.Release(ipnet, &ip)
}