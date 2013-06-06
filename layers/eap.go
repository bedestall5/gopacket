// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"code.google.com/p/gopacket"
	"encoding/binary"
)

type EAPCode uint8
type EAPType uint8

const (
	EAPCodeRequest  EAPCode = 1
	EAPCodeResponse EAPCode = 2
	EAPCodeSuccess  EAPCode = 3
	EAPCodeFailure  EAPCode = 4

	// EAPTypeNone means that this EAP layer has no Type or TypeData.
	// Success and Failure EAPs will have this set.
	EAPTypeNone EAPType = 0

	EAPTypeIdentity     EAPType = 1
	EAPTypeNotification EAPType = 2
	EAPTypeNACK         EAPType = 3
	EAPTypeOTP          EAPType = 4
	EAPTypeTokenCard    EAPType = 5
)

// EAP defines an Extensible Authentication Protocol (rfc 3748) layer.
type EAP struct {
	BaseLayer
	Code     EAPCode
	Id       uint8
	Length   uint16
	Type     EAPType
	TypeData []byte
}

// LayerType returns LayerTypeEAP.
func (e *EAP) LayerType() gopacket.LayerType { return LayerTypeEAP }

func decodeEAP(data []byte, p gopacket.PacketBuilder) error {
	e := &EAP{
		Code:   EAPCode(data[0]),
		Id:     data[1],
		Length: binary.BigEndian.Uint16(data[2:4]),
	}
	if e.Length > 4 {
		e.Type = EAPType(data[4])
		e.TypeData = data[5:]
	}
	e.BaseLayer.Contents = data[:e.Length]
	e.BaseLayer.Payload = data[e.Length:] // Should be 0 bytes
	p.AddLayer(e)
	// If we have any bytes left in the packet, we have no idea what they are,
	// so treat them as unknown data.
	return p.NextDecoder(gopacket.DecodeUnknown)
}
