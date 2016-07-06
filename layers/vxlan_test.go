// Copyright 2016 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"github.com/postwait/gopacket"
	"reflect"
	"testing"
)

// VXLAN is specifed in RFC 7348
//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  0             8               16              24              32
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |G|R|R|R|I|R|R|R|R|D|R|R|A|R|R|R|       Group Policy ID         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     24 bit VXLAN Network Identifier           |   Reserved    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// Ethernet[IP[UDP[VXLAN[Ethernet[IP[ICMP]]]]]]

// testPacketVXLAN
// 0000   00 16 3e 08 71 cf 36 dc 85 1e b3 40 08 00 45 00  ..>.q.6....@..E.
// 0010   00 86 d2 c0 40 00 40 11 51 52 c0 a8 cb 01 c0 a8  ....@.@.QR......
// 0020   ca 01 b0 5d 12 b5 00 72 00 00 08 00 00 00 00 00  ...]...r........
// 0030   00 00 00 30 88 01 00 02 00 16 3e 37 f6 04 08 00  ...0......>7....
// 0040   45 00 00 54 00 00 40 00 40 01 23 4f c0 a8 cb 03  E..T..@.@.#O....
// 0050   c0 a8 cb 05 08 00 f6 f2 05 0c 00 01 fc e2 97 51  ...............Q
// 0060   00 00 00 00 a6 f8 02 00 00 00 00 00 10 11 12 13  ................
// 0070   14 15 16 17 18 19 1a 1b 1c 1d 1e 1f 20 21 22 23  ............ !"#
// 0080   24 25 26 27 28 29 2a 2b 2c 2d 2e 2f 30 31 32 33  $%&'()*+,-./0123
// 0090   34 35 36 37                                      4567               ./01234567
var testPacketVXLAN = []byte{
	0x00, 0x16, 0x3e, 0x08, 0x71, 0xcf, 0x36, 0xdc, 0x85, 0x1e, 0xb3, 0x40, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x86, 0xd2, 0xc0, 0x40, 0x00, 0x40, 0x11, 0x51, 0x52, 0xc0, 0xa8, 0xcb, 0x01, 0xc0, 0xa8,
	0xca, 0x01, 0xb0, 0x5d, 0x12, 0xb5, 0x00, 0x72, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xff, 0x00, 0x00, 0x30, 0x88, 0x01, 0x00, 0x02, 0x00, 0x16, 0x3e, 0x37, 0xf6, 0x04, 0x08, 0x00,
	0x45, 0x00, 0x00, 0x54, 0x00, 0x00, 0x40, 0x00, 0x40, 0x01, 0x23, 0x4f, 0xc0, 0xa8, 0xcb, 0x03,
	0xc0, 0xa8, 0xcb, 0x05, 0x08, 0x00, 0xf6, 0xf2, 0x05, 0x0c, 0x00, 0x01, 0xfc, 0xe2, 0x97, 0x51,
	0x00, 0x00, 0x00, 0x00, 0xa6, 0xf8, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x11, 0x12, 0x13,
	0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
	0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33,
	0x34, 0x35, 0x36, 0x37,
}

func TestPacketVXLAN(t *testing.T) {
	p := gopacket.NewPacket(testPacketVXLAN, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeUDP, LayerTypeVXLAN, LayerTypeEthernet, LayerTypeIPv4, LayerTypeICMPv4, gopacket.LayerTypePayload}, t)
	if got, ok := p.Layer(LayerTypeVXLAN).(*VXLAN); ok {
		want := &VXLAN{
			BaseLayer: BaseLayer{
				Contents: []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00},
				Payload: []byte{0x00, 0x30, 0x88, 0x01, 0x00, 0x02, 0x00, 0x16, 0x3e, 0x37, 0xf6, 0x04, 0x08, 0x00,
					0x45, 0x00, 0x00, 0x54, 0x00, 0x00, 0x40, 0x00, 0x40, 0x01, 0x23, 0x4f, 0xc0, 0xa8, 0xcb, 0x03,
					0xc0, 0xa8, 0xcb, 0x05, 0x08, 0x00, 0xf6, 0xf2, 0x05, 0x0c, 0x00, 0x01, 0xfc, 0xe2, 0x97, 0x51,
					0x00, 0x00, 0x00, 0x00, 0xa6, 0xf8, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
					0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33,
					0x34, 0x35, 0x36, 0x37},
			},
			ValidIDFlag:      true,
			VNI:              255,
			GBPExtension:     false,
			GBPApplied:       false,
			GBPDontLearn:     false,
			GBPGroupPolicyID: 0,
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("VXLAN layer mismatch, \nwant %#v\ngot %#v\n", want, got)
		}
	}
}

func BenchmarkDecodePacketVXLAN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketVXLAN, LinkTypeEthernet, gopacket.NoCopy)
	}
}
