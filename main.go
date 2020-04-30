package main

var config *AppConf

func main() {
	config = NewAppConf()
	mh := NewMessageHandler()
	mh.Init()
}
