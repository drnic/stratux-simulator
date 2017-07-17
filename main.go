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

func serveTraffic(w http.ResponseWriter, r *http.Request) {
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
			Lat:           42.193336,
			Lng:           -83.92136,
			PositionValid: true,
			Alt:           3400,
			Track:         9,
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

func serveStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusReport := StatusReport{
		Version:               "v0.5b1",
		Devices:               0,
		ConnectedUsers:        1,
		ESMessagesLastMinute:  100,
		ESMessagesMax:         500,
		UATMessagesLastMinute: 0,
		UATMessagesMax:        0,
		GPSSatellitesLocked:   10,
		GPSConnected:          true,
		Uptime:                227068,
		CPUTemp:               42.236,
	}
	status, _ := json.Marshal(statusReport)
	w.Write([]byte(status))
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
	http.HandleFunc("/traffic", serveTraffic)
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
