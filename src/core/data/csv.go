package data

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

const (
	FileTimeFormat = "2006_01_02_15_04_05"
	CSVFileName    = "utcping-%s-%d-%s.csv"
)

type CSV struct {
	Contents [][][]string
}

func (ptr *CSV) WriteStat(stat *Statistics) error {

	var header [][]string

	delayMinStr := "0.00ms"
	if stat.DelayMin != math.MaxFloat64 {
		delayMinStr = fmt.Sprintf("%.2fms", stat.DelayMin)
	}

	titles := []string{"ip", "port", "timeout", "print_time", "total", "success", "fail", "packet_lose_ratio", "delay_min", "delay_max", "delay_avg", "delay_stdev"}
	titleContent := []string{stat.IP, fmt.Sprintf("%d", stat.Port), fmt.Sprintf("%ds", stat.Timeout), time.Now().Format(TimeFormat),
		fmt.Sprintf("%d", stat.Seq), fmt.Sprintf("%d", len(stat.Success)), fmt.Sprintf("%d", len(stat.Fail)), stat.PacketLoseRatio,
		delayMinStr, fmt.Sprintf("%.2fms", stat.DelayMax), fmt.Sprintf("%.2fms", stat.DelayAvg), fmt.Sprintf("%.2fms", stat.DelayStandardDeviation)}
	header = append(header, titles, titleContent)

	var content [][]string

	failHeader := []string{"Fail:"}
	failTitle := []string{"time", "seq", "reason"}
	content = append(content, []string{}, failHeader, failTitle)
	for _, v := range stat.Fail {
		content = append(content, []string{fmt.Sprintf("%s", v.Time.Format(TimeFormat)), fmt.Sprintf("%d", v.Seq), fmt.Sprintf("%s", v.Err)})
	}

	successHeader := []string{"Success:"}
	successTitle := []string{"time", "seq", "delay"}
	content = append(content, []string{}, successHeader, successTitle)
	for _, v := range stat.Success {
		content = append(content, []string{fmt.Sprintf("%s", v.Time.Format(TimeFormat)), fmt.Sprintf("%d", v.Seq), fmt.Sprintf("%.2fms", v.Delay)})
	}

	var res [][][]string
	res = append(res, header, content)
	ptr.Contents = res

	err := ptr.Write(fmt.Sprintf(CSVFileName, strings.Replace(stat.IP, ".", "_", -1), stat.Port, time.Now().Format(FileTimeFormat)))
	if nil != err {
		return err
	}

	return nil
}

func (ptr *CSV) Write(fileName string) error {
	if len(ptr.Contents) == 0 {
		return fmt.Errorf("csv content is empty")
	}
	f, err := os.Create(fmt.Sprintf("./%s", fileName))
	if nil != err {
		return err
	}
	defer f.Close()

	// 写入UTF-8 BOM
	_, err = f.WriteString("\xEF\xBB\xBF")
	if nil != err {
		return err
	}

	//创建一个新的写入文件流
	w := csv.NewWriter(f)

	for _, content := range ptr.Contents {
		err = w.WriteAll(content)
		if nil != err {
			return err
		}
	}

	w.Flush()

	return nil
}
