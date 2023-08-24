package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"

	"github.com/mdlayher/arp"
)

const (
	PrivateIPRanges = "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"
)

func GetLANIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("[GetLANIP]::Failed to get local IP: %v", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	log.Println("[GetLANIP]::Unable to determine LAN IP")
	return ""
}

func getLANConnectedInterface() *net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("[getLANConnectedInterface]::Error getting network interfaces:", err)
		return nil
	}

	for _, intf := range ifaces {
		if !isLANConnectedInterface(intf) {
			continue
		}
		return &intf
	}

	log.Println("[getLANConnectedInterface]::No LAN-connected network adapter found.")
	return nil
}

func isLANConnectedInterface(intf net.Interface) bool {
	if intf.Flags&net.FlagLoopback != 0 || intf.Flags&net.FlagUp == 0 {
		return false
	}
	if strings.HasPrefix(intf.Name, "vmnet") || strings.HasPrefix(intf.Name, "tun") {
		return false
	}
	addrs, err := intf.Addrs()
	if err != nil {
		log.Printf("[isLANConnectedInterface]::Error getting addresses for %s: %v", intf.Name, err)
		return false
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && isPrivateIP(ipnet.IP) {
			return true
		}
	}
	return false
}

func isPrivateIP(ip net.IP) bool {
	_, privateRange, _ := net.ParseCIDR(PrivateIPRanges)
	return privateRange.Contains(ip)
}

func ResolveMACFromIP(remoteIP string) (string, error) {
	iface := getLANConnectedInterface()
	c, err := arp.Dial(iface)
	if err != nil {
		log.Printf("[ResolveMACFromIP]::Error dialing ARP: %v", err)
		return "", err
	}
	defer c.Close()

	addr, _ := netip.ParseAddr(remoteIP)
	s := addr.AsSlice()
	ip := net.IP(s)
	addrFromIP, _ := netip.AddrFromSlice(ip)

	macAdress, err := c.Resolve(addrFromIP)
	if err != nil {
		log.Printf("[ResolveMACFromIP]::Error resolving IP: %v", err)
		return "", err
	}
	return macAdress.String(), nil
}

func ScanMacAddress(r *http.Request) (string, error) {
	if os.Geteuid() != 0 {
		return "", fmt.Errorf("this program requires root privileges to access raw sockets")
	}
	return ResolveMACFromIP(strings.Split(r.RemoteAddr, ":")[0])
}
