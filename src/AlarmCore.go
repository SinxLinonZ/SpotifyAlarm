package main

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
	"time"
)

func startAlarmService() {

	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				nowHour := now.Hour()
				nowMinute := now.Minute()

				//for _, alarm := range AlarmList {
				for i := 0; i < len(AlarmList); i++ {
					alarm := &AlarmList[i]
					if !alarm.Enabled {
						continue
					}

					if alarm.TrigTime.Hour == nowHour && alarm.TrigTime.Minute == nowMinute {
						ID := spotify.ID(alarm.Device.Id)
						URIs := make([]spotify.URI, 0)
						for _, track := range alarm.Tracks {
							URIs = append(URIs, spotify.URI(track.URI))
						}

						err := Client.PlayOpt(context.Background(), &spotify.PlayOptions{
							DeviceID: &ID,
							URIs:     URIs,
						})
						if err != nil {
							log.Println(err)
						}
						err = Client.Volume(context.Background(), alarm.Volume)
						if err != nil {
							log.Println(err)
						}

						alarm.Enabled = false
						SaveAlarmData(&AlarmList)
					}
				}
			}
		}
	}()
}
