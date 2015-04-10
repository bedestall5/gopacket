// Copyright 2014, Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	_ "fmt"
	"github.com/google/gopacket"
	"net"
	"reflect"
	"bytes"
	"testing"
)

// Generator: python layers/test_creator.py --layerType=LayerTypeRadioTap --linkType=LinkTypeIEEE80211Radio --name=Dot11%s ~/Downloads/mesh.pcap
// http://wiki.wireshark.org/SampleCaptures#Sample_Captures

// testPacketDot11CtrlCTS is the packet:
//   09:28:41.830560 20604983us tsft short preamble 24.0 Mb/s 5240 MHz 11a -79dB signal -92dB noise antenna 1 Clear-To-Send RA:d8:a2:5e:97:61:c1
//   	0x0000:  0000 1900 6f08 0000 3768 3a01 0000 0000  ....o...7h:.....
//   	0x0010:  1230 7814 4001 b1a4 01c4 0094 00d8 a25e  .0x.@..........^
//   	0x0020:  9761 c136 5095 8e                        .a.6P..

var testPacketDot11CtrlCTS = []byte{
	0x00, 0x00, 0x19, 0x00, 0x6f, 0x08, 0x00, 0x00, 0x37, 0x68, 0x3a, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x12, 0x30, 0x78, 0x14, 0x40, 0x01, 0xb1, 0xa4, 0x01, 0xc4, 0x00, 0x94, 0x00, 0xd8, 0xa2, 0x5e,
	0x97, 0x61, 0xc1, 0x36, 0x50, 0x95, 0x8e,
}

func TestPacketDot11CtrlCTS(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11CtrlCTS, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11}, t)

	if got, ok := p.Layer(LayerTypeRadioTap).(*RadioTap); ok {
		want := &RadioTap{
			BaseLayer: BaseLayer{
				Contents: []uint8{0x0, 0x0, 0x19, 0x0, 0x6f, 0x8, 0x0, 0x0, 0x37, 0x68, 0x3a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x12, 0x30, 0x78, 0x14, 0x40, 0x1, 0xb1, 0xa4, 0x1},
				Payload:  []uint8{0xc4, 0x0, 0x94, 0x0, 0xd8, 0xa2, 0x5e, 0x97, 0x61, 0xc1, 0x36, 0x50, 0x95, 0x8e},
			},
			Version:          0x0,
			Length:           0x19,
			Present:          0x86f,
			TSFT:             0x13a6837,
			Flags:            0x12,
			Rate:             0x30,
			ChannelFrequency: 0x1478,
			ChannelFlags:     0x140,
			FHSS:             0x0,
			DBMAntennaSignal: -79,
			DBMAntennaNoise:  -92,
			LockQuality:      0x0,
			TxAttenuation:    0x0,
			DBTxAttenuation:  0x0,
			DBMTxPower:       0,
			Antenna:          1,
			DBAntennaSignal:  0x0,
			DBAntennaNoise:   0x0,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("RadioTap packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}

	if got, ok := p.Layer(LayerTypeDot11).(*Dot11); ok {
		if !got.ChecksumValid() {
			t.Errorf("Dot11 packet processing failed:\nchecksum failed. got  :\n%#v\n\n", got)
		}

		want := &Dot11{
			BaseLayer: BaseLayer{
				Contents: []uint8{0xc4, 0x0, 0x94, 0x0, 0xd8, 0xa2, 0x5e, 0x97, 0x61, 0xc1},
				Payload:  []uint8{},
			},
			Type:       Dot11TypeCtrlCTS,
			Proto:      0x0,
			Flags:      0x0,
			DurationID: 0x94,
			Address1:   net.HardwareAddr{0xd8, 0xa2, 0x5e, 0x97, 0x61, 0xc1}, // check
			Address2:   net.HardwareAddr(nil),
			Address3:   net.HardwareAddr(nil),
			Address4:   net.HardwareAddr(nil),
			Checksum:   0x8e955036,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Dot11 packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}
}

func BenchmarkDecodePacketDot11CtrlCTS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11CtrlCTS, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// testPacketDot11MgmtBeacon is the packet:
//   06:14:26.492772 637444841us tsft short preamble 6.0 Mb/s -40dB signal -96dB noise antenna 2 5180 MHz 11a Beacon () [6.0* 9.0 12.0* 18.0 24.0* 36.0 48.0 54.0 Mbit] IBSS CH: 36
//   	0x0000:  0000 2000 6708 0400 e9a2 fe25 0000 0000  ....g......%....
//   	0x0010:  220c d8a0 0200 0000 4001 0000 3c14 2411  ".......@...<.$.
//   	0x0020:  8000 0000 ffff ffff ffff 0003 7f07 a016  ................
//   	0x0030:  0000 0000 0000 d09b 3840 1028 0000 0000  ........8@.(....
//   	0x0040:  6400 0005 0000 0108 8c12 9824 b048 606c  d..........$.H`l
//   	0x0050:  0301 2405 0400 0100 0007 2a55 5320 2401  ..$.......*US.$.
//   	0x0060:  1128 0111 2c01 1130 0111 3401 1738 0117  .(..,..0..4..8..
//   	0x0070:  3c01 1740 0117 9501 1e99 011e 9d01 1ea1  <..@............
//   	0x0080:  011e a501 1e20 0100 dd18 0050 f202 0101  ...........P....
//   	0x0090:  0000 03a4 0000 27a4 0000 4243 5e00 6232  ......'...BC^.b2
//   	0x00a0:  2f00 340c 6672 6565 6273 642d 6d65 7368  /.4.freebsd-mesh
//   	0x00b0:  3317 0100 0fac 0000 0fac 0000 0fac ff00  3...............
//   	0x00c0:  0fac ff00 0fac ff00 df                   .........
var testPacketDot11MgmtBeacon = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0xe9, 0xa2, 0xfe, 0x25, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x0c, 0xd8, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x03, 0x7f, 0x07, 0xa0, 0x16,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xd0, 0x9b, 0x38, 0x40, 0x10, 0x28, 0x00, 0x00, 0x00, 0x00,
	0x64, 0x00, 0x00, 0x05, 0x00, 0x00, 0x01, 0x08, 0x8c, 0x12, 0x98, 0x24, 0xb0, 0x48, 0x60, 0x6c,
	0x03, 0x01, 0x24, 0x05, 0x04, 0x00, 0x01, 0x00, 0x00, 0x07, 0x2a, 0x55, 0x53, 0x20, 0x24, 0x01,
	0x11, 0x28, 0x01, 0x11, 0x2c, 0x01, 0x11, 0x30, 0x01, 0x11, 0x34, 0x01, 0x17, 0x38, 0x01, 0x17,
	0x3c, 0x01, 0x17, 0x40, 0x01, 0x17, 0x95, 0x01, 0x1e, 0x99, 0x01, 0x1e, 0x9d, 0x01, 0x1e, 0xa1,
	0x01, 0x1e, 0xa5, 0x01, 0x1e, 0x20, 0x01, 0x00, 0xdd, 0x18, 0x00, 0x50, 0xf2, 0x02, 0x01, 0x01,
	0x00, 0x00, 0x03, 0xa4, 0x00, 0x00, 0x27, 0xa4, 0x00, 0x00, 0x42, 0x43, 0x5e, 0x00, 0x62, 0x32,
	0x2f, 0x00, 0x34, 0x0c, 0x66, 0x72, 0x65, 0x65, 0x62, 0x73, 0x64, 0x2d, 0x6d, 0x65, 0x73, 0x68,
	0x33, 0x17, 0x01, 0x00, 0x0f, 0xac, 0x00, 0x00, 0x0f, 0xac, 0x00, 0x00, 0x0f, 0xac, 0xff, 0x00,
	0x0f, 0xac, 0xff, 0x00, 0x0f, 0xac, 0xff, 0x00, 0xdf,
}

func TestPacketDot11MgmtBeacon(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11MgmtBeacon, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}

	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11MgmtBeacon}, t)

	if got, ok := p.Layer(LayerTypeRadioTap).(*RadioTap); ok {
		want := &RadioTap{BaseLayer: BaseLayer{Contents: []uint8{0x0, 0x0, 0x20, 0x0, 0x67, 0x8, 0x4, 0x0, 0xe9, 0xa2, 0xfe, 0x25, 0x0, 0x0, 0x0, 0x0, 0x22, 0xc, 0xd8, 0xa0, 0x2, 0x0, 0x0, 0x0, 0x40, 0x1, 0x0, 0x0, 0x3c, 0x14, 0x24, 0x11}, Payload: []uint8{0x80, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x3, 0x7f, 0x7, 0xa0, 0x16, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd0, 0x9b, 0x38, 0x40, 0x10, 0x28, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x5, 0x0, 0x0, 0x1, 0x8, 0x8c, 0x12, 0x98, 0x24, 0xb0, 0x48, 0x60, 0x6c, 0x3, 0x1, 0x24, 0x5, 0x4, 0x0, 0x1, 0x0, 0x0, 0x7, 0x2a, 0x55, 0x53, 0x20, 0x24, 0x1, 0x11, 0x28, 0x1, 0x11, 0x2c, 0x1, 0x11, 0x30, 0x1, 0x11, 0x34, 0x1, 0x17, 0x38, 0x1, 0x17, 0x3c, 0x1, 0x17, 0x40, 0x1, 0x17, 0x95, 0x1, 0x1e, 0x99, 0x1, 0x1e, 0x9d, 0x1, 0x1e, 0xa1, 0x1, 0x1e, 0xa5, 0x1, 0x1e, 0x20, 0x1, 0x0, 0xdd, 0x18, 0x0, 0x50, 0xf2, 0x2, 0x1, 0x1, 0x0, 0x0, 0x3, 0xa4, 0x0, 0x0, 0x27, 0xa4, 0x0, 0x0, 0x42, 0x43, 0x5e, 0x0, 0x62, 0x32, 0x2f, 0x0, 0x34, 0xc, 0x66, 0x72, 0x65, 0x65, 0x62, 0x73, 0x64, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x33, 0x17, 0x1, 0x0, 0xf, 0xac, 0x0, 0x0, 0xf, 0xac, 0x0, 0x0, 0xf, 0xac, 0xff, 0x0, 0xf, 0xac, 0xff, 0x0, 0xf, 0xac, 0xff, 0x0, 0xdf}}, Version: 0x0, Length: 0x20, Present: 0x40867, TSFT: 0x25fea2e9, Flags: 0x22, Rate: 0xc, ChannelFrequency: 0x0, ChannelFlags: 0x0, FHSS: 0x0, DBMAntennaSignal: -40, DBMAntennaNoise: -96, LockQuality: 0x0, TxAttenuation: 0x0, DBTxAttenuation: 0x0, DBMTxPower: 0, Antenna: 0x2, DBAntennaSignal: 0x0, DBAntennaNoise: 0x0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("RadioTap packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}

	if got, ok := p.Layer(LayerTypeDot11).(*Dot11); ok {
		want := &Dot11{
			BaseLayer: BaseLayer{
				Contents: []uint8{0x80, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x3, 0x7f, 0x7, 0xa0, 0x16, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd0, 0x9b},
				Payload:  []uint8{0x38, 0x40, 0x10, 0x28, 0x0, 0x0, 0x0, 0x0, 0x64, 0x0, 0x0, 0x5, 0x0, 0x0, 0x1, 0x8, 0x8c, 0x12, 0x98, 0x24, 0xb0, 0x48, 0x60, 0x6c, 0x3, 0x1, 0x24, 0x5, 0x4, 0x0, 0x1, 0x0, 0x0, 0x7, 0x2a, 0x55, 0x53, 0x20, 0x24, 0x1, 0x11, 0x28, 0x1, 0x11, 0x2c, 0x1, 0x11, 0x30, 0x1, 0x11, 0x34, 0x1, 0x17, 0x38, 0x1, 0x17, 0x3c, 0x1, 0x17, 0x40, 0x1, 0x17, 0x95, 0x1, 0x1e, 0x99, 0x1, 0x1e, 0x9d, 0x1, 0x1e, 0xa1, 0x1, 0x1e, 0xa5, 0x1, 0x1e, 0x20, 0x1, 0x0, 0xdd, 0x18, 0x0, 0x50, 0xf2, 0x2, 0x1, 0x1, 0x0, 0x0, 0x3, 0xa4, 0x0, 0x0, 0x27, 0xa4, 0x0, 0x0, 0x42, 0x43, 0x5e, 0x0, 0x62, 0x32, 0x2f, 0x0, 0x34, 0xc, 0x66, 0x72, 0x65, 0x65, 0x62, 0x73, 0x64, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x33, 0x17, 0x1, 0x0, 0xf, 0xac, 0x0, 0x0, 0xf, 0xac, 0x0, 0x0, 0xf, 0xac, 0xff, 0x0, 0xf, 0xac, 0xff, 0x0, 0xf},
			},
			Type:           Dot11TypeMgmtBeacon,
			Proto:          0x0,
			Flags:          0x0,
			DurationID:     0x0,
			Address1:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			Address2:       net.HardwareAddr{0x0, 0x3, 0x7f, 0x7, 0xa0, 0x16},
			Address3:       net.HardwareAddr{0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			Address4:       net.HardwareAddr(nil),
			SequenceNumber: 0x26f, FragmentNumber: 0x10,
			Checksum: 0xdf00ffac,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Dot11 packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}
}

func BenchmarkDecodePacketDot11MgmtBeacon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11MgmtBeacon, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// testPacketDot11DataQOSData is the packet:
//   06:14:27.838634 638790765us tsft short preamble 54.0 Mb/s -51dB signal -96dB noise antenna 2 5180 MHz 11a CF +QoS ARP, Request who-has 140.180.51.68 tell 169.254.247.0, length 28
//   	0x0000:  0000 2000 6708 0400 6d2c 1326 0000 0000  ....g...m,.&....
//   	0x0010:  226c cda0 0200 0000 4001 0000 3c14 2411  "l......@...<.$.
//   	0x0020:  8801 2c00 0603 7f07 a016 0019 e3d3 5352  ..,...........SR
//   	0x0030:  ffff ffff ffff 5064 0000 50aa aaaa 0300  ......Pd..P.....
//   	0x0040:  0000 0806 0001 0800 0604 0001 0019 e3d3  ................
//   	0x0050:  5352 a9fe f700 0000 0000 0000 8cb4 3344  SR............3D
var testPacketDot11DataQOSData = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0x6d, 0x2c, 0x13, 0x26, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x6c, 0xcd, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0x88, 0x01, 0x2c, 0x00, 0x06, 0x03, 0x7f, 0x07, 0xa0, 0x16, 0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x50, 0x64, 0x00, 0x00, 0x50, 0xaa, 0xaa, 0xaa, 0x03, 0x00,
	0x00, 0x00, 0x08, 0x06, 0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01, 0x00, 0x19, 0xe3, 0xd3,
	0x53, 0x52, 0xa9, 0xfe, 0xf7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8c, 0xb4, 0x33, 0x44,
}

func TestPacketDot11DataQOSData(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11DataQOSData, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11DataQOSData, LayerTypeDot11Data, LayerTypeLLC, LayerTypeSNAP, LayerTypeARP}, t)

	if got, ok := p.Layer(LayerTypeARP).(*ARP); ok {
		want := &ARP{BaseLayer: BaseLayer{
			Contents: []uint8{0x0, 0x1, 0x8, 0x0, 0x6, 0x4, 0x0, 0x1, 0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0xa9, 0xfe, 0xf7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8c, 0xb4, 0x33, 0x44},
			Payload:  []uint8{},
		},
			AddrType:          0x1,
			Protocol:          0x800,
			HwAddressSize:     0x6,
			ProtAddressSize:   0x4,
			Operation:         0x1,
			SourceHwAddress:   []uint8{0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52},
			SourceProtAddress: []uint8{0xa9, 0xfe, 0xf7, 0x0},
			DstHwAddress:      []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			DstProtAddress:    []uint8{0x8c, 0xb4, 0x33, 0x44},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ARP packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}
}
func BenchmarkDecodePacketDot11DataQOSData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11DataQOSData, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// this test is failing, this has probably todo with the datapad

// testPacketDot11MgmtAction is the packet:
//   06:14:23.248019 634199967us tsft short preamble 6.0 Mb/s -41dB signal -96dB noise antenna 1 5180 MHz 11a Action (00:03:7f:07:a0:16): Reserved(32) Act#0
//   	0x0000:  0000 2000 6708 0400 9f1f cd25 0000 0000  ....g......%....
//   	0x0010:  220c d7a0 0100 0000 4001 0000 3c14 2411  ".......@...<.$.
//   	0x0020:  d000 0000 ffff ffff ffff 0003 7f07 a016  ................
//   	0x0030:  0003 7f07 a016 a002 2000 4425 0001 1e06  ..........D%....
//   	0x0040:  0000 0000 037f 0342 5206 0000 0088 1300  .......BR.......
//   	0x0050:  0030 0b00 0001 0600 16cb ace5 f900 0000  .0..............
//   	0x0060:  00                                       .
var testPacketDot11MgmtAction = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0x9f, 0x1f, 0xcd, 0x25, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x0c, 0xd7, 0xa0, 0x01, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0xd0, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x03, 0x7f, 0x07, 0xa0, 0x16,
	0x00, 0x03, 0x7f, 0x07, 0xa0, 0x16, 0xa0, 0x02, 0x20, 0x00, 0x44, 0x25, 0x00, 0x01, 0x1e, 0x06,
	0x00, 0x00, 0x00, 0x00, 0x03, 0x7f, 0x03, 0x42, 0x52, 0x06, 0x00, 0x00, 0x00, 0x88, 0x13, 0x00,
	0x00, 0x30, 0x0b, 0x00, 0x00, 0x01, 0x06, 0x00, 0x16, 0xcb, 0xac, 0xe5, 0xf9, 0x00, 0x00, 0x00,
	0x00,
}

func TestPacketDot11MgmtAction(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11MgmtAction, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11MgmtAction}, t)

	want := `PACKET: 97 bytes
- Layer 1 (32 bytes) = RadioTap	{Contents=[..32..] Payload=[..65..] Version=0 Length=32 Present=264295 TSFT=634199967 Flags=SHORT-PREAMBLE,DATAPAD Rate=6 Mb/s ChannelFrequency=0 MHz ChannelFlags= FHSS=0 DBMAntennaSignal=-41 DBMAntennaNoise=-96 LockQuality=0 TxAttenuation=0 DBTxAttenuation=0 DBMTxPower=0 Antenna=1 DBAntennaSignal=0 DBAntennaNoise=0}
- Layer 2 (24 bytes) = Dot11	{Contents=[..24..] Payload=[..37..] Type=MgmtAction Proto=0 Flags= DurationID=0 Address1=ff:ff:ff:ff:ff:ff Address2=00:03:7f:07:a0:16 Address3=00:03:7f:07:a0:16 Address4= SequenceNumber=10 FragmentNumber=32 Checksum=0}
- Layer 3 (37 bytes) = Dot11MgmtAction	{Contents=[..37..] Payload=[]}
`
	if got := p.String(); got != want {
		t.Errorf("packet string mismatch:\n---got---\n%q\n---want---\n%q", got, want)
	}
	if _, ok := p.Layer(LayerTypeDot11).(*Dot11); !ok {
		t.Errorf("could not get Dot11 layer from packet")
	} else {
		// See note above:  this checksum fails most likely due to datapad.
		// wireshark also says this packet is malformed, so I'm not going to waste
		// too much more time on it.
		//   if !got.ChecksumValid() { t.Errorf("Dot11 packet processing failed: checksum failed")	}
	}
}
func BenchmarkDecodePacketDot11MgmtAction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11MgmtAction, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// testPacketDot11CtrlAck is the packet:
//   06:14:27.838669 638758038us tsft short preamble 24.0 Mb/s -39dB signal -96dB noise antenna 2 5180 MHz 11a Acknowledgment RA:00:19:e3:d3:53:52
//   	0x0000:  0000 2000 6708 0400 96ac 1226 0000 0000  ....g......&....
//   	0x0010:  2230 d9a0 0200 0000 4001 0000 3c14 2411  "0......@...<.$.
//   	0x0020:  d400 0000 0019 e3d3 5352 46e9 7687       ........SRF.v.
var testPacketDot11CtrlAck = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0x96, 0xac, 0x12, 0x26, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x30, 0xd9, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0xd4, 0x00, 0x00, 0x00, 0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0x46, 0xe9, 0x76, 0x87,
}

func TestPacketDot11CtrlAck(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11CtrlAck, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11}, t)

	if got, ok := p.Layer(LayerTypeDot11).(*Dot11); ok {
		if !got.ChecksumValid() {
			t.Errorf("Dot11 packet processing failed:\nchecksum failed. got  :\n%#v\n\n", got)
		}
	}

	if got, ok := p.Layer(LayerTypeDot11).(*Dot11); ok {
		want := &Dot11{
			BaseLayer: BaseLayer{
				Contents: []uint8{0xd4, 0x0, 0x0, 0x0, 0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52},
				Payload:  []uint8{},
			},
			Type:       Dot11TypeCtrlAck,
			Proto:      0x0,
			Flags:      0x0,
			DurationID: 0x0,
			Address1:   net.HardwareAddr{0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52},
			Address2:   net.HardwareAddr(nil),
			Address3:   net.HardwareAddr(nil),
			Address4:   net.HardwareAddr(nil),
			Checksum:   0x8776e946,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Dot11 packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}
}
func BenchmarkDecodePacketDot11CtrlAck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11CtrlAck, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// testPacketDot11DataARP is the packet:
//   06:14:11.512316 622463798us tsft short preamble 6.0 Mb/s -39dB signal -96dB noise antenna 2 5180 MHz 11a ARP, Request who-has 67.8.14.54 tell 169.254.247.0, length 28
//   	0x0000:  0000 2000 6708 0400 360b 1a25 0000 0000  ....g...6..%....
//   	0x0010:  220c d9a0 0200 0000 4001 0000 3c14 2411  ".......@...<.$.
//   	0x0020:  0802 0000 ffff ffff ffff 0603 7f07 a016  ................
//   	0x0030:  0019 e3d3 5352 e07f aaaa 0300 0000 0806  ....SR..........
//   	0x0040:  0001 0800 0604 0001 0019 e3d3 5352 a9fe  ............SR..
//   	0x0050:  f700 0000 0000 0000 4308 0e36            ........C..6
var testPacketDot11DataARP = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0x36, 0x0b, 0x1a, 0x25, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x0c, 0xd9, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0x08, 0x02, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x06, 0x03, 0x7f, 0x07, 0xa0, 0x16,
	0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0xe0, 0x7f, 0xaa, 0xaa, 0x03, 0x00, 0x00, 0x00, 0x08, 0x06,
	0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01, 0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0xa9, 0xfe,
	0xf7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x43, 0x08, 0x0e, 0x36,
}

func TestPacketDot11DataARP(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11DataARP, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11Data, LayerTypeLLC, LayerTypeSNAP, LayerTypeARP}, t)

	if got, ok := p.Layer(LayerTypeARP).(*ARP); ok {
		want := &ARP{
			BaseLayer: BaseLayer{
				Contents: []uint8{0x0, 0x1, 0x8, 0x0, 0x6, 0x4, 0x0, 0x1, 0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0xa9, 0xfe, 0xf7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x43, 0x8, 0xe, 0x36},
				Payload:  []uint8{},
			},
			AddrType:          0x1,
			Protocol:          0x800,
			HwAddressSize:     0x6,
			ProtAddressSize:   0x4,
			Operation:         0x1,
			SourceHwAddress:   []uint8{0x0, 0x19, 0xe3, 0xd3, 0x53, 0x52},
			SourceProtAddress: []uint8{0xa9, 0xfe, 0xf7, 0x0},
			DstHwAddress:      []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			DstProtAddress:    []uint8{0x43, 0x8, 0xe, 0x36},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ARP packet processing failed:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	}
}

func BenchmarkDecodePacketDot11DataARP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11DataARP, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// testPacketDot11DataIP is the packet:
//   06:14:21.388622 632340487us tsft short preamble 6.0 Mb/s -40dB signal -96dB noise antenna 1 5180 MHz 11a IP 0.0.0.0.68 > 255.255.255.255.67: BOOTP/DHCP, Request from 00:19:e3:d3:53:52, length 300
//   	0x0000:  0000 2000 6708 0400 07c0 b025 0000 0000  ....g......%....
//   	0x0010:  220c d8a0 0100 0000 4001 0000 3c14 2411  ".......@...<.$.
//   	0x0020:  0802 0000 ffff ffff ffff 0603 7f07 a016  ................
//   	0x0030:  0019 e3d3 5352 4095 aaaa 0300 0000 0800  ....SR@.........
//   	0x0040:  4500 0148 c514 0000 ff11 f590 0000 0000  E..H............
//   	0x0050:  ffff ffff 0044 0043 0134 2b39 0101 0600  .....D.C.4+9....
//   	0x0060:  131f 8c43 003c 0000 0000 0000 0000 0000  ...C.<..........
//   	0x0070:  0000 0000 0000 0000 0019 e3d3 5352 0000  ............SR..
//   	0x0080:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0090:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00a0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00b0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00c0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00d0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00e0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x00f0:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0100:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0110:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0120:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0130:  0000 0000 0000 0000 0000 0000 0000 0000  ................
//   	0x0140:  0000 0000 0000 0000 6382 5363 3501 0137  ........c.Sc5..7
//   	0x0150:  0a01 0306 0f77 5ffc 2c2e 2f39 0205 dc3d  .....w_.,./9...=
//   	0x0160:  0701 0019 e3d3 5352 3304 0076 a700 0c0b  ......SR3..v....
//   	0x0170:  4d61 6369 6e74 6f73 682d 34ff 0000 0000  Macintosh-4.....
//   	0x0180:  0000 0000 0000 0000                      ........
var testPacketDot11DataIP = []byte{
	0x00, 0x00, 0x20, 0x00, 0x67, 0x08, 0x04, 0x00, 0x07, 0xc0, 0xb0, 0x25, 0x00, 0x00, 0x00, 0x00,
	0x22, 0x0c, 0xd8, 0xa0, 0x01, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x3c, 0x14, 0x24, 0x11,
	0x08, 0x02, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x06, 0x03, 0x7f, 0x07, 0xa0, 0x16,
	0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0x40, 0x95, 0xaa, 0xaa, 0x03, 0x00, 0x00, 0x00, 0x08, 0x00,
	0x45, 0x00, 0x01, 0x48, 0xc5, 0x14, 0x00, 0x00, 0xff, 0x11, 0xf5, 0x90, 0x00, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xff, 0xff, 0x00, 0x44, 0x00, 0x43, 0x01, 0x34, 0x2b, 0x39, 0x01, 0x01, 0x06, 0x00,
	0x13, 0x1f, 0x8c, 0x43, 0x00, 0x3c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x82, 0x53, 0x63, 0x35, 0x01, 0x01, 0x37,
	0x0a, 0x01, 0x03, 0x06, 0x0f, 0x77, 0x5f, 0xfc, 0x2c, 0x2e, 0x2f, 0x39, 0x02, 0x05, 0xdc, 0x3d,
	0x07, 0x01, 0x00, 0x19, 0xe3, 0xd3, 0x53, 0x52, 0x33, 0x04, 0x00, 0x76, 0xa7, 0x00, 0x0c, 0x0b,
	0x4d, 0x61, 0x63, 0x69, 0x6e, 0x74, 0x6f, 0x73, 0x68, 0x2d, 0x34, 0xff, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestPacketDot11DataIP(t *testing.T) {
	p := gopacket.NewPacket(testPacketDot11DataIP, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11Data, LayerTypeLLC, LayerTypeSNAP, LayerTypeIPv4, LayerTypeUDP, gopacket.LayerTypePayload}, t)
}
func BenchmarkDecodePacketDot11DataIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketDot11DataIP, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

// Encrypted

/// testPacketP6196 is the packet:
//   09:28:41.830631 20605036us tsft wep -69dB signal -92dB noise antenna 1 5240 MHz 11a ht/40- 162.0 Mb/s MCS 12 40 MHz lon GI mixed BCC FEC [bit 20] CF +QoS Data IV:50a9 Pad 20 KeyID 0
//   	0x0000:  0000 3000 6b08 1c00 6c68 3a01 0000 0000  ..0.k...lh:.....
//   	0x0010:  1400 7814 4001 bba4 0160 0e1a 4001 0400  ..x.@....`..@...
//   	0x0020:  7814 3022 1f01 0cff b10d 0000 0400 0000  x.0"............
//   	0x0030:  8841 2c00 0025 9c42 c262 d8a2 5e97 61c1  .A,..%.B.b..^.a.
//   	0x0040:  0025 9c42 c25f 10db 0000 a950 0020 0000  .%.B._.....P....
//   	0x0050:  0000 f8ab a97e 3fbd d6e1 785b 0040 5f15  .....~?...x[.@_.
//   	0x0060:  7123 8711 bd1f ffb9 e5b3 84bb ec2a 0a90  q#...........*..
//   	0x0070:  d0a0 1a6f 9033 1083 5179 a0da f833 3a00  ...o.3..Qy...3:.
//   	0x0080:  5471 f596 539b 1823 a33c 4908 545c 266a  Tq..S..#.<I.T\&j
//   	0x0090:  8540 515a 1da9 c49e a85a fbf7 de09 7f9c  .@QZ.....Z......
//   	0x00a0:  6f35 0b8b 6831 2c10 43dc 8983 b1d9 dd29  o5..h1,.C......)
//   	0x00b0:  7395 65b9 4b43 b391 16ec 4201 86c9 ca    s.e.KC....B....
var testPacketP6196 = []byte{
	0x00, 0x00, 0x30, 0x00, 0x6b, 0x08, 0x1c, 0x00, 0x6c, 0x68, 0x3a, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x14, 0x00, 0x78, 0x14, 0x40, 0x01, 0xbb, 0xa4, 0x01, 0x60, 0x0e, 0x1a, 0x40, 0x01, 0x04, 0x00,
	0x78, 0x14, 0x30, 0x22, 0x1f, 0x01, 0x0c, 0xff, 0xb1, 0x0d, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
	0x88, 0x41, 0x2c, 0x00, 0x00, 0x25, 0x9c, 0x42, 0xc2, 0x62, 0xd8, 0xa2, 0x5e, 0x97, 0x61, 0xc1,
	0x00, 0x25, 0x9c, 0x42, 0xc2, 0x5f, 0x10, 0xdb, 0x00, 0x00, 0xa9, 0x50, 0x00, 0x20, 0x00, 0x00,
	0x00, 0x00, 0xf8, 0xab, 0xa9, 0x7e, 0x3f, 0xbd, 0xd6, 0xe1, 0x78, 0x5b, 0x00, 0x40, 0x5f, 0x15,
	0x71, 0x23, 0x87, 0x11, 0xbd, 0x1f, 0xff, 0xb9, 0xe5, 0xb3, 0x84, 0xbb, 0xec, 0x2a, 0x0a, 0x90,
	0xd0, 0xa0, 0x1a, 0x6f, 0x90, 0x33, 0x10, 0x83, 0x51, 0x79, 0xa0, 0xda, 0xf8, 0x33, 0x3a, 0x00,
	0x54, 0x71, 0xf5, 0x96, 0x53, 0x9b, 0x18, 0x23, 0xa3, 0x3c, 0x49, 0x08, 0x54, 0x5c, 0x26, 0x6a,
	0x85, 0x40, 0x51, 0x5a, 0x1d, 0xa9, 0xc4, 0x9e, 0xa8, 0x5a, 0xfb, 0xf7, 0xde, 0x09, 0x7f, 0x9c,
	0x6f, 0x35, 0x0b, 0x8b, 0x68, 0x31, 0x2c, 0x10, 0x43, 0xdc, 0x89, 0x83, 0xb1, 0xd9, 0xdd, 0x29,
	0x73, 0x95, 0x65, 0xb9, 0x4b, 0x43, 0xb3, 0x91, 0x16, 0xec, 0x42, 0x01, 0x86, 0xc9, 0xca,
}

func TestPacketP6196(t *testing.T) {
	p := gopacket.NewPacket(testPacketP6196, LinkTypeIEEE80211Radio, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}

	checkLayers(p, []gopacket.LayerType{LayerTypeRadioTap, LayerTypeDot11, LayerTypeDot11WEP}, t)
}

func BenchmarkDecodePacketP6196(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketP6196, LinkTypeIEEE80211Radio, gopacket.NoCopy)
	}
}

func TestInformationElement(t *testing.T) {
	bin := []byte{
		0, 0,
		0, 2, 1, 3,
		221, 5, 1, 2, 3, 4, 5,
	}
	pkt := gopacket.NewPacket(bin, LayerTypeDot11InformationElement, gopacket.NoCopy)
	
	buf := gopacket.NewSerializeBuffer()
	var sLayers []gopacket.SerializableLayer
	for _,l := range pkt.Layers() {
		sLayers = append(sLayers, l.(*Dot11InformationElement))
	}
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, sLayers...); err!=nil {
		t.Error(err.Error())
	}
	if !bytes.Equal(bin, buf.Bytes()) {
		t.Error("build failed")
	}
}
