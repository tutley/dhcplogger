package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"gopkg.in/mgo.v2"

	"flag"
	"log"
	"os"
	"time"
)

var (
	device                    = "en0"
	snapshotLen int32         = 4000
	promiscuous               = true
	timeout     time.Duration = 300 * time.Second
	handle      *pcap.Handle
)

func handlepacket(p gopacket.Packet, session *mgo.Session) {
	// Grab the DHCP info
	dhcp4Layer := p.Layer(layers.LayerTypeDHCPv4)
	if dhcp4Layer != nil {
		var a Assignment

		d, _ := dhcp4Layer.(*layers.DHCPv4)
		// dhcp operation 2 is a Reply
		if d.Operation == 2 {
			a.Mac = d.ClientHWAddr.String()
			ipv4 := d.ClientIP.String()
			if d.ClientIP.String() == "0.0.0.0" {
				ipv4 = d.YourClientIP.String()
			}
			a.Time = time.Now()
			addAssignment(ipv4, a, session)
		}
	}
}

func main() {
	// Read in command line arguments
	filterPtr := flag.String(
		"filter",
		"",
		"Additional packet filters such as host IP address")
	flag.StringVar(&device, "i", "eth0", "The device to inspect - default eht0")
	flag.Parse()
	var extraFilter = *filterPtr

	// Init DB connection
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Println("Error Connecting to Database: ", err)
		os.Exit(2)
	}
	defer db.Close()

	db.SetMode(mgo.Monotonic, true)

	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter = "(port 67 or port 68)"
	if extraFilter != "" {
		filter = filter + " and " + extraFilter
	}
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Capturing traffic on interface " + device + " with the filter: " + filter)

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlepacket(packet, db)
	}
}
