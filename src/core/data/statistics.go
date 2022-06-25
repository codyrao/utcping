package data

import (
	"core/config"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	TimeFormat            = "2006-01-02 15:04:05"
	DefaultSuccessListCap = 100
)

var (
	once             sync.Once
	globalStatistics *Statistics
)

type Writer interface {
	WriteStat(stat *Statistics) error
}

type Statistics struct {
	IP                     string
	Port                   int64
	Timeout                int64
	Fail                   []*Fail
	Seq                    int64
	Success                []*Success
	DelayTotal             float64
	DelayAvg               float64
	DelayMax               float64
	DelayMin               float64
	DelayStandardDeviation float64
	PacketLoseRatio        string
	LastOutputTime         time.Time
}

func GetGlobalStatistics() *Statistics {
	once.Do(func() {
		cap := DefaultSuccessListCap
		if config.GetGlobalFlag().Number > 0 {
			cap = config.GetGlobalFlag().Number
		}
		globalStatistics = &Statistics{
			IP:              config.GetGlobalFlag().IP,
			Port:            config.GetGlobalFlag().Port,
			Timeout:         config.GetGlobalFlag().Timeout,
			Fail:            make([]*Fail, 0),
			Success:         make([]*Success, 0, cap),
			Seq:             0,
			LastOutputTime:  time.Now(),
			PacketLoseRatio: "0%",
			DelayMin:        math.MaxFloat64,
		}
	})
	return globalStatistics
}

type Event struct {
	Now   time.Time
	Delay float64
	Err   error
}

func (ptr *Statistics) Do(event interface{}) {
	e, ok := event.(*Event)
	if !ok {
		return
	}

	ptr.Set(e.Now, e.Delay, e.Err)

	ptr.PrintOne(e.Now, e.Delay, e.Err)

	now := time.Now()
	if ptr.LastOutputTime.Add(time.Duration(config.GetGlobalFlag().StatisticsInterVal) * time.Minute).Before(now) {
		ptr.LastOutputTime = now
		ptr.Output()
	}

}

func (ptr *Statistics) Set(now time.Time, delay float64, err error) {
	ptr.Seq++
	if nil != err {
		ptr.Fail = append(ptr.Fail, &Fail{
			Seq:  ptr.Seq,
			Err:  err,
			Time: now,
		})
		ptr.PacketLoseRatio = fmt.Sprintf("%.2f%%", float64(len(ptr.Fail)*100)/float64(ptr.Seq))
		return
	}

	ptr.Success = append(ptr.Success, &Success{
		Seq:   ptr.Seq,
		Delay: delay,
		Time:  now,
	})

	ptr.DelayTotal += delay
	ptr.DelayMax = math.Max(ptr.DelayMax, delay)
	ptr.DelayMin = math.Min(ptr.DelayMin, delay)
	ptr.DelayAvg = ptr.GetDelayAvg()
	ptr.DelayStandardDeviation = ptr.GetDelayStandardDeviation()
}

func (ptr *Statistics) Output() {
	dataSource := config.GetGlobalFlag().DataSource

	switch dataSource {
	case "csv":
		ptr.Write(new(CSV))
		ptr.Write(new(Console))
	default:
		ptr.Write(new(Console))

	}

}

func (ptr *Statistics) Write(writer Writer) {
	err := writer.WriteStat(ptr)
	if nil != err {
		fmt.Printf("statistics write fail.error:%s", err)
		return
	}
}

func (ptr *Statistics) PrintOne(now time.Time, delay float64, err error) {

	if nil != err {
		fmt.Printf("%s tcping %s:%d fail[timeout:%ds]: %s\n", now.Format(TimeFormat), ptr.IP, ptr.Port, ptr.Timeout, err)
		return
	}

	fmt.Printf("%s tcping %s:%d success: seq=%d,time=%.2fms\n", now.Format(TimeFormat), ptr.IP, ptr.Port, ptr.Seq, delay)
}

func (ptr *Statistics) GetDelayAvg() float64 {
	count := len(ptr.Success)
	if count == 0 {
		return 0
	}
	return ptr.DelayTotal / float64(count)
}

func (ptr *Statistics) GetDelayMin() float64 {
	count := len(ptr.Success)
	if count == 0 {
		return 0
	}

	success := ptr.Success
	var total float64 = 0
	for _, v := range success {
		total += v.Delay
	}

	return total / float64(count)
}

func (ptr *Statistics) GetDelayMax() float64 {
	count := len(ptr.Success)
	if count == 0 {
		return 0
	}

	success := ptr.Success
	var total float64 = 0
	for _, v := range success {
		total += v.Delay
	}

	return total / float64(count)
}

func (ptr *Statistics) GetDelayStandardDeviation() float64 {
	count := len(ptr.Success)
	if count == 0 {
		return 0
	}
	var variance float64

	delayAvg := ptr.DelayAvg
	for _, v := range ptr.Success {
		variance += (v.Delay - delayAvg) * (v.Delay - delayAvg)
	}
	variance /= float64(count)

	return math.Sqrt(variance)

}

type Fail struct {
	Seq  int64
	Err  error
	Time time.Time
}

type Success struct {
	Seq   int64
	Delay float64
	Time  time.Time
}
