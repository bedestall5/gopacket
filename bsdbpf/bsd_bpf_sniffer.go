// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// +build darwin dragonfly freebsd netbsd openbsd

package bsdbpf

import (
	"github.com/google/gopacket"
	"golang.org/x/sys/unix"

	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))

func bpfWordAlign(x int) int {
	return (((x) + (wordSize - 1)) &^ (wordSize - 1))
}

type Options struct {
	bpfDeviceName    string
	readBufLen       int
	timeout          *syscall.Timeval
	promisc          bool
	immediate        bool
	preserveLinkAddr bool
}

var defaultOptions = Options{
	bpfDeviceName:    "",
	readBufLen:       0,
	timeout:          nil,
	promisc:          true,
	immediate:        false,
	preserveLinkAddr: true,
}

type BPFSniffer struct {
	options           *Options
	sniffDeviceName   string
	fd                int
	readBuffer        []byte
	lastReadLen       int
	readBytesConsumed int
}

func NewBPFSniffer(iface string, options *Options) *BPFSniffer {
	sniffer := BPFSniffer{
		sniffDeviceName: iface,
	}
	if options == nil {
		sniffer.options = &defaultOptions
	} else {
		sniffer.options = options
	}
	return &sniffer
}

func (b *BPFSniffer) Close() error {
	return syscall.Close(b.fd)
}

func (b *BPFSniffer) pickBpfDevice() {
	var err error
	for i := 0; i < 99; i++ {
		b.options.bpfDeviceName = fmt.Sprintf("/dev/bpf%d", i)
		b.fd, err = syscall.Open(b.options.bpfDeviceName, syscall.O_RDWR, 0)
		if err == nil {
			break
		}
	}
}

// Init is used to initialize a BPF device for promiscuous sniffing.
// It also starts a goroutine to continuously read frames.
func (b *BPFSniffer) Init() error {
	var err error
	enable := 1

	if b.options.bpfDeviceName == "" {
		b.pickBpfDevice()
	}

	// setup our read buffer
	if b.options.readBufLen == 0 {
		b.options.readBufLen, err = syscall.BpfBuflen(b.fd)
		if err != nil {
			panic(err)
		}
	} else {
		b.options.readBufLen, err = syscall.SetBpfBuflen(b.fd, b.options.readBufLen)
		if err != nil {
			panic(err)
		}
	}
	b.readBuffer = make([]byte, b.options.readBufLen)

	err = syscall.SetBpfInterface(b.fd, b.sniffDeviceName)
	if err != nil {
		return err
	}

	if b.options.immediate {
		// turn immediate mode on. This makes the snffer non-blocking.
		err = syscall.SetBpfImmediate(b.fd, enable)
		if err != nil {
			return err
		}
	}

	// the above call to syscall.SetBpfImmediate needs to be made
	// before setting a timer otherwise the reads will block for the
	// entire timer duration even if there are packets to return.
	if b.options.timeout != nil {
		err = syscall.SetBpfTimeout(b.fd, b.options.timeout)
		if err != nil {
			panic(err)
		}
	}

	if b.options.preserveLinkAddr {
		// preserves the link level source address...
		// higher level protocol analyzers will not need this
		err = syscall.SetBpfHeadercmpl(b.fd, enable)
		if err != nil {
			return err
		}
	}

	if b.options.promisc {
		// forces the interface into promiscuous mode
		err = syscall.SetBpfPromisc(b.fd, enable)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BPFSniffer) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	var err error
	if b.readBytesConsumed >= b.lastReadLen {
		b.readBytesConsumed = 0
		b.readBuffer = make([]byte, b.options.readBufLen)
		b.lastReadLen, err = syscall.Read(b.fd, b.readBuffer)
		if err != nil {
			b.lastReadLen = 0
			return nil, gopacket.CaptureInfo{}, err
		}
	}
	hdr := (*unix.BpfHdr)(unsafe.Pointer(&b.readBuffer[b.readBytesConsumed]))
	frameStart := b.readBytesConsumed + int(hdr.Hdrlen)
	b.readBytesConsumed += bpfWordAlign(int(hdr.Hdrlen) + int(hdr.Caplen))
	rawFrame := b.readBuffer[frameStart : frameStart+int(hdr.Caplen)]
	captureInfo := gopacket.CaptureInfo{
		Timestamp:     time.Unix(int64(hdr.Tstamp.Sec), int64(hdr.Tstamp.Usec)*1000),
		CaptureLength: len(rawFrame),
		Length:        len(rawFrame),
	}
	return rawFrame, captureInfo, nil
}
