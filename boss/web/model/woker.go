package model

type Worker struct {
	Addr string `json:"address"`
	Name string `json:"name"`
	Ping int64  `json:"ping"`
}
