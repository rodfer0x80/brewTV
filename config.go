package main

import (
	"log"
	"net"
)

const ALLOWED_MAC_ADDRESSES_PATH = "./allowed_mac_addresses.txt"

func AddToMacAddressAllowlist(macAddress string) error {
	if ip := net.ParseIP(macAddress); ip != nil {
		macFromIP, err := ResolveMACFromIP(ip.String())
		if err != nil {
			log.Printf("[AddToMacAddressAllowlist]::Error resolving MAC address from IP: %v\n", err)
			return err
		}

		macAddress = macFromIP
	}

	existingMACAddresses, err := ReadAllowedMacAddresses(ALLOWED_MAC_ADDRESSES_PATH)
	if err != nil {
		log.Printf("[AddToMacAddressAllowlist]::Error reading allowed MAC addresses: %v\n", err)
		return err
	}

	for _, existingMAC := range existingMACAddresses {
		if existingMAC.String() == macAddress {
			log.Printf("MAC address already in allowed list %s\n", macAddress)
			return nil
		}
	}

	if err := AppendToFile(ALLOWED_MAC_ADDRESSES_PATH, macAddress); err != nil {
		log.Printf("[AddToMacAddressAllowlist]::Error appending MAC address: %v\n", err)
		return err
	}

	return nil
}

func ReadAllowedMacAddresses(filename string) ([]net.HardwareAddr, error) {
	var macAddresses []net.HardwareAddr

	strMacAddresses, err := ReadlinesFromFile(filename)
	if err != nil {
		log.Printf("[ReadAllowedMacAddresses]::Error opening file: %v\n", err)
		return nil, err
	}

	for _, strMacAddress := range strMacAddresses {
		macAddress, err := net.ParseMAC(strMacAddress)
		if err != nil {
			log.Printf("[ReadAllowedMacAddresses]::Invalid MAC address: %s\n", strMacAddress)
			continue
		}
		macAddresses = append(macAddresses, macAddress)
	}

	return macAddresses, nil
}
