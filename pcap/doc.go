// Copyright 2012 Google, Inc. All rights reserved.

/*
Package pcap allows users of gopacket to read packets off the wire or from
pcap files.

Reading PCAP Files

The following code can be used to read in data from a pcap file.

 if handle, err := pcap.OpenOffline("/path/to/my/file"); err != nil {
   panic(err)
 } else {
   packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	 for packet := range packetSource.Packets() {
     handlePacket(packet)  // Do something with a packet here.
   }
 }

Reading Live Packets

The following code can be used to read in data from a live device, in this case
"eth0".

 if handle, err := pcap.OpenLive("eth0", 1600, true, 0); err != nil {
   panic(err)
 } else if err := handle.SetBPFFilter("tcp and port 80"); err != nil {  // optional
   panic(err)
 } else {
   packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	 for packet := range packetSource.Packets() {
     handlePacket(packet)  // Do something with a packet here.
   }
 }
*/
package pcap
