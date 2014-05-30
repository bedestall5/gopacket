// Copyright 2012, Google, Inc. All rights reserved.
// Copyright 2009-2011 Andreas Krennmair. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"reflect"
	"testing"

	"code.google.com/p/gopacket"
)

// testUDPPacketDNS is the packet:
//   10:33:07.883637 IP 172.16.255.1.53 > 172.29.20.15.35181: 47320 7/0/0 MX ALT2.ASPMX.L.GOOGLE.com. 20, MX ASPMX2.GOOGLEMAIL.com. 30, MX ASPMX3.GOOGLEMAIL.com. 30, MX ASPMX4.GOOGLEMAIL.com. 30, MX ASPMX5.GOOGLEMAIL.com. 30, MX ASPMX.L.GOOGLE.com. 10, MX ALT1.ASPMX.L.GOOGLE.com. 20 (202)
//      0x0000:  24be 0527 0b17 001f cab3 75c0 0800 4500  $..'......u...E.
//      0x0010:  00e6 68cf 0000 3f11 a6f9 ac10 ff01 ac1d  ..h...?.........
//      0x0020:  140f 0035 896d 00d2 754a b8d8 8180 0001  ...5.m..uJ......
//      0x0030:  0007 0000 0000 0478 6b63 6403 636f 6d00  .......xkcd.com.
//      0x0040:  000f 0001 c00c 000f 0001 0000 0258 0018  .............X..
//      0x0050:  0014 0441 4c54 3205 4153 504d 5801 4c06  ...ALT2.ASPMX.L.
//      0x0060:  474f 4f47 4c45 c011 c00c 000f 0001 0000  GOOGLE..........
//      0x0070:  0258 0016 001e 0641 5350 4d58 320a 474f  .X.....ASPMX2.GO
//      0x0080:  4f47 4c45 4d41 494c c011 c00c 000f 0001  OGLEMAIL........
//      0x0090:  0000 0258 000b 001e 0641 5350 4d58 33c0  ...X.....ASPMX3.
//      0x00a0:  53c0 0c00 0f00 0100 0002 5800 0b00 1e06  S.........X.....
//      0x00b0:  4153 504d 5834 c053 c00c 000f 0001 0000  ASPMX4.S........
//      0x00c0:  0258 000b 001e 0641 5350 4d58 35c0 53c0  .X.....ASPMX5.S.
//      0x00d0:  0c00 0f00 0100 0002 5800 0400 0ac0 2dc0  ........X.....-.
//      0x00e0:  0c00 0f00 0100 0002 5800 0900 1404 414c  ........X.....AL
//      0x00f0:  5431 c02d                                T1.-
// Packet generated by doing DNS query for 'xkcd.com'
var testUDPPacketDNS = []byte{
	0x24, 0xbe, 0x05, 0x27, 0x0b, 0x17, 0x00, 0x1f, 0xca, 0xb3, 0x75, 0xc0, 0x08, 0x00, 0x45, 0x00,
	0x00, 0xe6, 0x68, 0xcf, 0x00, 0x00, 0x3f, 0x11, 0xa6, 0xf9, 0xac, 0x10, 0xff, 0x01, 0xac, 0x1d,
	0x14, 0x0f, 0x00, 0x35, 0x89, 0x6d, 0x00, 0xd2, 0x75, 0x4a, 0xb8, 0xd8, 0x81, 0x80, 0x00, 0x01,
	0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x04, 0x78, 0x6b, 0x63, 0x64, 0x03, 0x63, 0x6f, 0x6d, 0x00,
	0x00, 0x0f, 0x00, 0x01, 0xc0, 0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00, 0x02, 0x58, 0x00, 0x18,
	0x00, 0x14, 0x04, 0x41, 0x4c, 0x54, 0x32, 0x05, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x01, 0x4c, 0x06,
	0x47, 0x4f, 0x4f, 0x47, 0x4c, 0x45, 0xc0, 0x11, 0xc0, 0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00,
	0x02, 0x58, 0x00, 0x16, 0x00, 0x1e, 0x06, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x32, 0x0a, 0x47, 0x4f,
	0x4f, 0x47, 0x4c, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0xc0, 0x11, 0xc0, 0x0c, 0x00, 0x0f, 0x00, 0x01,
	0x00, 0x00, 0x02, 0x58, 0x00, 0x0b, 0x00, 0x1e, 0x06, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x33, 0xc0,
	0x53, 0xc0, 0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00, 0x02, 0x58, 0x00, 0x0b, 0x00, 0x1e, 0x06,
	0x41, 0x53, 0x50, 0x4d, 0x58, 0x34, 0xc0, 0x53, 0xc0, 0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00,
	0x02, 0x58, 0x00, 0x0b, 0x00, 0x1e, 0x06, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x35, 0xc0, 0x53, 0xc0,
	0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00, 0x02, 0x58, 0x00, 0x04, 0x00, 0x0a, 0xc0, 0x2d, 0xc0,
	0x0c, 0x00, 0x0f, 0x00, 0x01, 0x00, 0x00, 0x02, 0x58, 0x00, 0x09, 0x00, 0x14, 0x04, 0x41, 0x4c,
	0x54, 0x31, 0xc0, 0x2d,
}

func TestUDPPacketDNS(t *testing.T) {
	p := gopacket.NewPacket(testUDPPacketDNS, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeUDP, LayerTypeDNS}, t)
	if got, ok := p.TransportLayer().(*UDP); ok {
		want := &UDP{
			BaseLayer: BaseLayer{
				Contents: []byte{0x0, 0x35, 0x89, 0x6d, 0x0, 0xd2, 0x75, 0x4a},
				Payload: []byte{0xb8, 0xd8, 0x81, 0x80, 0x0, 0x1, 0x0,
					0x7, 0x0, 0x0, 0x0, 0x0, 0x4, 0x78, 0x6b, 0x63, 0x64, 0x3, 0x63, 0x6f,
					0x6d, 0x0, 0x0, 0xf, 0x0, 0x1, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1, 0x0, 0x0,
					0x2, 0x58, 0x0, 0x18, 0x0, 0x14, 0x4, 0x41, 0x4c, 0x54, 0x32, 0x5, 0x41,
					0x53, 0x50, 0x4d, 0x58, 0x1, 0x4c, 0x6, 0x47, 0x4f, 0x4f, 0x47, 0x4c,
					0x45, 0xc0, 0x11, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0,
					0x16, 0x0, 0x1e, 0x6, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x32, 0xa, 0x47, 0x4f,
					0x4f, 0x47, 0x4c, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0xc0, 0x11, 0xc0, 0xc,
					0x0, 0xf, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0xb, 0x0, 0x1e, 0x6, 0x41,
					0x53, 0x50, 0x4d, 0x58, 0x33, 0xc0, 0x53, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1,
					0x0, 0x0, 0x2, 0x58, 0x0, 0xb, 0x0, 0x1e, 0x6, 0x41, 0x53, 0x50, 0x4d,
					0x58, 0x34, 0xc0, 0x53, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1, 0x0, 0x0, 0x2,
					0x58, 0x0, 0xb, 0x0, 0x1e, 0x6, 0x41, 0x53, 0x50, 0x4d, 0x58, 0x35, 0xc0,
					0x53, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0x4, 0x0,
					0xa, 0xc0, 0x2d, 0xc0, 0xc, 0x0, 0xf, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0,
					0x9, 0x0, 0x14, 0x4, 0x41, 0x4c, 0x54, 0x31, 0xc0, 0x2d},
			},
			SrcPort:  53,
			DstPort:  35181,
			Length:   210,
			Checksum: 30026,
			sPort:    []byte{0x0, 0x35},
			dPort:    []byte{0x89, 0x6d},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("UDP packet mismatch:\ngot  :\n%#v\n\nwant :\n%#v\n\n", got, want)
		}
	} else {
		t.Error("Transport layer packet not UDP")
	}
}

func loadDNS(dnspacket []byte, t *testing.T) *DNS {
	p := gopacket.NewPacket(dnspacket, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4,
		LayerTypeUDP, LayerTypeDNS}, t)

	dnsL := p.Layer(LayerTypeDNS)
	if dnsL == nil {
		t.Error("No DNS Layer found")
	}

	dns, ok := dnsL.(*DNS)
	if !ok {
		return nil
	}
	return dns
}

var testDNSQueryA = []byte{
	0xfe, 0x54, 0x00, 0x3e, 0x00, 0x96, 0x52, 0x54, /* .T.>..RT */
	0x00, 0xbd, 0x1c, 0x70, 0x08, 0x00, 0x45, 0x00, /* ...p..E. */
	0x00, 0x3c, 0x22, 0xe0, 0x00, 0x00, 0x40, 0x11, /* .<"...@. */
	0xe2, 0x38, 0xc0, 0xa8, 0x7a, 0x46, 0xc0, 0xa8, /* .8..zF.. */
	0x7a, 0x01, 0xc3, 0x35, 0x00, 0x35, 0x00, 0x28, /* z..5.5.( */
	0x75, 0xd2, 0x52, 0x41, 0x01, 0x00, 0x00, 0x01, /* u.RA.... */
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, /* .......w */
	0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, /* ww.googl */
	0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, /* e.com... */
	0x00, 0x01, /* .. */
}

func TestDNSQueryA(t *testing.T) {
	dns := loadDNS(testDNSQueryA, t)
	if dns == nil {
		t.Error("Failed to get a pointer to DNS struct")
		return
	}

	if uint16(len(dns.Questions)) != dns.Header.QDCount {
		t.Error("Invalid query decoding, not the right number of questions")
		return
	}

	if dns.Questions[0].Name != "www.google.com" {
		t.Errorf("Invalid query decoding, expecting 'www.google.com', got '%s'",
			dns.Questions[0].Name)
	}
	if dns.Questions[0].Class != DNSClassIN {
		t.Errorf("Invalid query decoding, expecting Class IN, got '%d'",
			dns.Questions[0].Class)
	}

	if dns.Questions[0].Type != DNSTypeA {
		t.Errorf("Invalid query decoding, expecting Type A, got '%d'",
			dns.Questions[0].Type)
	}
}

var testDNSRRA = []byte{
	0x52, 0x54, 0x00, 0xbd, 0x1c, 0x70, 0xfe, 0x54, /* RT...p.T */
	0x00, 0x3e, 0x00, 0x96, 0x08, 0x00, 0x45, 0x00, /* .>....E. */
	0x01, 0x24, 0x00, 0x00, 0x40, 0x00, 0x40, 0x11, /* .$..@.@. */
	0xc4, 0x30, 0xc0, 0xa8, 0x7a, 0x01, 0xc0, 0xa8, /* .0..z... */
	0x7a, 0x46, 0x00, 0x35, 0xc3, 0x35, 0x01, 0x10, /* zF.5.5.. */
	0x76, 0xba, 0x52, 0x41, 0x81, 0x80, 0x00, 0x01, /* v.RA.... */
	0x00, 0x06, 0x00, 0x04, 0x00, 0x04, 0x03, 0x77, /* .......w */
	0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, /* ww.googl */
	0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, /* e.com... */
	0x00, 0x01, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* ........ */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x67, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* .g...... */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x68, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* .h...... */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x69, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* .i...... */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x6a, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* .j...... */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x93, 0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, /* ........ */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x04, 0x4a, 0x7d, /* ...,..J} */
	0xc3, 0x63, 0xc0, 0x10, 0x00, 0x02, 0x00, 0x01, /* .c...... */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x06, 0x03, 0x6e, /* .......n */
	0x73, 0x32, 0xc0, 0x10, 0xc0, 0x10, 0x00, 0x02, /* s2...... */
	0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, 0x00, 0x06, /* ........ */
	0x03, 0x6e, 0x73, 0x33, 0xc0, 0x10, 0xc0, 0x10, /* .ns3.... */
	0x00, 0x02, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x06, 0x03, 0x6e, 0x73, 0x31, 0xc0, 0x10, /* ...ns1.. */
	0xc0, 0x10, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02, /* ........ */
	0xa3, 0x00, 0x00, 0x06, 0x03, 0x6e, 0x73, 0x34, /* .....ns4 */
	0xc0, 0x10, 0xc0, 0xb0, 0x00, 0x01, 0x00, 0x01, /* ........ */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x04, 0xd8, 0xef, /* ........ */
	0x20, 0x0a, 0xc0, 0x8c, 0x00, 0x01, 0x00, 0x01, /*  ....... */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x04, 0xd8, 0xef, /* ........ */
	0x22, 0x0a, 0xc0, 0x9e, 0x00, 0x01, 0x00, 0x01, /* "....... */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x04, 0xd8, 0xef, /* ........ */
	0x24, 0x0a, 0xc0, 0xc2, 0x00, 0x01, 0x00, 0x01, /* $....... */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x04, 0xd8, 0xef, /* ........ */
	0x26, 0x0a, /* &. */
}

func TestDNSRRA(t *testing.T) {
	dns := loadDNS(testDNSRRA, t)
	if dns == nil {
		t.Error("Failed to get a pointer to DNS struct")
		return
	}

	if uint16(len(dns.Questions)) != dns.Header.QDCount {
		t.Error("Invalid query decoding, not the right number of questions")
		return
	}
	if uint16(len(dns.Answers)) != dns.Header.ANCount {
		t.Error("Invalid query decoding, not the right number of answers")
		return
	}
	if uint16(len(dns.Authorities)) != dns.Header.NSCount {
		t.Error("Invalid query decoding, not the right number of authorities")
		return
	}
	if uint16(len(dns.Additionals)) != dns.Header.ARCount {
		t.Error("Invalid query decoding, not the right number of additionals info")
		return
	}

	if dns.Questions[0].Name != "www.google.com" {
		t.Errorf("Invalid query decoding, expecting 'www.google.com', got '%s'",
			dns.Questions[0].Name)
	}
	if dns.Answers[0].Name != "www.google.com" {
		t.Errorf("Invalid query decoding, expecting 'www.google.com', got '%d'",
			dns.Questions[0].Class)
	}
	if dns.Answers[0].Class != DNSClassIN {
		t.Errorf("Invalid query decoding, expecting Class IN, got '%d'",
			dns.Questions[0].Class)
	}
	if dns.Answers[0].Type != DNSTypeA {
		t.Errorf("Invalid query decoding, expecting Type A, got '%d'",
			dns.Questions[0].Type)
	}
	if !dns.Answers[0].IP.Equal([]byte{74, 125, 195, 103}) {
		t.Errorf("Invalid query decoding, invalid IP address,"+
			" expecting '74.125.195.103', got '%s'",
			dns.Answers[0].IP.String())
	}
	if len(dns.Answers) != 6 {
		t.Errorf("No correct number of answers, expecting 6, go '%d'",
			len(dns.Answers))
	}
	if len(dns.Authorities) != 4 {
		t.Errorf("No correct number of answers, expecting 4, go '%d'",
			len(dns.Answers))
	}
	if len(dns.Additionals) != 4 {
		t.Errorf("No correct number of answers, expecting 4, go '%d'",
			len(dns.Answers))
	}

}

var testDNSAAAA = []byte{
	0x52, 0x54, 0x00, 0xbd, 0x1c, 0x70, 0xfe, 0x54, /* RT...p.T */
	0x00, 0x3e, 0x00, 0x96, 0x08, 0x00, 0x45, 0x00, /* .>....E. */
	0x00, 0xe0, 0x00, 0x00, 0x40, 0x00, 0x40, 0x11, /* ....@.@. */
	0xc4, 0x74, 0xc0, 0xa8, 0x7a, 0x01, 0xc0, 0xa8, /* .t..z... */
	0x7a, 0x46, 0x00, 0x35, 0xdb, 0x13, 0x00, 0xcc, /* zF.5.... */
	0x76, 0x76, 0xf3, 0x03, 0x81, 0x80, 0x00, 0x01, /* vv...... */
	0x00, 0x01, 0x00, 0x04, 0x00, 0x04, 0x03, 0x77, /* .......w */
	0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, /* ww.googl */
	0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x1c, /* e.com... */
	0x00, 0x01, 0xc0, 0x0c, 0x00, 0x1c, 0x00, 0x01, /* ........ */
	0x00, 0x00, 0x01, 0x2c, 0x00, 0x10, 0x2a, 0x00, /* ...,..*. */
	0x14, 0x50, 0x40, 0x0c, 0x0c, 0x01, 0x00, 0x00, /* .P@..... */
	0x00, 0x00, 0x00, 0x00, 0x00, 0x69, 0xc0, 0x10, /* .....i.. */
	0x00, 0x02, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x06, 0x03, 0x6e, 0x73, 0x34, 0xc0, 0x10, /* ...ns4.. */
	0xc0, 0x10, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02, /* ........ */
	0xa3, 0x00, 0x00, 0x06, 0x03, 0x6e, 0x73, 0x32, /* .....ns2 */
	0xc0, 0x10, 0xc0, 0x10, 0x00, 0x02, 0x00, 0x01, /* ........ */
	0x00, 0x02, 0xa3, 0x00, 0x00, 0x06, 0x03, 0x6e, /* .......n */
	0x73, 0x31, 0xc0, 0x10, 0xc0, 0x10, 0x00, 0x02, /* s1...... */
	0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, 0x00, 0x06, /* ........ */
	0x03, 0x6e, 0x73, 0x33, 0xc0, 0x10, 0xc0, 0x6c, /* .ns3...l */
	0x00, 0x01, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x04, 0xd8, 0xef, 0x20, 0x0a, 0xc0, 0x5a, /* .... ..Z */
	0x00, 0x01, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x04, 0xd8, 0xef, 0x22, 0x0a, 0xc0, 0x7e, /* ...."..~ */
	0x00, 0x01, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x04, 0xd8, 0xef, 0x24, 0x0a, 0xc0, 0x48, /* ....$..H */
	0x00, 0x01, 0x00, 0x01, 0x00, 0x02, 0xa3, 0x00, /* ........ */
	0x00, 0x04, 0xd8, 0xef, 0x26, 0x0a, /* ....&. */
}

func TestDNSAAAA(t *testing.T) {
	dns := loadDNS(testDNSAAAA, t)
	if dns == nil {
		t.Error("Failed to get a pointer to DNS struct")
		return
	}

	if len(dns.Questions) != 1 {
		t.Error("Invalid number of question")
		return
	}
	if dns.Questions[0].Type != DNSTypeAAAA {
		t.Error("Invalid question, Type is not AAAA, found %d",
			dns.Questions[0].Type)
	}

	if len(dns.Answers) != 1 {
		t.Error("Invalid number of answers")
	}
	if !dns.Answers[0].IP.Equal([]byte{0x2a, 0x00, 0x14, 0x50, 0x40,
		0x0c, 0x0c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x69}) {
		t.Error("Invalid IP address, found ",
			dns.Answers[0].IP.String())
	}
}

var testDNSMXSOA = []byte{
	0x52, 0x54, 0x00, 0xbd, 0x1c, 0x70, 0xfe, 0x54, /* RT...p.T */
	0x00, 0x3e, 0x00, 0x96, 0x08, 0x00, 0x45, 0x00, /* .>....E. */
	0x00, 0x6e, 0x00, 0x00, 0x40, 0x00, 0x40, 0x11, /* .n..@.@. */
	0xc4, 0xe6, 0xc0, 0xa8, 0x7a, 0x01, 0xc0, 0xa8, /* ....z... */
	0x7a, 0x46, 0x00, 0x35, 0x9c, 0x60, 0x00, 0x5a, /* zF.5.`.Z */
	0x76, 0x04, 0xfc, 0x7a, 0x81, 0x80, 0x00, 0x01, /* v..z.... */
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x03, 0x77, /* .......w */
	0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, /* ww.googl */
	0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x0f, /* e.com... */
	0x00, 0x01, 0xc0, 0x10, 0x00, 0x06, 0x00, 0x01, /* ........ */
	0x00, 0x00, 0x00, 0x3c, 0x00, 0x26, 0x03, 0x6e, /* ...<.&.n */
	0x73, 0x31, 0xc0, 0x10, 0x09, 0x64, 0x6e, 0x73, /* s1...dns */
	0x2d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0xc0, 0x10, /* -admin.. */
	0x00, 0x17, 0x9f, 0x64, 0x00, 0x00, 0x1c, 0x20, /* ...d...  */
	0x00, 0x00, 0x07, 0x08, 0x00, 0x12, 0x75, 0x00, /* ......u. */
	0x00, 0x00, 0x01, 0x2c, /* ..., */
}

func TestDNSMXSOA(t *testing.T) {
	dns := loadDNS(testDNSMXSOA, t)
	if dns == nil {
		t.Error("Failed to get a pointer to DNS struct")
		return
	}

	if len(dns.Authorities) != 1 {
		t.Error("Invalid number of authoritative answers, found ",
			len(dns.Authorities))
		return
	}
}

func BenchmarkDecodeDNS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testDNSQueryA, LinkTypeEthernet, gopacket.NoCopy)
	}
}
