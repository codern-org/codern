package platform

import (
	"time"

	"github.com/sony/sonyflake"
)

func NewSonyFlake() (*sonyflake.Sonyflake, error) {
	return sonyflake.New(sonyflake.Settings{
		// Friday, September 1, 2023 12:00:00 AM GMT+07:00
		StartTime: time.Unix(1693501200, 0),
		// The lower 16 bits of the private IP address
		MachineID: nil,
	})
}
