//go:build !windows
// +build !windows

package wg

import "time"

const (
	ReportInterval = time.Second * 60
)

func (w *wireGuard) reportStatusTask() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			w.pingPeers()
			time.Sleep(ReportInterval)
		}
	}
}
