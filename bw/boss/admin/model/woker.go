package model

type Worker struct {
	Addr string `json:"address" gorm:"addr"`
	Name string `json:"name"gorm:"name"`
	Thumb string `json:"thumb" gorm:"thumb"`
	Ping int64  `json:"ping" gorm:"ping"`
	Status bool	`json:"status" gorm:"status"`
}
