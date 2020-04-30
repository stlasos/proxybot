package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type MessageHandler struct {
	bot *tb.Bot
	um  *UserManager
}

func NewMessageHandler() *MessageHandler {
	if config.ApiUrl == "" {
		config.ApiUrl = "https://api.telegram.org"
	}

	bot, err := tb.NewBot(tb.Settings{
		URL:    config.ApiUrl,
		Token:  config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	return &MessageHandler{bot: bot, um: NewUserManager()}
}

func (mh *MessageHandler) Init() {
	mh.bot.Handle("/start", mh.greetingsHandler)

	mh.bot.Handle("/ips", mh.ipsHandler)

	mh.bot.Handle("/ports", mh.portsHandler)

	mh.bot.Start()
}

func (mh *MessageHandler) greetingsHandler(m *tb.Message) {

	mh.bot.Send(m.Sender, `I can check urls on being able to act like a proxy.
Following commands are available:

/ips - You can load ip list separated by a newline character (below the command in the same message)
/ports - You can load ports the same way (below the command in the same message)
								
After both of these commands are done url check should start
You will be notified about progress and test results.`)
}

func (mh *MessageHandler) ipsHandler(m *tb.Message) {
	if !mh.um.IsInProgress(m.Sender.ID) {
		ips := strings.Fields(m.Text)[1:]
		if mh.isIpsValid(ips) {
			mh.um.SetUserIps(m.Sender.ID, ips)
			if len(mh.um.GetUserPorts(m.Sender.ID)) > 0 {
				mh.bot.Send(m.Sender, "Ips are accepted. Proxy check will begin now.")
				ch := make(chan string)
				mh.initMsgChannel(m.Sender, ch)
				mh.um.StartChecker(m.Sender.ID, ch)
			} else {
				mh.bot.Send(m.Sender, "Please fill ports for check to start.")
			}
		} else {
			mh.bot.Send(m.Sender, "Ip addresses seems to be invalid. Please check and fill them again.")
		}
	} else {
		mh.bot.Send(m.Sender, "Please wait. Proxy check is still in progress.")
	}
}

func (mh *MessageHandler) portsHandler(m *tb.Message) {
	if !mh.um.IsInProgress(m.Sender.ID) {
		ports := strings.Fields(m.Text)[1:]
		if mh.isPortsValid(ports) {
			mh.um.SetUserPorts(m.Sender.ID, ports)
			if len(mh.um.GetUserIps(m.Sender.ID)) > 0 {
				mh.bot.Send(m.Sender, "Ports are accepted. Proxy check will begin now.")
				ch := make(chan string)
				mh.initMsgChannel(m.Sender, ch)
				mh.um.StartChecker(m.Sender.ID, ch)
			} else {
				mh.bot.Send(m.Sender, "Please fill ips for check to start.")
			}
		} else {
			mh.bot.Send(m.Sender, "Ports seems to be invalid. Please check and fill them again.")
		}
	} else {
		mh.bot.Send(m.Sender, "Please wait. Proxy check is still in progress.")
	}
}

func (mh *MessageHandler) initMsgChannel(to *tb.User, msg chan string) {
	go func() {
		for {
			select {
			case m := <-msg:
				mh.bot.Send(to, m)
			}
		}
	}()
}

func (mh *MessageHandler) isIpsValid(ips []string) bool {
	for _, val := range ips {
		if net.ParseIP(val) == nil {
			return false
		}
	}
	return true
}

func (mh *MessageHandler) isPortsValid(ports []string) bool {
	for _, val := range ports {
		if res, _ := regexp.Match("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$", []byte(val)); !res {
			return false
		}
	}
	return true
}
