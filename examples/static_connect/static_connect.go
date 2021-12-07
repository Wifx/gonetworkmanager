package main

import (
	"fmt"
	"os"

	"github.com/Wifx/gonetworkmanager"
	"github.com/google/uuid"
)

const (
	ethernetType                 = "802-3-ethernet"
	ethernetSection              = "802-3-ethernet"
	ethernetSectionAutoNegotiate = "auto-negotiate"
	connectionSection            = "connection"
	connectionSectionID          = "id"
	connectionSectionType        = "type"
	connectionSectionUUID        = "uuid"
	connectionSectionIfaceName   = "interface-name"
	connectionSectionAutoconnect = "autoconnect"
	ip4Section                   = "ipv4"
	ip4SectionAddressData        = "address-data"
	ip4SectionAddresses          = "addresses"
	ip4SectionAddress            = "address"
	ip4SectionPrefix             = "prefix"
	ip4SectionMethod             = "method"
	ip4SectionGateway            = "gateway"
	ip4SectionNeverDefault       = "never-default"
	ip6Section                   = "ipv6"
	ip6SectionMethod             = "method"

	connectionID                   = "My Connection"
	interfaceName                  = "eth1"
	desiredIPAddress               = "192.168.1.1"
	desiredGatewayAddress          = "192.168.1.1"
	desiredIPAddressNumerical      = 16885952
	desiredIPPrefix                = 24
	desiredGatewayAddressNumerical = 16885952

	// Allows for static ip
	desiredIP4Method = "manual"

	// Would like this to be "disabled" however not supported
	// in the current network manager stack
	desiredIP6Method = "ignore"
)

func printVersion() error {
	/* Create new instance of gonetworkmanager */
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return err
	}

	// Don't really need the network manager object per se
	// however knowing the version isn't bad
	var nmVersion string
	nmVersion, err = nm.GetPropertyVersion()
	if err != nil {
		return err
	}

	fmt.Println("Network Manager Version: " + nmVersion)
	return nil
}

func checkForExistingConnection() (bool, error) {
	// See if our connection already exists
	settings, err := gonetworkmanager.NewSettings()

	if err != nil {
		return false, err
	}

	currentConnections, err := settings.ListConnections()
	if err != nil {
		return false, err
	}

	for _, v := range currentConnections {
		connectionSettings, settingsError := v.GetSettings()
		if settingsError != nil {
			fmt.Println("settings error, continuing")
			continue
		}
		currentConnectionSection := connectionSettings[connectionSection]
		if currentConnectionSection[connectionSectionID] == connectionID {
			return true, nil
		}
	}
	return false, nil
}

func createNewConnection() error {

	connection := make(map[string]map[string]interface{})
	connection[ethernetSection] = make(map[string]interface{})
	connection[ethernetSection][ethernetSectionAutoNegotiate] = false
	connection[connectionSection] = make(map[string]interface{})
	connection[connectionSection][connectionSectionID] = connectionID
	connection[connectionSection][connectionSectionType] = ethernetType
	connectionUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	connection[connectionSection][connectionSectionUUID] = connectionUUID.String()
	connection[connectionSection][connectionSectionIfaceName] = interfaceName
	connection[connectionSection][connectionSectionAutoconnect] = true
	connection[ip4Section] = make(map[string]interface{})

	addressData := make([]map[string]interface{}, 1)

	addressData[0] = make(map[string]interface{})
	addressData[0][ip4SectionAddress] = desiredIPAddress
	addressData[0][ip4SectionPrefix] = desiredIPPrefix

	connection[ip4Section][ip4SectionAddressData] = addressData

	// order defined by network manager
	addresses := make([]uint32, 3)
	addresses[0] = desiredIPAddressNumerical
	addresses[1] = desiredIPPrefix
	addresses[2] = desiredGatewayAddressNumerical

	addressArray := make([][]uint32, 1)
	addressArray[0] = addresses
	connection[ip4Section][ip4SectionAddresses] = addressArray

	connection[ip4Section][ip4SectionGateway] = desiredGatewayAddress
	connection[ip4Section][ip4SectionMethod] = desiredIP4Method
	connection[ip4Section][ip4SectionNeverDefault] = true

	connection[ip6Section] = make(map[string]interface{})
	connection[ip6Section][ip6SectionMethod] = desiredIP6Method

	settings, err := gonetworkmanager.NewSettings()

	if err != nil {
		return err
	}

	_, err = settings.AddConnection(connection)

	if err != nil {
		return err
	}
	return nil
}

func main() {

	// show the version
	if printVersion() != nil {
		fmt.Println("failed to find version.  Is NetworkManager running?")
		os.Exit(1)
	}

	// See if our connection already exists
	doesExist, err := checkForExistingConnection()

	// if an error then we are done.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// if the connection already exists we are done.
	if doesExist == true {
		fmt.Println("connection already exists, nothing to do.")
		os.Exit(0)
	}

	// create the new connection
	err = createNewConnection()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print("added " + connectionID + " to the system.")

	os.Exit(0)
}
