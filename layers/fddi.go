// Copyright 2012 Google, Inc. All rights reserved.

package layers

import (
	"github.com/gconnell/gopacket"
)

// FDDI contains the header for FDDI frames.
type FDDI struct {
	baseLayer
	FrameControl   FDDIFrameControl
	Priority       uint8
	SrcMAC, DstMAC []byte
}

func (f *FDDI) LayerType() gopacket.LayerType { return LayerTypeFDDI }

func (f *FDDI) LinkFlow() gopacket.Flow {
	return gopacket.NewFlow(EndpointMAC, f.SrcMAC, f.DstMAC)
}

func decodeFDDI(data []byte, p gopacket.PacketBuilder) error {
	f := &FDDI{
		FrameControl: FDDIFrameControl(data[0] & 0xF8),
		Priority:     data[0] & 0x07,
		SrcMAC:       data[1:7],
		DstMAC:       data[7:13],
		baseLayer:    baseLayer{data[:13], data[13:]},
	}
	p.SetLinkLayer(f)
	p.AddLayer(f)
	return p.NextDecoder(f.FrameControl)
}
