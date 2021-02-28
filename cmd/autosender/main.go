package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"discord_autosender/cmd/autosender/config"
	"discord_autosender/cmd/autosender/utils"
	"discord_autosender/pkg/curlparser"
	"github.com/getlantern/systray"
)

const MaxMinutes = 60

func main() {
	systray.Run(onReady, onExit)
}

func SendMessage() []byte {
	var c config.Conf
	var clientConfig curlparser.CurlParser

	clientConfig.GetConf()
	c.GetConf()
	url := fmt.Sprintf("https://discord.com/api/v8/channels/%s/messages?limit=50", c.Message.ChannelID)

	body := struct {
		Content string `json:"content"`
		Tts     bool   `json:"tts"`
	}{
		Content: c.Message.Text,
		Tts:     false,
	}

	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Fatalln(err)
	}

	for header, value := range clientConfig.Headers {
		req.Header.Add(header, value)
	}

	req.Header.Set("referer", fmt.Sprintf("https://discord.com/channels/%s/%s", c.Message.ServerID, c.Message.ChannelID))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://discord.com")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return res
}

func onReady() {
	var c config.Conf

	c.GetConf()

	var duration time.Duration

	if c.IsDebug {
		duration = time.Minute
	} else {
		duration = time.Hour
	}

	ticker := time.NewTicker(duration)

	systray.SetIcon(getIcon("./assets/icon.ico"))
	systray.SetTooltip("Discord AutoSender <3")

	quit := systray.AddMenuItem("Quit", "Stop autosend message")

	SendWithRandom(c)

	go func() bool {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Ticker triggered")
				SendWithRandom(c)
			case <-quit.ClickedCh:
				ticker.Stop()
				systray.Quit()
				return true
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here
	os.Exit(1)
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}

	return b
}

func isNeedRequireSend(c config.Conf) bool {
	currTime := time.Now()

	for _, hour := range c.Planned {
		if currTime.Hour() == hour { //nolint:wsl
			return true
		}
	}

	return false
}

func SendWithRandom(c config.Conf) {
	if isNeedRequireSend(c) {
		rand.Seed(time.Now().UnixNano())

		currentTime := time.Now()

		randomMinute := utils.RandInt(c.Settings.MinRandMinute, c.Settings.MaxRandMinute)

		if (MaxMinutes - currentTime.Minute()) < randomMinute {
			randomMinute = MaxMinutes - currentTime.Minute()
		}

		if c.IsDebug {
			randomMinute = 0
		}

		fmt.Printf("Started after %v minutes \n", randomMinute)

		go func() {
			defer fmt.Println("Exit of routine \n")

			<-time.After(time.Duration(randomMinute) * time.Minute)
			fmt.Printf("SendMessage at %s \n", time.Now())
			SendMessage()
		}()
	}
}
