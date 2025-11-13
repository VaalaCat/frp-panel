package multibind

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/conn"
)

var (
	_ conn.Bind = (*MultiBind)(nil)
)

type MultiBind struct {
	transports   []*Transport
	svcLogger    *logrus.Entry
	endpointPool sync.Pool
}

func NewMultiBind(logger *logrus.Entry, trans ...*Transport) *MultiBind {
	if logger == nil {
		logger = logrus.NewEntry(logrus.New())
	}

	if len(trans) == 0 {
		panic("no transport provided")
	}

	mb := &MultiBind{
		transports: trans,
		svcLogger:  logger,
	}

	mb.endpointPool.New = func() interface{} {
		eps := make([]conn.Endpoint, 0, 128)
		return &eps
	}

	return mb
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

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Open implements conn.Bind.
func (m *MultiBind) Open(port uint16) (fns []conn.ReceiveFunc, actualPort uint16, err error) {
	multiRecvFunc := []conn.ReceiveFunc{}

	for _, t := range m.transports {
		fns, p, err := t.bind.Open(port)
		if err != nil {
			// TODO: 这里处理一下
			return nil, 0, err
		}

		if p != 0 && t.name == "udp" {
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
	e, ok := ep.(*MultiEndpoint)
	if !ok {
		return fmt.Errorf("invalid endpoint type, not a MultiEndpoint")
	}

	return e.trans.bind.Send(bufs, e.inner)
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
		defer func() {
			if panicRecover := recover(); panicRecover != nil {
				err = fmt.Errorf("multibind recvWrapper panic: %v, debugStack: %s", panicRecover, debug.Stack())
				m.svcLogger.WithError(err).Error("multibind recvWrapper panic")
			} else if err != nil {
				m.svcLogger.WithError(err).Error("multibind recvWrapper error")
			}
		}()

		// 从对象池获取临时 endpoint 切片
		tmpEpsPtr := m.endpointPool.Get().(*[]conn.Endpoint)
		tmpEps := *tmpEpsPtr

		// 确保容量足够
		if cap(tmpEps) < len(eps) {
			tmpEps = make([]conn.Endpoint, len(eps))
		} else {
			tmpEps = tmpEps[:len(eps)]
		}

		n, err = fns(packets, sizes, tmpEps)

		// 批量转换 endpoint，只转换实际接收到的数量
		for i := 0; i < n; i++ {
			eps[i] = trans.loadOrNewEndpoint(tmpEps[i])
		}

		// 清空切片内容并归还到对象池
		for i := range tmpEps {
			tmpEps[i] = nil
		}
		tmpEps = tmpEps[:0]
		*tmpEpsPtr = tmpEps
		m.endpointPool.Put(tmpEpsPtr)

		return n, err
	}
}
