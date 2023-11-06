package generator

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

var instance *sonyflake.Sonyflake
var once sync.Once

func GetId() int {
	once.Do(func() {
		snowflake, err := sonyflake.New(sonyflake.Settings{
			// Friday, September 1, 2023 12:00:00 AM GMT+07:00
			StartTime: time.Unix(1693501200, 0),
			// The lower 16 bits of the private IP address
			MachineID: nil,
		})
		if err != nil {
			panic("cannot create sonyflake instance")
		}
		instance = snowflake
	})

	id, err := instance.NextID()
	if err != nil {
		panic(fmt.Sprintf(`cannot generate sonyflake id: %s`, err.Error()))
	}
	return int(id)
}
