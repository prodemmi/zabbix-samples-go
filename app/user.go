package main

import "time"

type User struct {
	ID          int
	Username    string
	OrderCount  int64
	LastOrderAt time.Time
	MaxPayment  int64
	IsOnline    bool
	IP          string
	Country     string
	LoginAt     time.Time
	CreatedAt   time.Time
}
