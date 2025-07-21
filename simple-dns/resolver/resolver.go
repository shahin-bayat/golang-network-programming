package resolver

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"simple-dns/cache"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

var RootServers = []string{
	"198.41.0.4:53",     // A.ROOT-SERVERS.NET
	"170.247.170.2:53",  // B.ROOT-SERVERS.NET
	"192.33.4.12:53",    // C.ROOT-SERVERS.NET
	"199.7.91.13:53",    // D.ROOT-SERVERS.NET
	"192.203.230.10:53", // E.ROOT-SERVERS.NET
}

type Resolver struct {
	cache *cache.DNSCache
}

func NewResolver() *Resolver {
	return &Resolver{
		cache: cache.New(),
	}
}

func (r *Resolver) Resolve(name dnsmessage.Name, qType dnsmessage.Type) ([]dnsmessage.Resource, error) {
	// Check cache first
	cacheKey := cache.CacheKey{Name: name, Type: qType}
	if entry, found := r.cache.Get(cacheKey); found {
		fmt.Printf("Cache hit for %s %s\n", name.String(), qType)
		return entry.Resources, nil
	}

	// If not found in cache, resolve the name using the root servers
	addr, err := net.ResolveUDPAddr("udp", RootServers[0])
	if err != nil {
		return nil, err
	}

	fmt.Printf("querying root server %s for %s\n", addr.String(), name.String())

	var response dnsmessage.Message

	// Iterate up to 10 times to follow referrals
MainLoop:
	for i := 0; i < 10; i++ {
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		msg := dnsmessage.Message{
			Header: dnsmessage.Header{
				ID:               uint16(rand.Intn(65535)),
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
			minTTL := minTTL(response.Answers)
			cacheEntry := cache.CacheEntry{Resources: response.Answers, ExpiresAt: minTTL}
			if cname, ok := response.Answers[0].Body.(*dnsmessage.CNAMEResource); ok {
				fmt.Printf("following CNAME from %s to %s\n", name.String(), cname.CNAME.String())
				r.cache.Set(cacheKey, cacheEntry)
				name = cname.CNAME
				continue MainLoop
			} else {
				finalCacheKey := cache.CacheKey{Name: name, Type: qType}
				r.cache.Set(finalCacheKey, cacheEntry)
				return response.Answers, nil
			}
		}

		// If there are no authorities, we can't go any further
		if len(response.Authorities) == 0 {
			return nil, fmt.Errorf("could not resolve domain %s: no answers and no authorities found", name)
		}

		// Otherwise, we need to follow the referral
		foundNextAddr := false

	Authorities:
		for _, authority := range response.Authorities {
			ns, ok := authority.Body.(*dnsmessage.NSResource)
			if !ok {
				continue
			}

		Additionals: // Look for a glue record (the IP of the nameserver)
			for _, additional := range response.Additionals {
				if additional.Header.Name.String() != ns.NS.String() {
					continue Additionals
				}

				a, ok := additional.Body.(*dnsmessage.AResource)
				if !ok {
					continue Additionals
				}

				// Found the IP of the next server to query
				addr = &net.UDPAddr{IP: a.A[:], Port: 53}
				foundNextAddr = true
				break Additionals
			}

			if foundNextAddr {
				break Authorities
			} else {
				// If there's no glue record, we have to resolve the nameserver's IP ourselves
				recursiveResult, err := r.Resolve(ns.NS, dnsmessage.TypeA)
				if err != nil {
					continue Authorities // Try the next authority if this one fails
				}
				// Find the first A record from the recursive call
				for _, rr := range recursiveResult {
					ra, ok := rr.Body.(*dnsmessage.AResource)
					if !ok {
						continue
					}
					// Found the IP, update addr and break to the main loop
					addr = &net.UDPAddr{IP: ra.A[:], Port: 53}
					foundNextAddr = true
					break Authorities
				}
			}
		}

		// If we went through all authorities and couldn't find the next address, we're stuck.
		if !foundNextAddr {
			return nil, fmt.Errorf("could not find next address to query for %s", name)
		}

		fmt.Printf("redirected to query %s for %s\n", addr.String(), name.String())
	}

	return nil, fmt.Errorf("resolver timeout: exceeded 10 iterations for %s", name)
}

func minTTL(answers []dnsmessage.Resource) time.Time {
	minTTL := uint32(math.MaxUint32)
	for _, answer := range answers {
		if answer.Header.TTL < minTTL {
			minTTL = answer.Header.TTL
		}
	}
	return time.Now().Add(time.Second * time.Duration(minTTL))
}
