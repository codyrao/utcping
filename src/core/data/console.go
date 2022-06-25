package data

import (
	"fmt"
	"time"
)

type Console struct {
}

func (ptr *Console) WriteStat(stat *Statistics) error {
	if nil == stat {
		return fmt.Errorf("stat is nil")
	}
	fmt.Println("----------------Statistics----------------")
	fmt.Println()
	fmt.Printf("%s:%d  timeout: %ds   print_time: %s\n", stat.IP, stat.Port, stat.Timeout, time.Now().Format(TimeFormat))
	fmt.Printf("total: %d success: %d fail: %d packet_lose_ratio: %s", stat.Seq, len(stat.Success), len(stat.Fail), stat.PacketLoseRatio)
	fmt.Println()
	fmt.Println()
	fmt.Println("----Fail:")
	for _, v := range stat.Fail {
		fmt.Printf("time:%s seq:%d reason:%s\n", v.Time.Format(TimeFormat), v.Seq, v.Err)
	}
	fmt.Println()
	fmt.Println("----Success:")
	if len(stat.Success) != 0 {
		fmt.Printf("delay_min:%.2fms delay_max:%.2fms delay_avg:%.2fms,delay_stdev:%.2fms\n", stat.DelayMin, stat.DelayMax, stat.DelayAvg, stat.DelayStandardDeviation)
	}

	return nil
}
