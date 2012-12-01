// Copyright (c) 2012 Google, Inc. All rights reserved.
// Copyright (c) 2009-2012 Andreas Krennmair. All rights reserved.

package gopacket

import (
	"encoding/binary"
)

// IPv6 is the layer for the IPv6 header.
type IPv6 struct {
	// http://www.networksorcery.com/enp/protocol/ipv6.htm
	Version      uint8      // 4 bits
	TrafficClass uint8      // 8 bits
	FlowLabel    uint32     // 20 bits
	Length       uint16     // 16 bits
	NextHeader   IPProtocol // 8 bits, same as Protocol in Iphdr
	HopLimit     uint8      // 8 bits
	SrcIP        IPAddress  // 16 bytes
	DstIP        IPAddress  // 16 bytes
}

// LayerType returns LayerTypeIPv6
func (i *IPv6) LayerType() LayerType { return LayerTypeIPv6 }
func (i *IPv6) SrcNetAddr() Address  { return i.SrcIP }
func (i *IPv6) DstNetAddr() Address  { return i.DstIP }

var decodeIPv6 decoderFunc = func(data []byte) (out DecodeResult, err error) {
	ip6 := &IPv6{
		Version:      uint8(data[0]) >> 4,
		TrafficClass: uint8((binary.BigEndian.Uint16(data[0:2]) >> 4) & 0x00FF),
		FlowLabel:    binary.BigEndian.Uint32(data[0:4]) & 0x000FFFFF,
		Length:       binary.BigEndian.Uint16(data[4:6]),
		NextHeader:   IPProtocol(data[6]),
		HopLimit:     data[7],
		SrcIP:        data[8:24],
		DstIP:        data[24:40],
	}
	out.DecodedLayer = ip6
	out.RemainingBytes = data[40:]
	out.NextDecoder = ip6.NextHeader
	out.NetworkLayer = ip6
	return
}
