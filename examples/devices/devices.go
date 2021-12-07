package main

import (
	"fmt"
	"github.com/Wifx/gonetworkmanager"
	"os"
)

func main() {

	/* Create new instance of gonetworkmanager */
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	/* Get devices */
	devices, err := nm.GetPropertyAllDevices()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	/* Show each device path and interface name */
	for _, device := range devices {

		deviceInterface, err := device.GetPropertyInterface()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(deviceInterface + " - " + string(device.GetPath()))
	}

	os.Exit(0)
}
