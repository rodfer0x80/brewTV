package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const ALLOWED_MAC_ADDRESSES_PATH = "./allowed_mac_addresses.txt"

func ConfigureAllowedMacAddresses() {
	user_input := 1337
	max_input := 1337
	mac_addresses := ScanLANMacAddresses()
	for user_input != 0 {
		fmt.Println("Configure MAC addresses:")
		for i, mac := range mac_addresses {
			fmt.Printf("%d: %s", i+1, mac.String())
		}
		user_input, err := strconv.Atoi(GetUserInput())
		if err != nil {
			fmt.Println("Invalid user input")
		} else {
			if user_input == 0 {
				fmt.Println("Quitting...")
				return
			}
			if user_input <= max_input {
				err := addToMacAddressAllowlist(mac_addresses[user_input-1].String())
				if err != nil {
					log.Println("[configureAllowedMacAddress]::", err)
					return
				}
			}
		}
	}
}

func addToMacAddressAllowlist(mac_address string) error {
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
