package model

import "time"

// Usually I do custom inputs, but for simplicity I decided to use one model

type Order struct {
	ID        string    `json:"id"`
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

type RoomAvailability struct {
	ID      string    `json:"id"`
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Date    time.Time `json:"date"`
	Quota   int32     `json:"quota"`
}
