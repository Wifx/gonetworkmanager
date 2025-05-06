package gonetworkmanager

import (
	"encoding/json"
	"errors"

	"github.com/godbus/dbus/v5"
)

const (
	IP4ConfigInterface = NetworkManagerInterface + ".IP4Config"

	/* Properties */
	IP4ConfigPropertyAddresses      = IP4ConfigInterface + ".Addresses"      // readable   aau
	IP4ConfigPropertyAddressData    = IP4ConfigInterface + ".AddressData"    // readable   aa{sv}
	IP4ConfigPropertyGateway        = IP4ConfigInterface + ".Gateway"        // readable   s
	IP4ConfigPropertyRoutes         = IP4ConfigInterface + ".Routes"         // readable   aau
	IP4ConfigPropertyRouteData      = IP4ConfigInterface + ".RouteData"      // readable   aa{sv}
	IP4ConfigPropertyNameservers    = IP4ConfigInterface + ".Nameservers"    // readable   au
	IP4ConfigPropertyNameserverData = IP4ConfigInterface + ".NameserverData" // readable   aa{sv}
	IP4ConfigPropertyDomains        = IP4ConfigInterface + ".Domains"        // readable   as
	IP4ConfigPropertySearches       = IP4ConfigInterface + ".Searches"       // readable   as
	IP4ConfigPropertyDnsOptions     = IP4ConfigInterface + ".DnsOptions"     // readable   as
	IP4ConfigPropertyDnsPriority    = IP4ConfigInterface + ".DnsPriority"    // readable   i
	IP4ConfigPropertyWinsServers    = IP4ConfigInterface + ".WinsServers"    // readable   au
	IP4ConfigPropertyWinsServerData = IP4ConfigInterface + ".WinsServerData" // readable   as
)

// Deprecated: use IP4AddressData instead
type IP4Address struct {
	Address string
	Prefix  uint8
	Gateway string
}

type IP4AddressData struct {
	Address string
	Prefix  uint32
}

// Deprecated: use IP4RouteData instead
type IP4Route struct {
	Route   string
	Prefix  uint8
	NextHop string
	Metric  uint8
}

type IP4RouteData struct {
	Destination          string
	Prefix               uint32
	NextHop              string
	Metric               uint32
	AdditionalAttributes map[string]string
}

type IP4NameserverData struct {
	Address string
}

type IP4Config interface {
	// GetPropertyAddresses Array of arrays of IPv4 address/prefix/gateway. All 3 elements of each array are in network byte order. Essentially: [(addr, prefix, gateway), (addr, prefix, gateway), ...]
	// Deprecated: use AddressData and Gateway
	GetPropertyAddresses() ([]IP4Address, error)

	// GetPropertyAddressData Array of IP address data objects. All addresses will include "address" (an IP address string), and "prefix" (a uint). Some addresses may include additional attributes.
	GetPropertyAddressData() ([]IP4AddressData, error)

	// GetPropertyGateway The gateway in use.
	GetPropertyGateway() (string, error)

	// GetPropertyRoutes Arrays of IPv4 route/prefix/next-hop/metric. All 4 elements of each tuple are in network byte order. 'route' and 'next hop' are IPv4 addresses, while prefix and metric are simple unsigned integers. Essentially: [(route, prefix, next-hop, metric), (route, prefix, next-hop, metric), ...]
	// Deprecated: use RouteData
	GetPropertyRoutes() ([]IP4Route, error)

	// GetPropertyRouteData Array of IP route data objects. All routes will include "dest" (an IP address string) and "prefix" (a uint). Some routes may include "next-hop" (an IP address string), "metric" (a uint), and additional attributes.
	GetPropertyRouteData() ([]IP4RouteData, error)

	// GetPropertyNameservers The nameservers in use.
	// Deprecated: use NameserverData
	GetPropertyNameservers() ([]string, error)

	// GetPropertyNameserverData The nameservers in use. Currently only the value "address" is recognized (with an IP address string).
	GetPropertyNameserverData() ([]IP4NameserverData, error)

	// GetPropertyDomains A list of domains this address belongs to.
	GetPropertyDomains() ([]string, error)

	// GetPropertySearches A list of dns searches.
	GetPropertySearches() ([]string, error)

	// GetPropertyDnsOptions A list of DNS options that modify the behavior of the DNS resolver. See resolv.conf(5) manual page for the list of supported options.
	GetPropertyDnsOptions() ([]string, error)

	// GetPropertyDnsPriority The relative priority of DNS servers.
	GetPropertyDnsPriority() (uint32, error)

	// GetPropertyWinsServerData The Windows Internet Name Service servers associated with the connection.
	GetPropertyWinsServerData() ([]string, error)

	MarshalJSON() ([]byte, error)
}

func NewIP4Config(objectPath dbus.ObjectPath) (IP4Config, error) {
	var c ip4Config
	return &c, c.init(NetworkManagerInterface, objectPath)
}

type ip4Config struct {
	dbusBase
}

// Deprecated: use GetPropertyAddressData
func (c *ip4Config) GetPropertyAddresses() ([]IP4Address, error) {
	addresses, err := c.getSliceSliceUint32Property(IP4ConfigPropertyAddresses)
	ret := make([]IP4Address, len(addresses))

	if err != nil {
		return ret, err
	}

	for i, parts := range addresses {
		ret[i] = IP4Address{
			Address: ip4ToString(parts[0]),
			Prefix:  uint8(parts[1]),
			Gateway: ip4ToString(parts[2]),
		}
	}

	return ret, nil
}

func (c *ip4Config) GetPropertyAddressData() ([]IP4AddressData, error) {
	addresses, err := c.getSliceMapStringVariantProperty(IP4ConfigPropertyAddressData)
	if err != nil {
		return []IP4AddressData{}, err
	}
	return DecodeIP4AddressData(addresses)
}

func DecodeIP4AddressData(addresses []map[string]dbus.Variant) ([]IP4AddressData, error) {
	addressesData := make([]IP4AddressData, 0, len(addresses))

	for _, address := range addresses {
		addressDataObj := IP4AddressData{}

		for addressDataAttributeName, addressDataAttribute := range address {
			switch addressDataAttributeName {
			case "address":
				address, ok := addressDataAttribute.Value().(string)
				if !ok {
					return addressesData, errors.New("unexpected variant type for address")
				}
				addressDataObj.Address = address
			case "prefix":
				prefix, ok := addressDataAttribute.Value().(uint32)
				if !ok {
					return addressesData, errors.New("unexpected variant type for prefix")
				}
				addressDataObj.Prefix = prefix
			}
		}

		addressesData = append(addressesData, addressDataObj)
	}
	return addressesData, nil
}

func (c *ip4Config) GetPropertyGateway() (string, error) {
	return c.getStringProperty(IP4ConfigPropertyGateway)
}

// Deprecated: use GetPropertyRouteData
func (c *ip4Config) GetPropertyRoutes() ([]IP4Route, error) {
	routes, err := c.getSliceSliceUint32Property(IP4ConfigPropertyRoutes)
	ret := make([]IP4Route, len(routes))

	if err != nil {
		return ret, err
	}

	for i, parts := range routes {
		ret[i] = IP4Route{
			Route:   ip4ToString(parts[0]),
			Prefix:  uint8(parts[1]),
			NextHop: ip4ToString(parts[2]),
			Metric:  uint8(parts[3]),
		}
	}

	return ret, nil
}

func (c *ip4Config) GetPropertyRouteData() ([]IP4RouteData, error) {
	routes, err := c.getSliceMapStringVariantProperty(IP4ConfigPropertyRouteData)
	if err != nil {
		return []IP4RouteData{}, err
	}
	return DecodeIP4RouteData(routes)
}

func DecodeIP4RouteData(routes []map[string]dbus.Variant) ([]IP4RouteData, error) {
	routesData := make([]IP4RouteData, 0, len(routes))

	for _, route := range routes {
		routeDataObj := IP4RouteData{}

		for routeDataAttributeName, routeDataAttribute := range route {
			switch routeDataAttributeName {
			case "dest":
				destination, ok := routeDataAttribute.Value().(string)
				if !ok {
					return routesData, errors.New("unexpected variant type for dest")
				}
				routeDataObj.Destination = destination
			case "prefix":
				prefix, ok := routeDataAttribute.Value().(uint32)
				if !ok {
					return routesData, errors.New("unexpected variant type for prefix")
				}
				routeDataObj.Prefix = prefix
			case "next-hop":
				nextHop, ok := routeDataAttribute.Value().(string)
				if !ok {
					return routesData, errors.New("unexpected variant type for next-hop")
				}
				routeDataObj.NextHop = nextHop
			case "metric":
				metric, ok := routeDataAttribute.Value().(uint32)
				if !ok {
					return routesData, errors.New("unexpected variant type for metric")
				}
				routeDataObj.Metric = metric
			default:
				if routeDataObj.AdditionalAttributes == nil {
					routeDataObj.AdditionalAttributes = make(map[string]string)
				}
				routeDataObj.AdditionalAttributes[routeDataAttributeName] = routeDataAttribute.String()
			}
		}

		routesData = append(routesData, routeDataObj)
	}
	return routesData, nil
}

// Deprecated: use GetPropertyNameserverData
func (c *ip4Config) GetPropertyNameservers() ([]string, error) {
	nameservers, err := c.getSliceUint32Property(IP4ConfigPropertyNameservers)
	ret := make([]string, len(nameservers))

	if err != nil {
		return ret, err
	}

	for i, ns := range nameservers {
		ret[i] = ip4ToString(ns)
	}

	return ret, nil
}

func (c *ip4Config) GetPropertyNameserverData() ([]IP4NameserverData, error) {
	nameserversData, err := c.getSliceMapStringVariantProperty(IP4ConfigPropertyNameserverData)
	if err != nil {
		return nil, err
	}
	return DecodeIP4NameserverData(nameserversData)
}

func DecodeIP4NameserverData(nameserverData []map[string]dbus.Variant) ([]IP4NameserverData, error) {
	nameservers := make([]IP4NameserverData, 0, len(nameserverData))

	for _, nsData := range nameserverData {
		address, ok := nsData["address"].Value().(string)
		if !ok {
			return nameservers, errors.New("unexpected variant type for address")
		}

		nameserver := IP4NameserverData{
			Address: address,
		}

		nameservers = append(nameservers, nameserver)
	}
	return nameservers, nil
}

func (c *ip4Config) GetPropertyDomains() ([]string, error) {
	return c.getSliceStringProperty(IP4ConfigPropertyDomains)
}

func (c *ip4Config) GetPropertySearches() ([]string, error) {
	return c.getSliceStringProperty(IP4ConfigPropertySearches)
}

func (c *ip4Config) GetPropertyDnsOptions() ([]string, error) {
	return c.getSliceStringProperty(IP4ConfigPropertyDnsOptions)
}

func (c *ip4Config) GetPropertyDnsPriority() (uint32, error) {
	return c.getUint32Property(IP4ConfigPropertyDnsPriority)
}

func (c *ip4Config) GetPropertyWinsServerData() ([]string, error) {
	return c.getSliceStringProperty(IP4ConfigPropertyWinsServerData)
}

func (c *ip4Config) MarshalJSON() ([]byte, error) {
	Addresses, err := c.GetPropertyAddressData()
	if err != nil {
		return nil, err
	}
	Routes, err := c.GetPropertyRouteData()
	if err != nil {
		return nil, err
	}
	Nameservers, err := c.GetPropertyNameserverData()
	if err != nil {
		return nil, err
	}
	Domains, err := c.GetPropertyDomains()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"Addresses":   Addresses,
		"Routes":      Routes,
		"Nameservers": Nameservers,
		"Domains":     Domains,
	})
}
