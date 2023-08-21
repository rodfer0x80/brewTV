package main

import (
	"log"
	"net"
)

const ALLOWED_MAC_ADDRESSES_PATH = "./allowed_mac_addresses.txt"

func AddToMacAddressAllowlist(mac_address string) error {
	return AppendToFile(ALLOWED_MAC_ADDRESSES_PATH, mac_address)
}

func ReadAllowedMacAddresses(filename string) ([]net.HardwareAddr, error) {
	var mac_addresses []net.HardwareAddr
	str_mac_addresses, err := ReadlinesFromFile(ALLOWED_MAC_ADDRESSES_PATH)
	if err != nil {
		log.Printf("[ReadAllowedMacAddresses]::Error opening file: %s\n", err)
	}
	for _, str_mac_address := range str_mac_addresses {
		mac_address, err := net.ParseMAC(str_mac_address)
		if err != nil {
			log.Printf("[ReadAllowedMacAddresses]::Invalid MAC address: %s\n", mac_address)
			continue
		}
		mac_addresses = append(mac_addresses, mac_address)
	}
	return mac_addresses, nil
}
