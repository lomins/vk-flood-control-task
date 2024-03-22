package floodcontrol

import "time"

type Payload struct {
	LastCallTime time.Time `json:"last_call_time"`
	CallCount    uint      `json:"call_count"`
}
