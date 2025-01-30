package trainer

import (
	"encoding/binary"
	"time"

	"tinygo.org/x/bluetooth"
)

type Ble struct {
	adapter    *bluetooth.Adapter
	adv        *bluetooth.Advertisement
	fff6Handle bluetooth.Characteristic

	yprCharacteristicHandle bluetooth.Characteristic

	buffer    []byte
	sendDelay time.Duration

	paired   bool
	address  string
	channels [8]uint16
}

func NewBle() *Ble {
	return &Ble{
		adapter:  bluetooth.DefaultAdapter,
		buffer:   make([]byte, 20),
		paired:   false,
		address:  "B1:6B:00:B5:BA:BE",
		channels: [8]uint16{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
	}
}

func (t *Ble) Configure() {
	t.adapter.Enable()

	yawPitchRoll_CharConfig_fff7 := bluetooth.CharacteristicConfig{
		Handle: &t.yprCharacteristicHandle,
		UUID:   bluetooth.New16BitUUID(0xFFF7),
		Value:  []byte{1, 2, 3, 4, 5, 6},
		Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
	}

	t.adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.New16BitUUID(0xFFF0),
		Characteristics: []bluetooth.CharacteristicConfig{
			yawPitchRoll_CharConfig_fff7,
		},
	})

	t.adv = t.adapter.DefaultAdvertisement()
	t.adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Headtracker",
		ServiceUUIDs: []bluetooth.UUID{bluetooth.New16BitUUID(0xFFF0)},
	})
	t.adv.Start()

	addr, _ := t.adapter.Address()
	t.address = addr.MAC.String()

	t.adapter.SetConnectHandler(func(device bluetooth.Address, connected bool) {
		if connected {
			t.sendDelay = 1 * time.Second
			t.paired = true
		} else {
			t.paired = false
		}
	})
}

func (t *Ble) Run() {
	period := 20 * time.Millisecond
	for {
		time.Sleep(period)
		if !t.paired {
			continue
		}

		//Convert yaw pitch roll to single byte array
		numChannels := 3
		numBytesPerChannel := 2
		b := make([]byte, numChannels*numBytesPerChannel)
		for i := 0; i < numChannels; i += 1 {
			offset := i * 2
			binary.BigEndian.PutUint16(b[offset:], t.channels[i])
		}
		t.yprCharacteristicHandle.Write(b)

	}

}

func (p *Ble) Paired() bool {
	return p.paired
}

func (p *Ble) Address() string {
	return p.address
}

func (p *Ble) Channels() []uint16 {
	return p.channels[:3]
}

func (p *Ble) SetChannel(n int, v uint16) {
	p.channels[n] = v
}
