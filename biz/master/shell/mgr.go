package shell

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
)

type PTYMgr struct {
	*utils.SyncMap[string, pb.Master_PTYConnectServer]                                   // sessionID
	doneMap                                            *utils.SyncMap[string, chan bool] // sessionID
}

var (
	mgr *PTYMgr
)

func (m *PTYMgr) IsSessionDone(sessionID string) bool {
	ch, ok := m.doneMap.Load(sessionID)
	if !ok {
		return true
	}
	return <-ch
}

func (m *PTYMgr) SetSessionDone(sessionID string) {
	ch, ok := m.doneMap.Load(sessionID)
	if !ok {
		return
	}
	ch <- true
}

func (m *PTYMgr) Add(sessionID string, conn pb.Master_PTYConnectServer) {
	m.Store(sessionID, conn)
	m.doneMap.Store(sessionID, make(chan bool))
}

func NewPTYMgr() *PTYMgr {
	return &PTYMgr{
		SyncMap: &utils.SyncMap[string, pb.Master_PTYConnectServer]{},
		doneMap: &utils.SyncMap[string, chan bool]{},
	}
}

func Mgr() *PTYMgr {
	if mgr == nil {
		mgr = NewPTYMgr()
	}
	return mgr
}
