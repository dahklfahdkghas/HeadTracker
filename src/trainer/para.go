package trainer

// Bluetooth (FrSKY's PARA trainer protocol) link

import (
	"encoding/binary"
	"time"

	"tinygo.org/x/bluetooth"
)

// var serviceUUID = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
var serviceUUID = bluetooth.NewUUID([16]byte{0x61, 0xf4, 0x54, 0x33, 0x79, 0x32, 0x49, 0xda, 0xa7, 0xe8, 0xbb, 0x13, 0x06, 0x04, 0xd5, 0x03})
var charUUID = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}

//var suid = []byte("Helloservice!!!!")
//var cuid = []byte("Hellocharacteris")

// var serviceUUID = [16]byte{suid[0], suid[1], suid[2], suid[3], suid[4], suid[5], suid[6], suid[7], suid[8], suid[9], suid[10], suid[11], suid[12], suid[13], suid[14], suid[15]}
//var charUUID = [16]byte{cuid[0], cuid[1], cuid[2], cuid[3], cuid[4], cuid[5], cuid[6], cuid[7], cuid[8], cuid[9], cuid[10], cuid[11], cuid[12], cuid[13], cuid[14], cuid[15]}

type Para struct {
	adapter *bluetooth.Adapter
	adv     *bluetooth.Advertisement

	yprCharacteristicHandle bluetooth.Characteristic

	yCharacteristicHandle bluetooth.Characteristic
	pCharacteristicHandle bluetooth.Characteristic
	rCharacteristicHandle bluetooth.Characteristic

	buffer    []byte
	sendDelay time.Duration

	paired   bool
	address  string
	channels [8]uint16

	y uint16
	p uint16
	r uint16
}

func NewPara() *Para {
	return &Para{
		adapter: bluetooth.DefaultAdapter,
		//buffer:   make([]byte, 20),
		paired:   false,
		address:  "B1:6B:00:B5:BA:BE",
		channels: [8]uint16{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500},
		y:        uint16(11),
		p:        uint16(22),
		r:        uint16(33),
	}
}

func (t *Para) Configure() {

	t.adapter.Enable()

	t.adapter.SetConnectHandler(func(device bluetooth.Address, connected bool) {
		//t.adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		if connected {
			t.paired = true
		} else {
			t.paired = false
		}
	})

	t.adv = t.adapter.DefaultAdvertisement()
	t.adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Headtracker",

		//ServiceUUIDs: []bluetooth.UUID{bluetooth.NewUUID(serviceUUID)},
		//ServiceUUIDs: []bluetooth.UUID{bluetooth.New16BitUUID(0xFFF0)},
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	})
	t.adv.Start()

	//var yprCharacteristic bluetooth.Characteristic

	t.adapter.AddService(&bluetooth.Service{
		//UUID: bluetooth.NewUUID(serviceUUID),
		UUID: bluetooth.New16BitUUID(0xFFF0),
		Characteristics: []bluetooth.CharacteristicConfig{
			/*
				{
					Handle: &t.yprCharacteristicHandle,
					//UUID:   bluetooth.NewUUID(charUUID),
					UUID:  bluetooth.New16BitUUID(0xFFF1),
					Value: []byte{0, 0, 0, 0, 0, 0, 0, 0},
					Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,

				},
			*/

			{
				Handle: &t.yCharacteristicHandle,
				UUID:   bluetooth.New16BitUUID(0xFFF1),
				Value:  []byte{1, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},

			{
				Handle: &t.pCharacteristicHandle,
				UUID:   bluetooth.New16BitUUID(0xFFF2),
				Value:  []byte{2, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},

			{
				Handle: &t.rCharacteristicHandle,
				UUID:   bluetooth.New16BitUUID(0xFFF3),
				Value:  []byte{3, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})

	addr, _ := t.adapter.Address()
	t.address = addr.MAC.String()

}

func (t *Para) Run() {
	period := 20 * time.Millisecond
	for {
		time.Sleep(period)
		if !t.paired {
			continue
		}

		/*
			numChannels := 8
			numBytesPerChannel := 2
			b := make([]byte, numChannels*numBytesPerChannel)
			for i := 0; i < numChannels; i += 1 {
				offset := i * 2
				binary.BigEndian.PutUint16(b[offset:], t.channels[i])
			}
			t.yprCharacteristicHandle.Write(b)
		*/

		{
			by := make([]byte, 2)
			binary.BigEndian.PutUint16(by, t.channels[2])
			t.yCharacteristicHandle.Write(by)
		}

		{
			bp := make([]byte, 2)
			binary.BigEndian.PutUint16(bp, t.channels[1])
			t.pCharacteristicHandle.Write(bp)
		}

		{
			br := make([]byte, 2)
			binary.BigEndian.PutUint16(br, t.channels[0])
			t.rCharacteristicHandle.Write(br)
		}

	}
}

func (p *Para) Paired() bool {
	return p.paired
}

func (p *Para) Address() string {
	return p.address
}

func (p *Para) Channels() []uint16 {
	return p.channels[:3]
}

func (p *Para) SetChannel(n int, v uint16) {
	p.channels[n] = v
}
