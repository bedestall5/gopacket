// Copyright 2012 Google, Inc. All rights reserved.
// Copyright 2009-2011 Andreas Krennmair. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package pcap

/*
#cgo linux LDFLAGS: -lpcap
#cgo freebsd LDFLAGS: -lpcap
#cgo darwin LDFLAGS: -lpcap
#cgo windows CFLAGS: -I C:/WpdPack/Include
#cgo windows LDFLAGS: -L C:/WpdPack/Lib -lwpcap
#include <stdlib.h>
#include <pcap.h>
*/
import "C"

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"errors"
	"io"
	"net"
	"runtime"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const errorBufferSize = 256

// Handle provides a connection to a pcap handle, allowing users to read packets
// off the wire (Next), inject packets onto the wire (Inject), and
// perform a number of other functions to affect and understand packet output.
type Handle struct {
	// cptr is the handle for the actual pcap C object.
	cptr *C.pcap_t
}

// Stats contains statistics on how many packets were handled by a pcap handle,
// and what was done with those packets.
type Stats struct {
	PacketsReceived  int
	PacketsDropped   int
	PacketsIfDropped int
}

// Interface describes a single network interface on a machine.
type Interface struct {
	Name        string
	Description string
	Addresses   []InterfaceAddress
	// TODO: add more elements
}

// InterfaceAddress describes an address associated with an Interface.
// Currently, it's IPv4/6 specific.
type InterfaceAddress struct {
	IP      net.IP
	Netmask net.IPMask
	// TODO: add broadcast + PtP dst ?
}

// BlockForever, when passed into OpenLive, causes it to block forever waiting for packets.
const BlockForever = time.Duration(0)

// OpenLive opens a device and returns a *Handle.
// It takes as arguments the name of the device ("eth0"), the maximum size to
// read for each packet (snaplen), whether to put the interface in promiscuous
// mode, and a timeout.
func OpenLive(device string, snaplen int32, promisc bool, timeout time.Duration) (handle *Handle, _ error) {
	var buf *C.char
	buf = (*C.char)(C.calloc(errorBufferSize, 1))
	defer C.free(unsafe.Pointer(buf))
	var pro int32
	if promisc {
		pro = 1
	}

	dev := C.CString(device)
	defer C.free(unsafe.Pointer(dev))

	cptr := C.pcap_open_live(dev, C.int(snaplen), C.int(pro), C.int(timeout/time.Millisecond), buf)
	if cptr == nil {
		return nil, errors.New(C.GoString(buf))
	}
	return newHandle(cptr), nil
}

func newHandle(cptr *C.pcap_t) (handle *Handle) {
	handle = &Handle{cptr: cptr}
	return
}

// OpenOffline opens a file and returns its contents as a *Handle.
func OpenOffline(file string) (handle *Handle, err error) {
	var buf *C.char
	buf = (*C.char)(C.calloc(errorBufferSize, 1))
	defer C.free(unsafe.Pointer(buf))
	cf := C.CString(file)
	defer C.free(unsafe.Pointer(cf))

	cptr := C.pcap_open_offline(cf, buf)
	if cptr == nil {
		return nil, errors.New(C.GoString(buf))
	}
	return newHandle(cptr), nil
}

// NextError is the return code from a call to Next.
type NextError int32

// NextError implements the error interface.
func (n NextError) Error() string {
	switch n {
	case NextErrorOk:
		return "OK"
	case NextErrorTimeoutExpired:
		return "Timeout Expired"
	case NextErrorReadError:
		return "Read Error"
	case NextErrorNoMorePackets:
		return "No More Packets In File"
	}
	return strconv.Itoa(int(n))
}

const (
	NextErrorOk             NextError = 1
	NextErrorTimeoutExpired NextError = 0
	NextErrorReadError      NextError = -1
	// NextErrorNoMorePackets is returned when reading from a file (OpenOffline) and
	// EOF is reached.  When this happens, Next() returns io.EOF instead of this.
	NextErrorNoMorePackets NextError = -2
)

// NextError returns the next packet read from the pcap handle, along with an error
// code associated with that packet.  If the packet is read successfully, the
// returned error is nil.
func (p *Handle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	var pkthdr *C.struct_pcap_pkthdr

	var buf_ptr *C.u_char
	var buf unsafe.Pointer

	runtime.LockOSThread()
	result := NextError(C.pcap_next_ex(p.cptr, &pkthdr, &buf_ptr))

	buf = unsafe.Pointer(buf_ptr)

	if nil == buf {
		if result == NextErrorNoMorePackets {
			err = io.EOF
		} else {
			err = result
		}
		runtime.UnlockOSThread()
		return
	}
	data = C.GoBytes(buf, C.int(pkthdr.caplen))
	ci.Timestamp = time.Unix(int64(pkthdr.ts.tv_sec), int64(pkthdr.ts.tv_usec))
	ci.CaptureLength = int(pkthdr.caplen)
	ci.Length = int(pkthdr.len)
	runtime.UnlockOSThread()
	return
}

// Close closes the underlying pcap handle.
func (p *Handle) Close() {
	C.pcap_close(p.cptr)
}

// Error returns the current error associated with a pcap handle (pcap_geterr).
func (p *Handle) Error() error {
	return errors.New(C.GoString(C.pcap_geterr(p.cptr)))
}

// Stats returns statistics on the underlying pcap handle.
func (p *Handle) Stats() (stat *Stats, err error) {
	var cstats _Ctype_struct_pcap_stat
	if -1 == C.pcap_stats(p.cptr, &cstats) {
		return nil, p.Error()
	}
	return &Stats{
		PacketsReceived:  int(cstats.ps_recv),
		PacketsDropped:   int(cstats.ps_drop),
		PacketsIfDropped: int(cstats.ps_ifdrop),
	}, nil
}

// SetBPFFilter compiles and sets a BPF filter for the pcap handle.
func (p *Handle) SetBPFFilter(expr string) (err error) {
	var bpf _Ctype_struct_bpf_program
	cexpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cexpr))

	if -1 == C.pcap_compile(p.cptr, &bpf, cexpr, 1, 0) {
		return p.Error()
	}

	if -1 == C.pcap_setfilter(p.cptr, &bpf) {
		C.pcap_freecode(&bpf)
		return p.Error()
	}
	C.pcap_freecode(&bpf)
	return nil
}

// Version returns pcap_lib_version.
func Version() string {
	return C.GoString(C.pcap_lib_version())
}

// LinkType returns pcap_datalink, as a layers.LinkType.
func (p *Handle) LinkType() layers.LinkType {
	return layers.LinkType(C.pcap_datalink(p.cptr))
}

// SetLinkType calls pcap_set_datalink on the pcap handle.
func (p *Handle) SetLinkType(dlt layers.LinkType) error {
	if -1 == C.pcap_set_datalink(p.cptr, C.int(dlt)) {
		return p.Error()
	}
	return nil
}

// FindAllDevs attempts to enumerate all interfaces on the current machine.
func FindAllDevs() (ifs []Interface, err error) {
	var buf *C.char
	buf = (*C.char)(C.calloc(errorBufferSize, 1))
	defer C.free(unsafe.Pointer(buf))
	var alldevsp *C.pcap_if_t

	if -1 == C.pcap_findalldevs((**C.pcap_if_t)(&alldevsp), buf) {
		return nil, errors.New(C.GoString(buf))
	}
	defer C.pcap_freealldevs((*C.pcap_if_t)(alldevsp))
	dev := alldevsp
	var i uint32
	for i = 0; dev != nil; dev = (*C.pcap_if_t)(dev.next) {
		i++
	}
	ifs = make([]Interface, i)
	dev = alldevsp
	for j := uint32(0); dev != nil; dev = (*C.pcap_if_t)(dev.next) {
		var iface Interface
		iface.Name = C.GoString(dev.name)
		iface.Description = C.GoString(dev.description)
		iface.Addresses = findalladdresses(dev.addresses)
		// TODO: add more elements
		ifs[j] = iface
		j++
	}
	return
}

func findalladdresses(addresses *_Ctype_struct_pcap_addr) (retval []InterfaceAddress) {
	// TODO - make it support more than IPv4 and IPv6?
	retval = make([]InterfaceAddress, 0, 1)
	for curaddr := addresses; curaddr != nil; curaddr = (*_Ctype_struct_pcap_addr)(curaddr.next) {
		var a InterfaceAddress
		var err error
		if a.IP, err = sockaddr_to_IP((*syscall.RawSockaddr)(unsafe.Pointer(curaddr.addr))); err != nil {
			continue
		}
		if a.Netmask, err = sockaddr_to_IP((*syscall.RawSockaddr)(unsafe.Pointer(curaddr.addr))); err != nil {
			continue
		}
		retval = append(retval, a)
	}
	return
}

func sockaddr_to_IP(rsa *syscall.RawSockaddr) (IP []byte, err error) {
	switch rsa.Family {
	case syscall.AF_INET:
		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
		IP = make([]byte, 4)
		for i := 0; i < len(IP); i++ {
			IP[i] = pp.Addr[i]
		}
		return
	case syscall.AF_INET6:
		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		IP = make([]byte, 16)
		for i := 0; i < len(IP); i++ {
			IP[i] = pp.Addr[i]
		}
		return
	}
	err = errors.New("Unsupported address type")
	return
}

// WritePacketData calls pcap_sendpacket, injecting the given data into the pcap handle.
func (p *Handle) WritePacketData(data []byte) (err error) {
	buf := C.CString(string(data))
	defer C.free(unsafe.Pointer(buf))

	if -1 == C.pcap_sendpacket(p.cptr, (*C.u_char)(unsafe.Pointer(buf)), (C.int)(len(data))) {
		err = p.Error()
	}
	return
}
