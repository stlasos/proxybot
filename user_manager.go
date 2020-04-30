package main

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserManager struct {
	dataToTest     map[int]*ProxyData
	dataToTestLock sync.RWMutex
	proxyChecker   map[int]*ProxyChecker
}

type ProxyData struct {
	Ips   []string
	Ports []string
}

func NewUserManager() *UserManager {
	return &UserManager{
		dataToTest:   make(map[int]*ProxyData),
		proxyChecker: make(map[int]*ProxyChecker),
	}
}

func (um *UserManager) IsInProgress(id int) bool {
	if _, ok := um.proxyChecker[id]; ok {
		return um.proxyChecker[id].GetIsInProgress()
	}
	return false
}

func (um *UserManager) createDataForUser(id int) {
	if _, exists := um.dataToTest[id]; !exists {
		um.dataToTest[id] = &ProxyData{}
	}
}

func (um *UserManager) SetUserIps(id int, ips []string) {
	um.dataToTestLock.Lock()
	defer um.dataToTestLock.Unlock()
	um.createDataForUser(id)
	um.dataToTest[id].Ips = ips
}

func (um *UserManager) SetUserPorts(id int, ports []string) {
	um.dataToTestLock.Lock()
	defer um.dataToTestLock.Unlock()
	um.createDataForUser(id)
	for i, _ := range ports {
		ports[i] = strings.TrimPrefix(ports[i], ":")
	}
	um.dataToTest[id].Ports = ports
}

func (um *UserManager) GetUserIps(id int) []string {
	um.dataToTestLock.RLock()
	defer um.dataToTestLock.RUnlock()
	return um.dataToTest[id].Ips
}

func (um *UserManager) GetUserPorts(id int) []string {
	um.dataToTestLock.RLock()
	defer um.dataToTestLock.RUnlock()
	return um.dataToTest[id].Ports
}

func (um *UserManager) getDataToTest(id int) []string {
	um.dataToTestLock.RLock()
	defer um.dataToTestLock.RUnlock()
	res := make([]string, 0)
	for _, ip := range um.dataToTest[id].Ips {
		for _, port := range um.dataToTest[id].Ports {
			res = append(res, ip+":"+port)
		}
	}
	return res
}

func (um *UserManager) StartChecker(id int, msg chan string) {
	ks := make(chan bool)
	urls := um.getDataToTest(id)
	um.proxyChecker[id] = NewProxyChecker()
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			select {
			case <-ks:
				return
			default:
				msg <- "Progress: " + strconv.Itoa(um.proxyChecker[id].GetDoneCount()) + "/" + strconv.Itoa(len(urls))
			}
		}
	}()
	result := um.proxyChecker[id].Init(urls)
	ks <- true
	var b bytes.Buffer

	b.WriteString("All urls were checked.\n")
	if len(result) > 0 {
		b.WriteString("Successful:\n")
		for _, res := range result {
			b.WriteString(res)
			b.WriteString("\n")
		}
	} else {
		b.WriteString("Proxy urls not found.\n")
	}
	msg <- b.String()
}
