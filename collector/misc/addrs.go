package misc

import (
	"sync"
)

// Addrs 地址信息，用来做addr和appname、agentid映射
type Addrs struct {
	sync.RWMutex
	Addrs map[string]*Addr
}

// AddrStore ...
var AddrStore *Addrs

// Addr ...
type Addr struct {
	AppName string // 应用名
	// Agents  map[string]struct{} // agent信息
}

// NewAddrs ...
func NewAddrs() *Addrs {
	return &Addrs{
		Addrs: make(map[string]*Addr),
	}
}

// Add ...
func (a *Addrs) Add(appName, ip string) {
	a.RLock()
	_, ok := a.Addrs[ip]
	a.RUnlock()
	if !ok {
		a.Lock()
		a.Addrs[ip] = &Addr{
			AppName: appName,
		}
		a.Unlock()
	}

}

// Remove ...
func (a *Addrs) Remove(ip string) {
	a.Lock()
	delete(a.Addrs, ip)
	a.Lock()
}

func (a *Addrs) Get(ip string) (string, bool) {
	a.RLock()
	addrInfo, ok := a.Addrs[ip]
	a.RUnlock()
	if !ok {
		return "", false
	}
	return addrInfo.AppName, true
}

func initAddrStore() {
	AddrStore = NewAddrs()
}
