package main

import (
	"net/http"
	"net/url"
	"sync"
)

const workers = 500

type ProxyChecker struct {
	results    []string
	doneCount  int
	doneLock   sync.RWMutex
	inProgress bool
}

func NewProxyChecker() *ProxyChecker {
	return &ProxyChecker{
		results:    make([]string, 0),
		doneCount:  0,
		inProgress: false,
	}
}

func (p *ProxyChecker) Init(requestUrls []string) []string {
	p.inProgress = true
	defer func() { p.inProgress = false }()
	workersAmount := workers
	if len(requestUrls) < workers {
		workersAmount = len(requestUrls)
	}
	req := make(chan string)
	resultsCh := make(chan string)
	done := make(chan bool)
	ks := make(chan bool)
	p.results = make([]string, 0)
	p.doneCount = 0
	for i := 0; i < workersAmount; i++ {
		go p.initCheckerWorker(req, done, resultsCh, ks)
	}
	for j := 0; j < len(requestUrls); j++ {
		go func(url string) {
			req <- url
		}(requestUrls[j])
	}
	go func() {
		for {
			select {
			case res := <-resultsCh:
				p.results = append(p.results, res)
			}
		}
	}()
	for c := 0; c < len(requestUrls); c++ {
		<-done
		p.doneLock.Lock()
		p.doneCount++
		p.doneLock.Unlock()
	}
	ks <- true
	return p.results
}

func (p *ProxyChecker) initCheckerWorker(req chan string, done chan bool, res chan string, ks chan bool) {
	for {
		select {
		case urlToCheck := <-req:
			{
				if p.makeRequest(urlToCheck) {
					res <- urlToCheck
				}
				done <- true
			}
		case <-ks:
			return
		}
	}
}

func (p *ProxyChecker) makeRequest(pUrl string) bool {
	proxyUrl, err := url.Parse("http://" + pUrl)

	if err == nil {
		client := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
		resp, reqErr := client.Get(config.ConnectUrl)
		return reqErr == nil && resp.StatusCode == 200
	}
	return false
}

func (p *ProxyChecker) GetDoneCount() int {
	p.doneLock.RLock()
	defer p.doneLock.RUnlock()
	return p.doneCount
}

func (p *ProxyChecker) GetIsInProgress() bool {
	return p.inProgress
}
