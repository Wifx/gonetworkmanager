package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus/v5"
)

const (
	DeviceGenericInterface = DeviceInterface + ".Generic"

	// Properties
	DeviceGenericPropertyTypeDescription = DeviceGenericInterface + ".TypeDescription" // readable   s
)

type DeviceGeneric interface {
	Device

	// GetPropertyTypeDescription A (non-localized) description of the interface type, if known.
	GetPropertyTypeDescription() (string, error)
}

func NewDeviceGeneric(objectPath dbus.ObjectPath) (DeviceGeneric, error) {
	var d deviceGeneric
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type deviceGeneric struct {
	device
}

func (d *deviceGeneric) GetPropertyTypeDescription() (string, error) {
	return d.getStringProperty(DeviceGenericPropertyTypeDescription)
}

func (d *deviceGeneric) MarshalJSON() ([]byte, error) {
	m, err := d.device.marshalMap()
	if err != nil {
		return nil, err
	}

	m["HwAddress"], _ = d.GetPropertyHwAddress()
	m["TypeDescription"], _ = d.GetPropertyTypeDescription()
	return json.Marshal(m)
}
