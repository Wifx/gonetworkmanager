package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus/v5"
)

const (
	DeviceDummyInterface = DeviceInterface + ".Dummy"

	/* Properties */
)

type DeviceDummy interface {
	Device
}

func NewDeviceDummy(objectPath dbus.ObjectPath) (DeviceDummy, error) {
	var d deviceDummy
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type deviceDummy struct {
	device
}

func (d *deviceDummy) MarshalJSON() ([]byte, error) {
	m, err := d.device.marshalMap()
	if err != nil {
		return nil, err
	}

	m["HwAddress"], _ = d.GetPropertyHwAddress()
	return json.Marshal(m)
}
