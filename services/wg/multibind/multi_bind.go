package multibind

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/conn"
)

var (
	_ conn.Bind = (*MultiBind)(nil)
)

type MultiBind struct {
	opened atomic.Bool

	transports []*Transport
	svcLogger  *logrus.Entry
}

func NewMultiBind(logger *logrus.Entry, trans ...*Transport) *MultiBind {
	if logger == nil {
		logger = logrus.NewEntry(logrus.New())
	}

	if len(trans) == 0 {
		panic("no transport provided")
	}

	return &MultiBind{
		transports: trans,
		svcLogger:  logger,
	}
}

// BatchSize implements conn.Bind.
func (m *MultiBind) BatchSize() int {
	bs := 1
	for _, t := range m.transports {
		bs = max(bs, t.bind.BatchSize())
	}
	return bs
}

// Close implements conn.Bind.
func (m *MultiBind) Close() error {
	var errs []error
	for _, t := range m.transports {
		if err := t.bind.Close(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.name, err))
		}
	}

	m.opened.Store(false)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Open implements conn.Bind.
func (m *MultiBind) Open(port uint16) (fns []conn.ReceiveFunc, actualPort uint16, err error) {
	if m.opened.Load() {
		return nil, 0, conn.ErrBindAlreadyOpen
	}
	m.opened.Store(true)

	multiRecvFunc := []conn.ReceiveFunc{}

	for _, t := range m.transports {
		fns, p, err := t.bind.Open(port)
		if err != nil {
			// TODO: 这里处理一下
			return nil, 0, err
		}

		if p != 0 {
			actualPort = p
		}

		for _, fn := range fns {
			multiRecvFunc = append(multiRecvFunc, m.recvWrapper(t, fn))
		}
	}

	return multiRecvFunc, actualPort, nil
}

// ParseEndpoint implements conn.Bind.
func (m *MultiBind) ParseEndpoint(s string) (conn.Endpoint, error) {
	log := m.svcLogger.WithField("op", "ParseEndpoint")
	var errs []error

	for _, t := range m.transports {
		ep, err := t.bind.ParseEndpoint(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.name, err))
			continue
		}

		log.Debugf("multibind parsed endpoint: %s, transport: %s", ep, t.name)
		return NewMultiEndpoint(t, ep), nil
	}

	if len(errs) == 0 {
		return nil, conn.ErrWrongEndpointType
	}
	return nil, errors.Join(conn.ErrWrongEndpointType,
		fmt.Errorf("failed to parse endpoint: %s, errors: %v", s, errs))
}

// Send implements conn.Bind.
func (m *MultiBind) Send(bufs [][]byte, ep conn.Endpoint) error {
	log := m.svcLogger.WithField("op", "Send")

	e, ok := ep.(*MultiEndpoint)
	if !ok {
		return fmt.Errorf("invalid endpoint type, not a MultiEndpoint")
	}

	log.Tracef("multibind sending packets to endpoint: %s, transport: %s", e.inner.DstToString(), e.trans.name)
	e.trans.bind.Send(bufs, e.inner)
	return nil
}

// SetMark implements conn.Bind.
func (m *MultiBind) SetMark(mark uint32) error {
	var errs []error
	for _, t := range m.transports {
		if err := t.bind.SetMark(mark); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.name, err))
		}
	}
	if len(errs) > 0 {
		return errors.Join(fmt.Errorf("failed to set mark: %d", mark), errors.Join(errs...))
	}

	return nil
}

// recvWrapper trans endpoint from inner transport to multiEndpoint
// for faster endpoint type classification
func (m *MultiBind) recvWrapper(trans *Transport, fns conn.ReceiveFunc) conn.ReceiveFunc {

	return func(packets [][]byte, sizes []int, eps []conn.Endpoint) (n int, err error) {
		log := m.svcLogger.WithField("op", "recvWrapper").WithField("transport", trans.name)

		tmpEps := make([]conn.Endpoint, len(eps))
		n, err = fns(packets, sizes, tmpEps)
		log.Tracef("multibind received packets: [%d] from transport: [%s], with endpoints length: [%d], sizes length: [%d], packets length: [%d]",
			n, trans.name, len(tmpEps), len(sizes), len(packets))

		for i := range eps {
			eps[i] = trans.loadOrNewEndpoint(tmpEps[i])
		}
		return n, err
	}
}
