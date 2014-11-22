package lock

import (
	"sync"
)

var rwlock sync.RWMutex = sync.RWMutex{}
var unitLockDict map[string]*sync.RWMutex = map[string]*sync.RWMutex{}

type LockService struct{}

func (o LockService) GetUnitLock(key string) *sync.RWMutex {
	rwlock.Lock()
	defer rwlock.Unlock()
	
	if unitLockDict[key] == nil {
		unitLockDict[key] = &sync.RWMutex{}
	}
	return unitLockDict[key]
}
