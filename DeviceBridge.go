package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus/v5"
)

const (
	DeviceBridgeInterface = DeviceInterface + ".Bridge"

	// Properties
	DeviceBridgePropertySlaves  = DeviceBridgeInterface + ".Slaves"  // readable   as
	DeviceBridgePropertyCarrier = DeviceBridgeInterface + ".Carrier" // readable   b
)

type DeviceBridge interface {
	Device

	// GetPropertySlaves Array of paths of enslaved and active devices.
	// DEPRECATED. Use the "Ports" property in "org.freedesktop.NetworkManager.Device" instead which exists since version NetworkManager 1.34.0.
	GetPropertySlaves() ([]Device, error)

	// GetPropertyCarrier Indicates whether the physical carrier is found (e.g. whether a cable is plugged in or not).
	// DEPRECATED: check for the "carrier" flag in the "InterfaceFlags" property on the "org.freedesktop.NetworkManager.Device" interface.
	GetPropertyCarrier() (bool, error)
}

func NewDeviceBridge(objectPath dbus.ObjectPath) (DeviceBridge, error) {
	var d deviceBridge
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type deviceBridge struct {
	device
}

func (d *deviceBridge) GetPropertySlaves() ([]Device, error) {
	paths, err := d.getSliceObjectProperty(DeviceBridgePropertySlaves)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, len(paths))
	for i, path := range paths {
		devices[i], err = DeviceFactory(path)
		if err != nil {
			return nil, err
		}
	}
	return devices, nil
}

func (d *deviceBridge) GetPropertyCarrier() (bool, error) {
	return d.getBoolProperty(DeviceBridgePropertyCarrier)
}

func (d *deviceBridge) MarshalJSON() ([]byte, error) {
	m, err := d.device.marshalMap()
	if err != nil {
		return nil, err
	}

	m["HwAddress"], _ = d.GetPropertyHwAddress()
	m["Slaves"], _ = d.GetPropertySlaves()
	m["Carrier"], _ = d.GetPropertyCarrier()
	return json.Marshal(m)
}
