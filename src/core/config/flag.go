package config

import (
	"core/utils"
	"flag"
	"fmt"
	"sync"
)

var (
	once       sync.Once
	globalFlag *Flag
)

type Flag struct {
	IP                 string
	Port               int64
	Number             int
	Timeout            int64
	Interval           int64
	DataSource         string
	StatisticsInterVal int64
}

func (ptr *Flag) Init() error {

	err := ptr.Input()
	if nil != err {
		return err
	}

	flag.Parse()

	err = ptr.Check()
	if nil != err {
		return fmt.Errorf("%s\nPlease execute 'tcping -h' for help documentation", err)
	}

	return nil
}

func (ptr *Flag) Input() error {
	flag.StringVar(&ptr.IP, "i", "", "target ipv4/ipv6")
	flag.Int64Var(&ptr.Port, "p", 80, "target port")
	flag.Int64Var(&ptr.Timeout, "t", 3, "tcping timeout(s)")
	flag.IntVar(&ptr.Number, "n", 0, "numbers of tcping 0=unlimited (default unlimited)")
	flag.Int64Var(&ptr.Interval, "s", 1000, "tcping interval(ms)")
	flag.StringVar(&ptr.DataSource, "d", "console", "datasource for tcping statistics output,supported[console,csv]")
	flag.Int64Var(&ptr.StatisticsInterVal, "o", 60, "tcping statistics time interval(min)")
	return nil
}

func (ptr *Flag) Check() error {
	if !utils.IsIP(ptr.IP) {
		return fmt.Errorf("target[%s] ip is illegal", ptr.IP)
	}

	if !utils.IsPort(ptr.Port) {
		return fmt.Errorf("port[%d] is illegal, 0 < port < 65536", ptr.Port)
	}

	if !utils.IsLegalNumber(ptr.Number) {
		return fmt.Errorf("number[%d] is illegal, 0 < number < 2^31-1", ptr.Number)
	}

	if !utils.IsLegalTimeout(ptr.Timeout) {
		return fmt.Errorf("timeout[%d] is illegal,0 < timeouts < 2^31-1", ptr.Timeout)
	}

	if !utils.IsLegalInterval(ptr.Interval) {
		return fmt.Errorf("time interval[%d] is illegal,0 < interval < 2^31-1", ptr.Interval)
	}

	if !utils.IsLegalStatisticsInterVal(ptr.StatisticsInterVal) {
		return fmt.Errorf("statistics[%d] time interVal is illegal,0 < interval < 2^31-1", ptr.StatisticsInterVal)
	}

	return nil
}

func GetGlobalFlag() *Flag {
	once.Do(func() {
		globalFlag = new(Flag)
	})
	return globalFlag
}
