// Copyright 2012 Google, Inc. All rights reserved.

package gopacket

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// CaptureInfo contains capture metadata for a packet.  If a packet was captured
// off the wire or read from a pcap file (see the 'pcap' subdirectory), this
// information will be attached to the packet.
type CaptureInfo struct {
	// Populated is set to true if the rest of the CaptureInfo has been populated
	// with actual information.  If Populated is false, there's no point in
	// reading any of the other fields.
	Populated             bool
	Timestamp             time.Time
	CaptureLength, Length int
}

// Packet is the primary object used by gopacket.  Packets are created by a
// Decoder's Decode call.  A packet is made up of a set of Data, which
// is broken into a number of Layers as it is decoded.
type Packet interface {
	// Data returns all data associated with this packet
	Data() []byte
	// Layers returns all layers in this packet, computing them as necessary
	Layers() []Layer
	// Layer returns the first layer in this packet of the given type, or nil
	Layer(LayerType) Layer
	// LayerClass returns the first layer in this packet of the given class,
	// or nil.
	LayerClass(LayerClass) Layer
	// String returns a human-readable string.
	String() string
	// CaptureInfo returns the caputure information for this packet.  This returns
	// a pointer to the packet's struct, so it can be used both for reading and
	// writing the information.
	CaptureInfo() *CaptureInfo

	// LinkLayer returns the first link layer in the packet
	LinkLayer() LinkLayer
	// NetworkLayer returns the first network layer in the packet
	NetworkLayer() NetworkLayer
	// TransportLayer returns the first transport layer in the packet
	TransportLayer() TransportLayer
	// ApplicationLayer returns the first application layer in the packet
	ApplicationLayer() ApplicationLayer
	// ErrorLayer is particularly useful, since it returns nil if the packet
	// was fully decoded successfully, and non-nil if an error was encountered
	// in decoding and the packet was only partially decoded.  Thus, its output
	// can be used to determine if the entire packet was able to be decoded.
	ErrorLayer() ErrorLayer
}

type packet struct {
	// data contains the entire packet data for a packet
	data []byte
	// encoded contains all the packet data we have yet to decode
	encoded []byte
	// layers contains each layer we've already decoded
	layers []Layer
	// decoder is the next decoder we should call (lazily)
	decoder Decoder
	// capInfo is the CaptureInfo for this packet
	capInfo CaptureInfo

	// Pointers to the various important layers
	link        LinkLayer
	network     NetworkLayer
	transport   TransportLayer
	application ApplicationLayer
	failure     ErrorLayer
}

func (p *packet) CaptureInfo() *CaptureInfo {
	return &p.capInfo
}

func (p *packet) copySpecificLayersFrom(r *DecodeResult) {
	if p.link == nil {
		p.link = r.LinkLayer
	}
	if p.network == nil {
		p.network = r.NetworkLayer
	}
	if p.transport == nil {
		p.transport = r.TransportLayer
	}
	if p.application == nil {
		p.application = r.ApplicationLayer
	}
	if p.failure == nil {
		p.failure = r.ErrorLayer
	}
}

func (p *packet) LinkLayer() LinkLayer {
	return p.link
}
func (p *packet) NetworkLayer() NetworkLayer {
	return p.network
}
func (p *packet) TransportLayer() TransportLayer {
	return p.transport
}
func (p *packet) ApplicationLayer() ApplicationLayer {
	return p.application
}
func (p *packet) ErrorLayer() ErrorLayer {
	return p.failure
}

func (p *packet) Data() []byte {
	return p.data
}

func (p *packet) appendLayer(l Layer) {
	p.layers = append(p.layers, l)
}

// DecodeOptions tells gopacket how to decode a packet.
type DecodeOptions struct {
	// Lazy decoding decodes the minimum number of layers needed to return data
	// for a packet at each function call.  Be careful using this with concurrent
	// packet processors, as each call to packet.* could mutate the packet, and
	// two concurrent function calls could interact poorly.
	Lazy bool
	// NoCopy decoding doesn't copy its input buffer into storage that's owned by
	// the packet.  If you can guarantee that the bytes underlying the slice
	// passed into NewPacket aren't going to be modified, this can be faster.  If
	// there's any chance that those bytes WILL be changed, this will invalidate
	// your packets.
	NoCopy bool
}

// Default decoding provides the safest (but slowest) method for decoding
// packets.  It eagerly processes all layers (so it's concurrency-safe) and it
// copies its input buffer upon creation of the packet (so the packet remains
// valid if the underlying slice is modified.  Both of these take time,
// though, so beware.  If you can guarantee that the packet will only be used
// by one goroutine at a time, set Lazy decoding.  If you can guarantee that
// the underlying slice won't change, set NoCopy decoding.
var Default DecodeOptions = DecodeOptions{}

// Lazy is a DecodeOptions with just Lazy set.
var Lazy DecodeOptions = DecodeOptions{Lazy: true}

// NoCopy is a DecodeOptions with just NoCopy set.
var NoCopy DecodeOptions = DecodeOptions{NoCopy: true}

// NewPacket creates a new Packet object from a set of bytes.  The
// firstLayerDecoder tells it how to interpret the first layer from the bytes,
// future layers will be generated from that first layer automatically.
func NewPacket(data []byte, firstLayerDecoder Decoder, options DecodeOptions) Packet {
	if !options.NoCopy {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		data = dataCopy
	}
	p := &packet{
		data: data,
		// We start off with a size-4 slice since growing a size-zero slice actually
		// can take us a large amount of time, and we expect most packets to give us
		// 4 layers (link, network, transport, application).  This gives our 4-layer
		// benchmark (DecodeNotLazy) a speedup of ~10% (2150ns -> 1922ns)
		layers: make([]Layer, 0, 4),
	}
	if !options.Lazy {
		err := firstLayerDecoder.Decode(data, eagerCollector(p.layers))
		if err != nil {
			errData := data
			if len(p.layers) > 0 {
				errData = p.layers[len(p.layers)-1].LayerPayload()
			}
			p.layers = append(p.layers, &DecodeFailure{
				data: errData,
				err:  err,
			})
		}
		p.Layers()
	} else {
		p.decoder = firstLayerDecoder
		p.encoded = data
	}
	return p
}

func (p *packet) Layers() []Layer {
	return p.layers
}

func (p *packet) Layer(t LayerType) Layer {
	for _, l := range p.layers {
		if t == l.LayerType() {
			return l
		}
	}
	return nil
}

func (p *packet) LayerClass(lc LayerClass) Layer {
	for _, l := range p.layers {
		if lc.Contains(l.LayerType()) {
			return l
		}
	}
	return nil
}

func (p *packet) String() string {
	layers := []string{}
	for l := range p.Layers() {
		layers = append(layers, fmt.Sprintf("%#v", l))
	}
	return fmt.Sprintf("PACKET [%s]", strings.Join(layers, ", "))
}

type PacketDataSource interface {
	// ReadPacketData returns the next packet available from this data source.
	ReadPacketData() (data []byte, ci CaptureInfo, err error)
}

type PacketSource struct {
	source PacketDataSource
	Decoder
	DecodeOptions
}

func NewPacketSource(source PacketDataSource, decoder Decoder) *PacketSource {
	return &PacketSource{
		source:  source,
		Decoder: decoder,
	}
}

func (p *PacketSource) NextPacket() (Packet, error) {
	data, ci, err := p.source.ReadPacketData()
	if err != nil {
		return nil, err
	}
	packet := NewPacket(data, p.Decoder, p.DecodeOptions)
	*packet.CaptureInfo() = ci
	return packet, nil
}

func (p *PacketSource) PacketChannel() <-chan Packet {
	c := make(chan Packet, 100)
	go func() {
		for {
			packet, err := p.NextPacket()
			if err == io.EOF {
				close(c)
			} else if err == nil {
				c <- packet
			}
		}
	}()
	return c
}
