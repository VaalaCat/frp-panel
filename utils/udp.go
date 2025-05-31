package utils

import (
	"context"
	"errors"
	"math"
	"net"
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/defs"
)

// ProbeEndpoint sends a small UDP packet to addr and waits for a reply.
// It returns the measured RTT or an error.
func ProbeEndpoint(ctx context.Context, addr EndpointGettable, timeout time.Duration) (time.Duration, error) {
	// Resolve UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", addr.GetEndpoint())
	if err != nil {
		return 0, err
	}

	// Dial UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Prepare a simple ping payload
	payload := []byte(defs.VaalaMagicBytes)

	// Set deadlines
	deadline := time.Now().Add(timeout)
	conn.SetDeadline(deadline)

	start := time.Now()
	if _, err := conn.Write(payload); err != nil {
		return 0, err
	}

	// Buffer for response
	buf := make([]byte, 64)
	if _, _, err := conn.ReadFrom(buf); err != nil {
		return 0, err
	}
	rtt := time.Since(start)

	return rtt, nil
}

type EndpointGettable interface {
	GetEndpoint() string
}

// SelectFastestEndpoint concurrently probes all candidates and returns the fastest.
func SelectFastestEndpoint(ctx context.Context, candidates []EndpointGettable, timeout time.Duration) (EndpointGettable, error) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		bestEP   EndpointGettable
		bestRTT  = time.Duration(math.MaxInt64)
		firstErr error
	)

	wg.Add(len(candidates))
	for _, addr := range candidates {
		go func(addr EndpointGettable) {
			defer wg.Done()

			rtt, err := ProbeEndpoint(ctx, addr, timeout)
			if err != nil {
				// record first error
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			if rtt < bestRTT {
				bestRTT = rtt
				bestEP = addr
			}
			mu.Unlock()
		}(addr)
	}

	wg.Wait()

	if bestEP == nil {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, errors.New("no endpoint reachable")
	}

	return bestEP, nil
}
