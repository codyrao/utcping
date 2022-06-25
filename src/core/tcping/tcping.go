package utcping

import (
	"core/config"
	"core/data"
	"core/utils"
	"fmt"
	"net"
	"time"
)

type TCPPing struct {
}

func GetTCPing() *TCPPing {
	return new(TCPPing)
}

func (ptr *TCPPing) Do(flag *config.Flag) error {
	if nil == flag {
		return fmt.Errorf("tcping error:fail is nil")
	}
	ip := flag.IP
	port := flag.Port
	timeout := time.Duration(flag.Timeout) * time.Second
	number := flag.Number
	interval := time.Duration(flag.Interval) * time.Millisecond

	for i := 0; i < number || number == 0; i++ {
		//多协程并发探测
		go ptr.asynchronous(ip, port, timeout)
		time.Sleep(interval)
	}

	for {
		if data.GetGlobalStatistics().Seq >= int64(number) {
			break
		}
		time.Sleep(time.Second)
	}
	data.GetGlobalStatistics().Output()

	return nil
}

func (ptr *TCPPing) asynchronous(ip string, port int64, timeout time.Duration) {
	now := time.Now()
	callback := utils.PutCallbackEvent(&utils.CallbackEvent{
		F: TcpPing, Arg: &TcpPingEvent{
			Ip:      ip,
			Port:    port,
			Timeout: timeout,
		},
	})
	if nil == callback {
		panic(fmt.Errorf("callback is nil"))
	}

	res, err := callback.Get()
	eventRes, ok := res.(*TcpPingEventRes)
	if !ok {
		panic(fmt.Errorf("eventRes is nil"))
	}

	//协程队列中处理
	utils.PutEvent(&utils.Event{
		F: data.GetGlobalStatistics().Do, Arg: &data.Event{
			Now:   now,
			Delay: eventRes.Delay,
			Err:   err,
		},
	})
}

type TcpPingEvent struct {
	Ip      string
	Port    int64
	Timeout time.Duration
}

type TcpPingEventRes struct {
	Time  time.Time
	Delay float64
}

func TcpPing(event interface{}) (interface{}, error) {
	e, ok := event.(*TcpPingEvent)
	if !ok {
		return nil, fmt.Errorf("tcping event is nil")
	}
	addr := fmt.Sprintf("[%s]:%d", e.Ip, e.Port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, e.Timeout)
	delay := float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond)
	if nil != err {
		return &TcpPingEventRes{
			Time: start,
		}, err
	}
	defer conn.Close()

	return &TcpPingEventRes{
		Time:  start,
		Delay: delay,
	}, nil

}
