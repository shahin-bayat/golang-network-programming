package resolver

import (
	"fmt"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

var RootServers = []string{
	"198.41.0.4:53",     // A.ROOT-SERVERS.NET
	"170.247.170.2:53",  // B.ROOT-SERVERS.NET
	"192.33.4.12:53",    // C.ROOT-SERVERS.NET
	"199.7.91.13:53",    // D.ROOT-SERVERS.NET
	"192.203.230.10:53", // E.ROOT-SERVERS.NET
}

func Resolve(name dnsmessage.Name, qType dnsmessage.Type) ([]dnsmessage.Resource, error) {
	addr, err := net.ResolveUDPAddr("udp", RootServers[0])
	if err != nil {
		return nil, err
	}
	var response dnsmessage.Message

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		msg := dnsmessage.Message{
			Header: dnsmessage.Header{
				ID:               1,
				RecursionDesired: false,
			},
			Questions: []dnsmessage.Question{
				{
					Name:  name,
					Type:  qType,
					Class: dnsmessage.ClassINET,
				},
			},
		}
		packed, err := msg.Pack()
		if err != nil {
			return nil, err
		}
		if _, err := conn.Write(packed); err != nil {
			return nil, err
		}
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		if err := response.Unpack(buf[:n]); err != nil {
			return nil, err
		}
		if len(response.Answers) > 0 {
			break
		} else if len(response.Answers) == 0 && len(response.Authorities) == 0 {
			return nil, fmt.Errorf("could not resolve domain %s, no answers and no authorities found", name)
		} else {
			foundNextAddr := false

		Authorities:
			for _, authority := range response.Authorities {
				ns, ok := authority.Body.(*dnsmessage.NSResource)
				if !ok {
					continue Authorities
				}

			Additionals:
				for _, additional := range response.Additionals {
					if additional.Header.Name.String() != ns.NS.String() {
						continue Additionals
					}

					a, ok := additional.Body.(*dnsmessage.AResource)
					if !ok {
						continue Additionals
					}
					fmt.Println("found a resource", a)
					addr = &net.UDPAddr{IP: a.A[:], Port: 53}

					foundNextAddr = true
					break Additionals
				}

				if foundNextAddr {
					break Authorities
				}
			}
		}

	}
	fmt.Println("found answer")
	return response.Answers, nil
}
