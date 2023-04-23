package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asher/goDocker/types"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

func GetDiskServices(path string) disk.UsageStat {
	diskInfo, _ := disk.Usage(path)
	return *diskInfo
}

func GetPcInfoServices() host.InfoStat {
	hostinfo, _ := host.Info()
	return *hostinfo
}

func GetTemperatureStat() []host.TemperatureStat {
	tempetature, _ := host.SensorsTemperatures()

	return tempetature

}

func GetTotalProcesses() []types.ProcessList {

	// infoStat, _ := host.Info()

	// fmt.Printf("Total processes: %d\n", infoStat.Procs)

	// miscStat, _ := load.Misc()
	// fmt.Printf("Running processes: %d\n", miscStat.ProcsRunning)

	var res []types.ProcessList

	ps, err := process.Processes()
	if err != nil {
		fmt.Printf("d1")
	}

	for _, v := range ps {
		var t types.ProcessList
		t.Id = v.Pid
		t.Name, err = v.Name()
		t.Name, _ = Truncate(t.Name, 30)
		if err != nil {
			t.Name = "err"
			// fmt.Print(err)
		}
		t.CpuPercent, err = v.CPUPercent()
		// fmt.Sprintf("%.4f", float32( t.CpuPercent))
		if err != nil {
			t.CpuPercent = 0.0
			// fmt.Printf("d3")
		}

		res = append(res, t)

	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CpuPercent > res[j].CpuPercent
	})

	// fmt.Println("--->", len(res))
	return res
}
func Truncate(text string, width int) (string, error) {
	if width < 0 {
		return "", fmt.Errorf("invalid width size")
	}

	r := []rune(text)
	trunc := r[:width]
	return string(trunc), nil
}
func GetStartCpu() types.Cpu {
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()
	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	var res types.Cpu

	res.CpuUsage = cpuUsage
	res.Busy = totalTicks - idleTicks
	res.Total = totalTicks
	return res
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func GetLocalIP() types.Ip {
	addrs, err := net.InterfaceAddrs()
	var ip types.Ip
	if err != nil {
		ip.LocalIp = "error"
		return ip
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				ip.LocalIp = ipnet.IP.String()
				return ip
			}
		}
	}
	ip.LocalIp = "error"
	return ip
}

func parseLine(raw string) (key string, value int) {
	// fmt.Println(raw)
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	keyValue := strings.Split(text, ":")
	return keyValue[0], toInt(keyValue[1])
}

func toInt(raw string) int {
	if raw == "" {
		return 0
	}
	res, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return res
}
func ReadMemoryStats() types.Memory {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewScanner(file)
	scanner := bufio.NewScanner(file)
	res := types.Memory{}
	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "MemTotal":
			res.MemTotal = float64(value)
		case "MemFree":
			res.MemFree = float64(value)
		case "MemAvailable":
			res.MemAvailable = float64(value)
		}
	}
	return res
}
