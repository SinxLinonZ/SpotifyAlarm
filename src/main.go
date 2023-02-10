package main

import (
	"context"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/zmb3/spotify/v2"
	"log"
)

var (
	Client    *spotify.Client
	AlarmList []Alarm
)

func onReady() {
	/*****************
	 * 初始化 UI
	 *****************/
	systray.SetIcon(icon.Data)
	systray.SetTitle("Spotify Alarm")
	systray.SetTooltip("Spotify Alarm")

	mUsername := systray.AddMenuItem("-", "")
	mOpenUI := systray.AddMenuItem("Open UI", "Open UI")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	/*****************
	 * 注册菜单点击事件
	 *****************/
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mOpenUI.ClickedCh:
				openBrowser("http://localhost:13333")
			}
		}
	}()

	/*****************
	 * 初始化Spotify客户端
	 * 初始化闹钟列表，启动前端
	 */
	Client = getClient()
	AlarmList = InitAlarmDataFile()
	StartFrontend()
	startAlarmService()

	// Get username and display (also for token renew)
	username, err := Client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	mUsername.SetTitle("Login as " + username.DisplayName)
}

func onExit() {
	token, _ := Client.Token()
	saveTokenToFile(token)
	SaveAlarmData(&AlarmList)
}

func main() {
	systray.Run(onReady, onExit)
}
