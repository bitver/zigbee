package main

import (
	"fmt"
	"github.com/bitver/zigbee"
	"github.com/bitver/zigbee/configuration"
	"github.com/bitver/zigbee/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/dyrkin/zcl-go/cluster"
	"sync"
)

//simple device database
var devices = map[string]*model.Device{}

func main() {

	conf := configuration.Default()
	conf.PermitJoin = true

	stewie := steward.New(conf)

	eventListener := func() {
		for {
			select {
			case device := <-stewie.Channels().OnDeviceRegistered():
				saveDevice(device)
			case device := <-stewie.Channels().OnDeviceUnregistered():
				deleteDevice(device)
			case device := <-stewie.Channels().OnDeviceBecameAvailable():
				saveDevice(device)
			case deviceIncomingMessage := <-stewie.Channels().OnDeviceIncomingMessage():
				fmt.Printf("Device received incoming message:\n%s", spew.Sdump(deviceIncomingMessage))
				toggleIkeaBulb(stewie, deviceIncomingMessage)
			}
		}
	}

	go eventListener()
	stewie.Start()
	infiniteWait()
}

func toggleIkeaBulb(stewie *steward.Steward, message *model.DeviceIncomingMessage) {
	if isXiaomiButtonSingleClick(message) {
		if ikeaBulb, registered := devices["TRADFRI bulb E27 W opal 1000lm"]; registered {
			toggleTarget(stewie, ikeaBulb.NetworkAddress)
		} else {
			fmt.Println("IKEA bulb is not available")
		}
	}
}

func toggleTarget(stewie *steward.Steward, networkAddress string) {
	go func() {
		stewie.Functions().Cluster().Local().OnOff().Toggle(networkAddress, 0xFF)
	}()
}

func isXiaomiButtonSingleClick(message *model.DeviceIncomingMessage) bool {
	command, ok := message.IncomingMessage.Data.Command.(*cluster.ReportAttributesCommand)

	return ok && message.Device.Manufacturer == "LUMI" &&
		message.Device.Model == "lumi.remote.b186acn01\x00\x00\x00" &&
		isSingleClick(command)
}

func isSingleClick(command *cluster.ReportAttributesCommand) bool {
	click, ok := command.AttributeReports[0].Attribute.Value.(uint64)
	return ok && click == uint64(1)
}

func saveDevice(device *model.Device) {
	fmt.Printf("Registering device:\n%s", spew.Sdump(device))
	devices[device.Model] = device
}

func deleteDevice(device *model.Device) {
	fmt.Printf("Unregistering device:\n%s", spew.Sdump(device))
	delete(devices, device.Model)
}

func infiniteWait() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

//TODO Remove this
func subscribeForLevelControlEvents(stewie *steward.Steward, device *model.Device) {
	if device.Manufacturer == "IKEA of Sweden" && device.Model == "TRADFRI wireless dimmer" {
		go func() {
			rsp, err := stewie.Functions().Generic().Bind(device.NetworkAddress, device.IEEEAddress, 1,
				uint16(cluster.LevelControl), stewie.Configuration().IEEEAddress, 1)
			fmt.Printf("Bind result: [%v] [%s]", rsp, err)
		}()
	}
}
