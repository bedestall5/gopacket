// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.
package layers

import (
	"github.com/lrsk/gopacket"
	"reflect"
	"testing"
)

// testPacketGRE is the packet:
//   15:08:08.003196 IP 192.168.1.1 > 192.168.1.2: GREv0, length 88: IP 172.16.1.1 > 172.16.2.1: ICMP echo request, id 4724, seq 1, length 64
//      0x0000:  3a56 6b69 595e 8e7a 12c3 a971 0800 4500  :VkiY^.z...q..E.
//      0x0010:  006c 843c 4000 402f 32d3 c0a8 0101 c0a8  .l.<@.@/2.......
//      0x0020:  0102 0000 0800 4500 0054 0488 4000 4001  ......E..T..@.@.
//      0x0030:  dafe ac10 0101 ac10 0201 0800 82c4 1274  ...............t
//      0x0040:  0001 c892 a354 0000 0000 380c 0000 0000  .....T....8.....
//      0x0050:  0000 1011 1213 1415 1617 1819 1a1b 1c1d  ................
//      0x0060:  1e1f 2021 2223 2425 2627 2829 2a2b 2c2d  ...!"#$%&'()*+,-
//      0x0070:  2e2f 3031 3233 3435 3637                 ./01234567
var testPacketGRE = []byte{
	0x3a, 0x56, 0x6b, 0x69, 0x59, 0x5e, 0x8e, 0x7a, 0x12, 0xc3, 0xa9, 0x71, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x6c, 0x84, 0x3c, 0x40, 0x00, 0x40, 0x2f, 0x32, 0xd3, 0xc0, 0xa8, 0x01, 0x01, 0xc0, 0xa8,
	0x01, 0x02, 0x00, 0x00, 0x08, 0x00, 0x45, 0x00, 0x00, 0x54, 0x04, 0x88, 0x40, 0x00, 0x40, 0x01,
	0xda, 0xfe, 0xac, 0x10, 0x01, 0x01, 0xac, 0x10, 0x02, 0x01, 0x08, 0x00, 0x82, 0xc4, 0x12, 0x74,
	0x00, 0x01, 0xc8, 0x92, 0xa3, 0x54, 0x00, 0x00, 0x00, 0x00, 0x38, 0x0c, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
	0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d,
	0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
}

func TestPacketGRE(t *testing.T) {
	p := gopacket.NewPacket(testPacketGRE, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeGRE, LayerTypeIPv4, LayerTypeICMPv4, gopacket.LayerTypePayload}, t)
	if got, ok := p.Layer(LayerTypeGRE).(*GRE); ok {
		want := &GRE{
			BaseLayer: BaseLayer{testPacketGRE[34:38], testPacketGRE[38:]},
			Protocol:  EthernetTypeIPv4,
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("GRE layer mismatch, \nwant %#v\ngot  %#v\n", want, got)
		}
	}
}

func BenchmarkDecodePacketGRE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketGRE, LinkTypeEthernet, gopacket.NoCopy)
	}
}

// testPacketEthernetOverGRE is the packet:
//   11:01:38.124768 IP 192.168.1.1 > 192.168.1.2: GREv0, length 102: IP 172.16.1.1 > 172.16.1.2: ICMP echo request, id 3842, seq 1, length 64
//      0x0000:  ea6b 4cd3 5513 d6b9 d880 56ef 0800 4500  .kL.U.....V...E.
//      0x0010:  007a 0acd 4000 402f ac34 c0a8 0101 c0a8  .z..@.@/.4......
//      0x0020:  0102 0000 6558 aa6a 36e6 c630 6e32 3ec7  ....eX.j6..0n2>.
//      0x0030:  9def 0800 4500 0054 d970 4000 4001 0715  ....E..T.p@.@...
//      0x0040:  ac10 0101 ac10 0102 0800 3f15 0f02 0001  ..........?.....
//      0x0050:  82d9 b154 0000 0000 b5e6 0100 0000 0000  ...T............
//      0x0060:  1011 1213 1415 1617 1819 1a1b 1c1d 1e1f  ................
//      0x0070:  2021 2223 2425 2627 2829 2a2b 2c2d 2e2f  .!"#$%&'()*+,-./
//      0x0080:  3031 3233 3435 3637                      01234567
var testPacketEthernetOverGRE = []byte{
	0xea, 0x6b, 0x4c, 0xd3, 0x55, 0x13, 0xd6, 0xb9, 0xd8, 0x80, 0x56, 0xef, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x7a, 0x0a, 0xcd, 0x40, 0x00, 0x40, 0x2f, 0xac, 0x34, 0xc0, 0xa8, 0x01, 0x01, 0xc0, 0xa8,
	0x01, 0x02, 0x00, 0x00, 0x65, 0x58, 0xaa, 0x6a, 0x36, 0xe6, 0xc6, 0x30, 0x6e, 0x32, 0x3e, 0xc7,
	0x9d, 0xef, 0x08, 0x00, 0x45, 0x00, 0x00, 0x54, 0xd9, 0x70, 0x40, 0x00, 0x40, 0x01, 0x07, 0x15,
	0xac, 0x10, 0x01, 0x01, 0xac, 0x10, 0x01, 0x02, 0x08, 0x00, 0x3f, 0x15, 0x0f, 0x02, 0x00, 0x01,
	0x82, 0xd9, 0xb1, 0x54, 0x00, 0x00, 0x00, 0x00, 0xb5, 0xe6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
}

func TestPacketEthernetOverGRE(t *testing.T) {
	p := gopacket.NewPacket(testPacketEthernetOverGRE, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeGRE, LayerTypeEthernet, LayerTypeIPv4, LayerTypeICMPv4, gopacket.LayerTypePayload}, t)
	if got, ok := p.Layer(LayerTypeGRE).(*GRE); ok {
		want := &GRE{
			BaseLayer: BaseLayer{testPacketEthernetOverGRE[34:38], testPacketEthernetOverGRE[38:]},
			Protocol:  EthernetTypeTransparentEthernetBridging,
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("GRE layer mismatch, \nwant %#v\ngot  %#v\n", want, got)
		}
	}
}

func BenchmarkDecodePacketEthernetOverGRE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketEthernetOverGRE, LinkTypeEthernet, gopacket.NoCopy)
	}
}
