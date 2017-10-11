package blelib

import (
	"github.com/cksmith/gatt"
)

const (
	DeviceName = "Go Meetup Service"
)

var (
	ServiceUUID                           = gatt.MustParseUUID("1CE37703-5588-411C-9D8E-5FF6FA1FE57B")
	LEDCharacteristicUUID                 = gatt.MustParseUUID("1BED2229-5AE6-42CB-BF20-F79E8CEA8822")
	TemperatureHumidityCharacteristicUUID = gatt.MustParseUUID("8065CDC2-C278-4812-8488-D92D661425D8")
)

var CharacteristicUUIDs = []gatt.UUID{
	LEDCharacteristicUUID,
	TemperatureHumidityCharacteristicUUID,
}
