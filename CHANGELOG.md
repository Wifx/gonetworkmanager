# Changelog

All notable changes to this project since will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.2.0] - 2025-05-06

### Added

- Add GetPath() method to IP4Config and IP6Config interfaces
- Add public decoding functions for IP4 and IP6 data structs

## [3.1.0] - 2025-04-28

### Added

- Generate stringer implementations for all enum types
- Add newest device properties and centralize HwAddress property

## [3.0.0] - 2025-01-24

### Fixed

- BREAKING CHANGE: fix prefixes and metrics as uint8 to uint32 by @tom-wegener in #46
- BREAKING CHANGE: Fix #39 : Device.Reapply using wrong settings argument
- Fix #45: Improve Device types consistency by using DeviceFactory instead of NewDevice

## [2.1.0] - 2023-01-19

### Added

- Add device auto-connect setter

### Fixed

- Connection settings route and address data types

## [2.0.0] - 2022-11-22

### Changes

- Update go-dbus to v5.1.0
- **BREAKING CHANGE** : Generic recursive settings map dbus variants decoding

## [0.5.0] - 2022-11-22

### Added

- Godoc standard comment prefix
- DnsManager
- Device: GetIp4Connectivity
- Constants for Nm80211APSec
- Add gitignore

### Fixed

- CheckpointCreate: fix devicePaths variable scope

## [0.4.0] - 2022-01-17

### Added

- AccessPoint: add LastSeen property

### Changed

- Examples: move examples to their own subfolders

### Fixed

- DeviceWireless: remove duplicated fields
- PrimaryConnection: use ActiveConnection type
- SubscribeState: add the path to the recieved chan type, catch connect event for ActiveConnection

## [0.3.0] - 2020-03-26

### Added

- SetPropertyManaged (@joseffilzmaier)
- GetConnectionByUUID (@paulburlumi)
- VpnConnection
- ActiveConnectionSignalStateChanged
- CheckpointRollback
- SetPropertyWirelessEnabled (@Raqbit)
- Settings.ReloadConnections (@appnostic-io)

Static connection example (@everactivemilligan)

### Fixed

- GetPropertyRouteData panic (@zhengdelun)

## [0.2.0] - 2020-03-06

### Fixed

- added missing flag for Reload
- added parameter specific_object for AddAndActivateConnection
- Fix CheckpointCreateand GetPropertyCheckpoints

### Added

- Add property setter helper
- Add Device.SetPropertyRefreshRateMs
- Add Device.Reapply