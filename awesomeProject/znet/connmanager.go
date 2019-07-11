package znet

import (
	"awesomeProject/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	Connections map[uint32]ziface.IConnection
	ConnLock    sync.RWMutex
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	cm.Connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager sucessfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.ConnLock.Lock() //写锁
	defer cm.ConnLock.Unlock()
	if _, ok := cm.Connections[conn.GetConnID()]; ok {
		delete(cm.Connections, conn.GetConnID())
	}
	fmt.Println("connection removed from ConnManager sucessfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.ConnLock.RLock() //读锁
	defer cm.ConnLock.RUnlock()
	if conn, ok := cm.Connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection is not found")
}

func (cm *ConnManager) Len() int {
	//cm.ConnLock.RLock()
	//defer cm.ConnLock.RUnlock()
	return len(cm.Connections)
}

func (cm *ConnManager) ClearConn() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	for connID, conn := range cm.Connections {
		conn.Stop()
		fmt.Println("connection ", connID, " is Stop!!!")
		delete(cm.Connections, connID)
	}

	fmt.Println("Clear all the Connections, conn num = ", cm.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}
