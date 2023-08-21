package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/mdlayher/arp"
)

func GetLANIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("[GetLANIP]::Failed to get local IP:", err)
		os.Exit(1)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	log.Println("[GetLANIP]::Unable to determine LAN IP")
	return ""
}

func getNetworkAdapterInterface() *net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("[getNetworkAdapterInterface]::", err)
		return nil
	}
	for _, intf := range ifaces {
		// Skip loopback and down interfaces
		if intf.Flags&net.FlagLoopback != 0 || intf.Flags&net.FlagUp == 0 {
			continue
		}
		// Skip virtual interfaces and tunnels
		if strings.HasPrefix(intf.Name, "vmnet") || strings.HasPrefix(intf.Name, "tun") {
			continue
		}
		// Call the Addrs function to get the list of addresses
		addrs, err := intf.Addrs()
		if err != nil {
			log.Println("[getNetworkAdapterInterface]::Error getting addresses for", intf.Name, ":", err)
			continue
		}
		// Check if the interface has an IP in the private IP range
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && isPrivateIP(ipnet.IP) {
				return &intf
			}
		}
	}
	log.Println("[getNetworkAdapterInterface]::No LAN-connected network adapter found.")
	return nil
}

func isPrivateIP(ip net.IP) bool {
	privateRanges := []net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}
	for _, pr := range privateRanges {
		if pr.Contains(ip) {
			return true
		}
	}
	return false
}

func ScanLANMacAddresses() []net.HardwareAddr {
	if os.Geteuid() != 0 {
		fmt.Println("[scanLANMacAddresses]::This program requires root privileges to access raw sockets.")
		os.Exit(-1)
	}
	iface := getNetworkAdapterInterface()
	c, err := arp.Dial(iface)
	if err != nil {
		log.Println("[scanLANMacAddresses]::Error creating ARP client:", err)
		return nil
	}
	defer c.Close()
	var macAddresses []net.HardwareAddr
	addr, _ := netip.ParseAddr("192.168.0.1")
	s := addr.AsSlice()
	ip := net.IP(s)
	addrFromIp, _ := netip.AddrFromSlice(ip)
	// Resolve the MAC address using ARP
	hwAddr, err := c.Resolve(addrFromIp)
	if err != nil {
		log.Printf("[scanLANMacAddresses]::Error resolving %s: %v\n", ip, err)
	} else {
		macAddresses = append(macAddresses, hwAddr)
	}
	// Give some time for ARP requests and responses
	time.Sleep(2 * time.Second)
	return macAddresses
}
