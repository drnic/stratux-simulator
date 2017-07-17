// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

	traffic := TrafficReport{}
	for {
		<-time.After(time.Second)
		ws.WriteJSON(&traffic)
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
	http.HandleFunc("/", serveHome)
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
                var conn = new WebSocket("ws://{{.Host}}/traffic");
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
