package main

import (
	"context"
	"encoding/json"
	"github.com/zmb3/spotify/v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func StartFrontend() {
	m := http.NewServeMux()
	s := http.Server{Addr: ":13333", Handler: m}

	/*****************
	 * 前端页面
	 */
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/index.html")
	})

	/*****************
	 * 后端API
	 */

	// 获取当前登录的用户名
	m.HandleFunc("/api/username", func(w http.ResponseWriter, r *http.Request) {
		username, _ := Client.CurrentUser(context.Background())
		_, _ = w.Write([]byte(username.DisplayName))
	})

	// 转换URL为URI
	m.HandleFunc("/api/url/toURI", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")

		linkRegex := regexp.MustCompile(`spotify.com/(track|album|artist|playlist)/([a-zA-Z0-9]+)`)
		if !linkRegex.MatchString(url) {
			_, _ = w.Write([]byte("invalid url"))
			return
		}

		linkMatch := linkRegex.FindStringSubmatch(url)
		URI := "spotify:" + linkMatch[1] + ":" + linkMatch[2]
		_, _ = w.Write([]byte(URI))
	})

	// 获取URL详情
	m.HandleFunc("/api/uri/query", func(w http.ResponseWriter, r *http.Request) {
		URIParts := strings.Split(r.URL.Query().Get("uri"), ":")
		URIType := URIParts[1]
		URIId := URIParts[2]

		switch URIType {
		case "track":
			track, _ := Client.GetTrack(context.Background(), spotify.ID(URIId))
			jsonData, err := json.Marshal(track)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = w.Write(jsonData)
		case "album":
			album, _ := Client.GetAlbum(context.Background(), spotify.ID(URIId))
			jsonData, err := json.Marshal(album)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = w.Write(jsonData)
		case "artist":
			artist, _ := Client.GetArtist(context.Background(), spotify.ID(URIId))
			jsonData, err := json.Marshal(artist)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = w.Write(jsonData)
		case "playlist":
			playlist, _ := Client.GetPlaylist(context.Background(), spotify.ID(URIId))
			jsonData, err := json.Marshal(playlist)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = w.Write(jsonData)
		}
	})

	// 获取闹钟列表
	m.HandleFunc("/api/alarmList", func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := json.Marshal(AlarmList)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonData)
	})

	// 获取设备列表
	m.HandleFunc("/api/deviceList", func(w http.ResponseWriter, r *http.Request) {
		deviceList, _ := Client.PlayerDevices(context.Background())
		jsonData, err := json.Marshal(deviceList)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(jsonData)
	})

	// 添加闹钟
	m.HandleFunc("/api/alarm/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		newAlarmData := r.Form

		/***************
		 * 数字转换
		 */
		_volume, _ := strconv.Atoi(newAlarmData.Get("Volume"))
		DataVolume := _volume
		_hour, _ := strconv.Atoi(newAlarmData.Get("TrigTime[Hour]"))
		DataTrigTimeHour := _hour
		_minute, _ := strconv.Atoi(newAlarmData.Get("TrigTime[Minute]"))
		DataTrigTimeMinute := _minute

		/***************
		 * 校验数据
		 */
		if DataVolume < 0 || DataVolume > 100 {
			http.Error(w, "Volume must be between 0 and 100", http.StatusBadRequest)
			return
		}
		if DataTrigTimeHour < 0 || DataTrigTimeHour > 23 {
			http.Error(w, "Hour must be between 0 and 23", http.StatusBadRequest)
			return
		}
		if DataTrigTimeMinute < 0 || DataTrigTimeMinute > 59 {
			http.Error(w, "Minute must be between 0 and 59", http.StatusBadRequest)
			return
		}

		var tracks []Track
		err = json.Unmarshal([]byte(newAlarmData.Get("Tracks")), &tracks)
		if err != nil {
			return
		}
		AlarmList = append(AlarmList, Alarm{
			Id:   newAlarmData.Get("Id"),
			Name: newAlarmData.Get("Name"),
			Device: PlaybackDevice{
				Id:          newAlarmData.Get("Device[Id]"),
				DisplayName: newAlarmData.Get("Device[DisplayName]"),
			},
			Tracks: tracks,
			Volume: DataVolume,
			Option: PlaybackOptions{0, false},
			TrigTime: AlarmTime{
				Hour:   DataTrigTimeHour,
				Minute: DataTrigTimeMinute,
			},
			Enabled: true,
		})

		SaveAlarmData(&AlarmList)
		_, _ = w.Write([]byte("OK"))
	})

	// 开关闹钟
	m.HandleFunc("/api/alarm/switch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		alarmId := r.Form.Get("Id")
		toStatus := r.Form.Get("Enabled")

		for i, alarm := range AlarmList {
			if alarm.Id == alarmId {
				if toStatus == "true" {
					AlarmList[i].Enabled = true
				} else {
					AlarmList[i].Enabled = false
				}
				log.Println("Changed ", alarmId, AlarmList[i].Enabled)
				break
			}
		}

		SaveAlarmData(&AlarmList)
		_, _ = w.Write([]byte("OK"))
	})

	// 更新闹钟
	//m.HandleFunc("/api/alarm/modify", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != http.MethodPost {
	//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	//		return
	//	}
	//
	//	// Parse form
	//	err := r.ParseForm()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(r.Form)
	//})

	/*****************
	 * 启动服务
	 */
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}
