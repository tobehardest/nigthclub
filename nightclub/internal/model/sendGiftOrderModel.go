package model

import "time"

type SendGiftOrder struct {
	Id         int64     `json:"id"`
	Sn         string    `json:"sn"`
	FromId     int64     `json:"from_id"`
	ToId       int64     `json:"to_id"`
	GiftPrice  int64     `json:"gift_price"`
	Version    int64     `json:"version"`
	DelState   int64     `json:"del_state"`
	TradeState int64     `json:"trade_state"`
	TradeCode  string    `json:"trade_code"`
	CreatTime  time.Time `json:"creat_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
