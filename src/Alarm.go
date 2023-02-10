package main

import (
	"encoding/json"
	"log"
	"os"
)

type Alarm struct {
	Id       string          `json:"Id"`
	Name     string          `json:"Name"`
	Device   PlaybackDevice  `json:"Device"`
	Tracks   []Track         `json:"Tracks"`
	Volume   int             `json:"Volume"`
	Option   PlaybackOptions `json:"Option"`
	TrigTime AlarmTime       `json:"TrigTime"`
	Enabled  bool            `json:"Enabled"`
}

type AlarmTime struct {
	Hour   int `json:"Hour"`
	Minute int `json:"Minute"`
}

type PlaybackDevice struct {
	Id          string `json:"Id"`
	DisplayName string `json:"DisplayName"`
}

type PlaybackOptions struct {
	StartPosition int  `json:"StartPosition"`
	RandomOrder   bool `json:"RandomOrder"`
}

type AlarmListStruct struct {
	Alarms []Alarm `json:"alarms"`
}

type Track struct {
	DisplayName string `json:"DisplayName"`
	Id          string `json:"Id"`
	URI         string `json:"URI"`
}

func InitAlarmDataFile() []Alarm {
	// file does not exist, create empty file
	if _, err := os.Stat("alarms.json"); os.IsNotExist(err) {
		f, err := os.Create("alarms.json")
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString("{\"alarms\": []}")
		_ = f.Close()

		return []Alarm{}
	}

	// file exists, load data
	b, err := os.ReadFile("alarms.json")
	if err != nil {
		log.Fatal(err)
	}

	var alarmList AlarmListStruct
	err = json.Unmarshal(b, &alarmList)
	if err != nil {
		log.Fatal(err)
	}
	return alarmList.Alarms
}

func SaveAlarmData(alarmList *[]Alarm) {
	CurrentAlarmListStruct := &AlarmListStruct{
		Alarms: *alarmList,
	}

	// 2 indent spaces
	b, err := json.MarshalIndent(CurrentAlarmListStruct, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("alarms.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
