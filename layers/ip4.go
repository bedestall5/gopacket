// Copyright 2012 Google, Inc. All rights reserved.
// Copyright 2009-2011 Andreas Krennmair. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"code.google.com/p/gopacket"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type IPv4Flag uint8

const (
	IPv4MoreFragments IPv4Flag = 1 << 0
	IPv4DontFragment  IPv4Flag = 1 << 1
	IPv4EvilBit       IPv4Flag = 1 << 2 // http://tools.ietf.org/html/rfc3514 ;) (As an April Fools' joke, proposed for use in RFC 3514 as the "Evil bit".)
)

func (f IPv4Flag) String() string {
	var s []string

	if f&IPv4EvilBit != 0 {
		s = append(s, "Evil")
	}
	if f&IPv4DontFragment != 0 {
		s = append(s, "DF")
	}
	if f&IPv4MoreFragments != 0 {
		s = append(s, "MF")
	}
	return strings.Join(s, "|")
}

// IPv4 is the header of an IP packet.
type IPv4 struct {
	BaseLayer
	Version    uint8
	IHL        uint8
	TOS        uint8
	Length     uint16
	Id         uint16
	Flags      IPv4Flag
	FragOffset uint16
	TTL        uint8
	Protocol   IPProtocol
	Checksum   uint16
	SrcIP      net.IP
	DstIP      net.IP
	Options    []IPv4Option
	Padding    []byte
}

// LayerType returns LayerTypeIPv4
func (i *IPv4) LayerType() gopacket.LayerType { return LayerTypeIPv4 }
func (i *IPv4) NetworkFlow() gopacket.Flow {
	return gopacket.NewFlow(EndpointIPv4, i.SrcIP, i.DstIP)
}

type IPv4Option struct {
	OptionType   uint8
	OptionLength uint8
	OptionData   []byte
}

func (i IPv4Option) String() string {
	return fmt.Sprintf("IPv4Option(%v:%v)", i.OptionType, i.OptionData)
}

// SerializeTo writes the serialized form of this layer into the
// SerializationBuffer, implementing gopacket.SerializableLayer.
func (ip *IPv4) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	if len(ip.Options) > 0 {
		return fmt.Errorf("cannot currently serialize IPv4 options")
	}
	bytes, err := b.PrependBytes(20)
	if err != nil {
		return err
	}
	if opts.FixLengths {
		ip.IHL = 5 // Fix when we add support for options.
		ip.Length = uint16(len(b.Bytes()))
	}
	bytes[0] = (ip.Version << 4) | ip.IHL
	bytes[1] = ip.TOS
	binary.BigEndian.PutUint16(bytes[2:], ip.Length)
	binary.BigEndian.PutUint16(bytes[4:], ip.Id)
	binary.BigEndian.PutUint16(bytes[6:], ip.flagsfrags())
	bytes[8] = ip.TTL
	bytes[9] = byte(ip.Protocol)
	if len(ip.SrcIP) != 4 {
		return fmt.Errorf("invalid src IP %v", ip.SrcIP)
	}
	if len(ip.DstIP) != 4 {
		return fmt.Errorf("invalid dst IP %v", ip.DstIP)
	}
	copy(bytes[12:16], ip.SrcIP)
	copy(bytes[16:20], ip.DstIP)
	if opts.ComputeChecksums {
		// Clear checksum bytes
		bytes[10] = 0
		bytes[11] = 0
		// Compute checksum
		var csum uint32
		for i := 0; i < len(bytes); i += 2 {
			csum += uint32(bytes[i]) << 8
			csum += uint32(bytes[i+1])
		}
		ip.Checksum = ^uint16((csum >> 16) + csum)
	}
	binary.BigEndian.PutUint16(bytes[10:], ip.Checksum)
	return nil
}

func (ip *IPv4) flagsfrags() (ff uint16) {
	ff |= uint16(ip.Flags) << 13
	ff |= ip.FragOffset
	return
}

// DecodeFromBytes decodes the given bytes into this layer.
func (ip *IPv4) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	flagsfrags := binary.BigEndian.Uint16(data[6:8])

	ip.Version = uint8(data[0]) >> 4
	ip.IHL = uint8(data[0]) & 0x0F
	ip.TOS = data[1]
	ip.Length = binary.BigEndian.Uint16(data[2:4])
	ip.Id = binary.BigEndian.Uint16(data[4:6])
	ip.Flags = IPv4Flag(flagsfrags >> 13)
	ip.FragOffset = flagsfrags & 0x1FFF
	ip.TTL = data[8]
	ip.Protocol = IPProtocol(data[9])
	ip.Checksum = binary.BigEndian.Uint16(data[10:12])
	ip.SrcIP = data[12:16]
	ip.DstIP = data[16:20]
	// Set up an initial guess for contents/payload... we'll reset these soon.
	ip.BaseLayer = BaseLayer{Contents: data}

	if ip.Length < 20 {
		return fmt.Errorf("Invalid (too small) IP length (%d < 20)", ip.Length)
	} else if ip.IHL < 5 {
		return fmt.Errorf("Invalid (too small) IP header length (%d < 5)", ip.IHL)
	} else if int(ip.IHL*4) > int(ip.Length) {
		return fmt.Errorf("Invalid IP header length > IP length (%d > %d)", ip.IHL, ip.Length)
	}
	if cmp := len(data) - int(ip.Length); cmp > 0 {
		data = data[:ip.Length]
	} else if cmp < 0 {
		df.SetTruncated()
		if int(ip.IHL)*4 > len(data) {
			return fmt.Errorf("Not all IP header bytes available")
		}
	}
	ip.Contents = data[:ip.IHL*4]
	ip.Payload = data[ip.IHL*4:]
	// From here on, data contains the header options.
	data = data[20 : ip.IHL*4]
	// Pull out IP options
	for len(data) > 0 {
		if ip.Options == nil {
			// Pre-allocate to avoid growing the slice too much.
			ip.Options = make([]IPv4Option, 0, 4)
		}
		opt := IPv4Option{OptionType: data[0]}
		switch opt.OptionType {
		case 0: // End of options
			opt.OptionLength = 1
			ip.Options = append(ip.Options, opt)
			ip.Padding = data[1:]
			break
		case 1: // 1 byte padding
			opt.OptionLength = 1
		default:
			opt.OptionLength = data[1]
			opt.OptionData = data[2:opt.OptionLength]
		}
		if len(data) >= int(opt.OptionLength) {
			data = data[opt.OptionLength:]
		} else {
			return fmt.Errorf("IP option length exceeds remaining IP header size, option type %v length %v", opt.OptionType, opt.OptionLength)
		}
		ip.Options = append(ip.Options, opt)
	}
	return nil
}

func (i *IPv4) CanDecode() gopacket.LayerClass {
	return LayerTypeIPv4
}

func (i *IPv4) NextLayerType() gopacket.LayerType {
	return i.Protocol.LayerType()
}

func decodeIPv4(data []byte, p gopacket.PacketBuilder) error {
	ip := &IPv4{}
	err := ip.DecodeFromBytes(data, p)
	p.AddLayer(ip)
	p.SetNetworkLayer(ip)
	if err != nil {
		return err
	}
	return p.NextDecoder(ip.Protocol)
}
