package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type StatusReport struct {
	Version               string  `json:"Version"`
	Devices               int     `json:"Devices"`
	ConnectedUsers        int     `json:"Connected_Users"`
	UATMessagesLastMinute int     `json:"UAT_messages_last_minute"`
	UATMessagesMax        int     `json:"UAT_messages_max"`
	ESMessagesLastMinute  int     `json:"ES_messages_last_minute"`
	ESMessagesMax         int     `json:"ES_messages_max"`
	GPSSatellitesLocked   int     `json:"GPS_satellites_locked"`
	GPSConnected          bool    `json:"GPS_connected"`
	GPSSolution           string  `json:"GPS_solution"`
	Uptime                int     `json:"Uptime"`
	CPUTemp               float64 `json:"CPUTemp"`
}

type SituationReport struct {
	LastFixSinceMidnightUTC  float64 `json:"LastFixSinceMidnightUTC"`
	Lat                      float64 `json:"Lat"`
	Lng                      float64 `json:"Lng"`
	Quality                  int     `json:"Quality"`
	HeightAboveEllipsoid     float64 `json:"HeightAboveEllipsoid"`
	GeoidSep                 float64 `json:"GeoidSep"`
	Satellites               int     `json:"Satellites"`
	SatellitesTracked        int     `json:"SatellitesTracked"`
	SatellitesSeen           int     `json:"SatellitesSeen"`
	Accuracy                 float64 `json:"Accuracy"`
	NACp                     int     `json:"NACp"`
	Alt                      float64 `json:"Alt"`
	AccuracyVert             float64 `json:"AccuracyVert"`
	GPSVertVel               float64 `json:"GPSVertVel"`
	LastFixLocalTime         string  `json:"LastFixLocalTime"`
	TrueCourse               int     `json:"TrueCourse"`
	GroundSpeed              int     `json:"GroundSpeed"`
	LastGroundTrackTime      string  `json:"LastGroundTrackTime"`
	GPSTime                  string  `json:"GPSTime"`
	LastGPSTimeTime          string  `json:"LastGPSTimeTime"`
	LastValidNMEAMessageTime string  `json:"LastValidNMEAMessageTime"`
	LastValidNMEAMessage     string  `json:"LastValidNMEAMessage"`
	Temp                     int     `json:"Temp"`
	PressureAlt              int     `json:"Pressure_alt"`
	LastTempPressTime        string  `json:"LastTempPressTime"`
	Pitch                    int     `json:"Pitch"`
	Roll                     int     `json:"Roll"`
	GyroHeading              int     `json:"Gyro_heading"`
	LastAttitudeTime         string  `json:"LastAttitudeTime"`
}

type TrafficReport struct {
	IcaoAddr      int       `json:"Icao_addr"`
	OnGround      bool      `json:"OnGround"`
	Lat           float64   `json:"Lat"`
	Lng           float64   `json:"Lng"`
	PositionValid bool      `json:"Position_valid"`
	Alt           int       `json:"Alt"`
	Track         int       `json:"Track"`
	Speed         int       `json:"Speed"`
	SpeedValid    bool      `json:"Speed_valid"`
	Vvel          int       `json:"Vvel"`
	Tail          string    `json:"Tail"`
	LastSeen      time.Time `json:"Last_seen"`
	LastSource    int       `json:"Last_source"`
}

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func wsTraffic(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	defer ws.Close()

	for {
		<-time.After(time.Second)
		traffic := TrafficReport{
			IcaoAddr:      2837120,
			OnGround:      false,
			Lat:           -27.57,  // S27°34.22' - YBAF
			Lng:           152.997, // E152°59.83'
			PositionValid: true,
			Alt:           3400,
			Track:         158,
			Speed:         92,
			SpeedValid:    true,
			Vvel:          0,
			Tail:          "",
			LastSeen:      time.Now(),
			LastSource:    2,
		}
		ws.WriteJSON(&traffic)
	}

}

func getStatus() *StatusReport {
	return &StatusReport{
		Version:               "v1.0r1",
		Devices:               0,
		ConnectedUsers:        1,
		ESMessagesLastMinute:  100,
		ESMessagesMax:         500,
		UATMessagesLastMinute: 0,
		UATMessagesMax:        0,
		GPSSatellitesLocked:   10,
		GPSConnected:          true,
		GPSSolution:           "3D GPS",
		Uptime:                227068,
		CPUTemp:               42.236,
	}
}

func serveStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, _ := json.Marshal(getStatus())
	w.Write([]byte(status))
}

func wsStatus(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	defer ws.Close()

	for {
		<-time.After(time.Second)
		ws.WriteJSON(getStatus())
	}
}

func getSituation() *SituationReport {
	return &SituationReport{
		LastFixSinceMidnightUTC: 23550.8,
		Lat:                      -27.5,
		Lng:                      152.97,
		Quality:                  1,
		HeightAboveEllipsoid:     208.35303,
		GeoidSep:                 125.98426,
		Satellites:               6,
		SatellitesTracked:        8,
		SatellitesSeen:           8,
		Accuracy:                 4.4,
		NACp:                     10,
		Alt:                      82.36877,
		AccuracyVert:             5.2,
		GPSVertVel:               -0.36417323,
		LastFixLocalTime:         "0001-01-01T00:04:11.62Z",
		TrueCourse:               0,
		GroundSpeed:              0,
		LastGroundTrackTime:      "0001-01-01T00:04:11.62Z",
		GPSTime:                  "2017-07-17T06:32:29.6Z",
		LastGPSTimeTime:          "0001-01-01T00:04:10.42Z",
		LastValidNMEAMessageTime: "0001-01-01T00:04:11.62Z",
		LastValidNMEAMessage:     "$PUBX,00,063230.80,2731.72269,S,15258.06057,E,63.506,G3,2.2,2.6,0.523,0.00,0.111,,1.22,1.75,1.14,6,0,0*46",
		Temp:                     0,
		PressureAlt:              0,
		LastTempPressTime:        "0001-01-01T00:00:00Z",
		Pitch:                    0,
		Roll:                     0,
		GyroHeading:              0,
		LastAttitudeTime:         "0001-01-01T00:00:00Z",
	}
}

func serveSituation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	situation, _ := json.Marshal(getSituation())
	w.Write([]byte(situation))
}

func wsSituation(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	defer ws.Close()

	for {
		<-time.After(time.Second)
		ws.WriteJSON(getSituation())
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		"initial",
	}
	homeTempl.Execute(w, &v)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/getStatus", serveStatus)
	http.HandleFunc("/getSituation", serveSituation)
	http.HandleFunc("/status", wsStatus)
	http.HandleFunc("/situation", wsSituation)
	http.HandleFunc("/traffic", wsTraffic)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws:
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
