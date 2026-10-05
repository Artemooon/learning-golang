package interfaces

import "net"

type IPv4 struct {
        // Version is the IP version number.
        // For IPv4 packets, this is always 4.
        // It occupies 4 bits in the real packet header.
        Version uint8

        // IHL means Internet Header Length.
        // It says how long the IPv4 header is in 32-bit (4-byte) words.
        // Minimum value: 5, meaning a 20-byte header.
        // Values greater than 5 mean IPv4 options are present.
        // It occupies 4 bits in the real packet header.
        IHL uint8

        // DSCP means Differentiated Services Code Point.
        // It is a 6-bit traffic-priority/quality-of-service value.
        DSCP uint8

        // ECN means Explicit Congestion Notification.
        // It is a 2-bit value used to signal network congestion
        // without necessarily dropping the packet.
        ECN uint8

        // TotalLength is the complete IPv4 packet size in bytes:
        // IPv4 header + payload.
        TotalLength uint16

        // Identification helps reassemble fragments that came from
        // the same original IP packet.
        Identification uint16

        // Flags control fragmentation.
        // The meaningful flags are:
        // - Don't Fragment (DF)
        // - More Fragments (MF)
        // This occupies 3 bits in the real packet header.
        Flags uint8

        // FragmentOffset says where this fragment belongs in the
        // original packet. It is measured in units of 8 bytes.
        // It occupies 13 bits in the real packet header.
        FragmentOffset uint16

        // TimeToLive limits how many routers may forward the packet.
        // Each router decreases it by one. If it reaches zero,
        // the router discards the packet.
        TimeToLive uint8

        // Protocol identifies the payload carried inside this IP packet.
        // Common values:
        //   1  = ICMP
        //   6  = TCP
        //   17 = UDP
        Protocol uint8

        // HeaderChecksum detects corruption in the IPv4 header only.
        // Routers recalculate it because TTL changes at every hop.
        // It does not cover the payload.
        HeaderChecksum uint16

        // SrcAddr is the packet sender's IPv4 address.
        // An IPv4 address is exactly 4 bytes on the wire.
        SrcAddr net.IP

        // DestAddr is the intended receiver's IPv4 address.
        // Routers use it to choose the next hop.
        DestAddr net.IP
}

type UDP struct {
	SrcPort  uint16
	DestPort uint16
	Length   uint16
	Checksum uint16
}

const (
	IPv4HeaderLength = 20
)