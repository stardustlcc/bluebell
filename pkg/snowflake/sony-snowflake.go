package snowflake

// import (
// 	"fmt"
// 	"time"

// 	"github.com/sony/sonyflake"
// )

// var (
// 	sonyFlake     *sonyflake.Sonyflake
// 	sonyMachineID uint16
// )

// func getMachineID() (uint16, error) {
// 	return sonyMachineID, nil
// }

// func Init(startTime string, machineId uint16) (err error) {
// 	sonyMachineID = machineId
// 	var st time.Time
// 	st, err = time.Parse("2006-01-02", startTime)
// 	if err != nil {
// 		return err
// 	}
// 	settings := sonyflake.Settings{
// 		StartTime: st,
// 		MachineID: getMachineID,
// 	}
// 	sonyFlake = sonyflake.NewSonyflake(settings)
// 	return
// }

// func GetID() (id uint64, err error) {
// 	if sonyFlake == nil {
// 		err = fmt.Errorf("snoy flake not inited")
// 		return
// 	}
// 	id, err = sonyFlake.NextID()
// 	return
// }

// func main() {
// 	if err := Init("2020-07-01", 1); err != nil {
// 		fmt.Printf("init failed err:%v\n", err)
// 		return
// 	}
// 	id, errmsg := GetID()
// 	if errmsg != nil {
// 		fmt.Println(errmsg)
// 	}
// 	fmt.Println(id)
// }
