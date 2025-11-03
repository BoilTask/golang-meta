package metasystem

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	metaerror "meta/meta-error"
	"time"
)

func GetCpuUsage() (float64, error) {
	cpuPercent := 0.0
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) > 0 {
		cpuPercent = percentages[0]
	}
	return cpuPercent, nil
}

func GetVirtualMemory() (uint64, uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return vmStat.Used, vmStat.Total, nil
}

func GetAvgMessage() (string, error) {
	loadStat, err := load.Avg()
	if err != nil {
		return "", metaerror.Wrap(err)
	}
	return fmt.Sprintf("1min %.2f, 5min %.2f, 15min %.2f", loadStat.Load1, loadStat.Load5, loadStat.Load15), nil
}
