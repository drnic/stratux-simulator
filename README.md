# Stratux Simulator

I wanted to build a little app that watches [Stratux](https://github.com/cyoung/stratux) web socket data, so this is a simple simulator to produce traffic. Stratux API documented at https://github.com/cyoung/stratux/blob/master/notes/app-vendor-integration.md

Currently supports:

* GET /getStatus
* WS /traffic

In one terminal:

```
go run main.go
```

To watch simulated traffic:

```
wsd -url ws://localhost:8080/traffic -origin http://localhost:8080/
```

Get static response to Stratux status:

```
$ curl localhost:8080/getStatus
{"Version":"v0.5b1","Devices":0,"Connected_Users":1,"UAT_messages_last_minute":0,"UAT_messages_max":0,"ES_messages_last_minute":100,"ES_messages_max":500,"GPS_satellites_locked":10,"GPS_connected":true,"GPS_solution":"","Uptime":227068,"CPUTemp":42.236}
```
