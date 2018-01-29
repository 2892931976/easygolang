package pprof

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"
	"errors"
)

func PPMem(suffix string) error {
	fmt.Printf("PPMem start. suffix: [%v]\n", suffix)

	for _, name := range []string{
		"heap", "block", "goroutine", "threadcreate",
	} {
		fp, err := os.OpenFile(fmt.Sprintf("%v.pprof.%v", name, suffix), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return errors.New(fmt.Sprintf("PPMem failed. suffix: [%v], err: [%v]", suffix, err))
		}
		if err := pprof.Lookup(name).WriteTo(fp, 1); err != nil {
			fp.Close()
			return errors.New(fmt.Sprintf("PPMem failed. suffix: [%v], err: [%v]", suffix, err))
		}
		fp.Close()
	}

	return nil
}

func PPCpu(suffix string, duration time.Duration) error {
	fp, err := os.OpenFile(fmt.Sprintf("cpu.pprof.%v", suffix), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("PPCpu failed. suffix: [%v], duration: [%v], err: [%v]", suffix, duration, err))
	}
	defer fp.Close()

	if err := pprof.StartCPUProfile(fp); err != nil {
		return errors.New(fmt.Sprintf("PPCpu failed. suffix: [%v], duration: [%v], err: [%v]", suffix, duration, err))
	}
	time.Sleep(duration)
	pprof.StopCPUProfile()

	return nil
}

// "mem [suffix=yyyymmddhh]" 生成内存 pprof 文件
// "cpu [duration=60s] [suffix=yyyymmddhh]" 生成 cpu pprof 文件
func PPCmd(command string) error {
	fields := strings.Split(command, " ")
	cmdType := fields[0]

	switch cmdType {
	case "mem":
		suffix := time.Now().Format("200601021504")
		if len(fields) >= 2 {
			suffix = fields[1]
		}

		return PPMem(suffix)
	case "cpu":
		duration := 60 * time.Second
		if len(fields) >= 2 {
			var err error
			duration, err = time.ParseDuration(fields[1])
			if err != nil {
				return errors.New(fmt.Sprintf("PPCmd failed. command: [%v], err: [%v]", command, err))
			}
		}

		suffix := time.Now().Format("200601021504")
		if len(fields) >= 3 {
			suffix = fields[2]
		}

		return PPCpu(suffix, duration)
	default:
		return errors.New(fmt.Sprintf("PPCmd failed. invalid command: [%v]", command))
	}
}
