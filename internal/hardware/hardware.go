package hardware

import (
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

type SystemInfo struct {
	OS            string
	Platform      string
	Hostname      string
	Processes     uint64
	TotalMemory   uint64
	FreeMemory    uint64
	UsedMemoryPct float64
}

type DiskInfo struct {
	TotalSpace   uint64
	UsedSpace    uint64
	FreeSpace    uint64
	UsedSpacePct float64
}

type CPUInfo struct {
	ModelName  string
	Family     string
	SpeedMHz   float64
	CoresUsage []float64
}

func GetSystemInfo() (SystemInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return SystemInfo{}, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		OS:            runtime.GOOS,
		Platform:      hostStat.Platform,
		Hostname:      hostStat.Hostname,
		Processes:     hostStat.Procs,
		TotalMemory:   vmStat.Total,
		FreeMemory:    vmStat.Free,
		UsedMemoryPct: vmStat.UsedPercent,
	}, nil
}

func GetDiskInfo() (DiskInfo, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskInfo{}, err
	}

	return DiskInfo{
		TotalSpace:   diskStat.Total,
		UsedSpace:    diskStat.Used,
		FreeSpace:    diskStat.Free,
		UsedSpacePct: diskStat.UsedPercent,
	}, nil
}

func GetCPUInfo() (CPUInfo, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return CPUInfo{}, err
	}

	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return CPUInfo{}, err
	}

	return CPUInfo{
		ModelName:  cpuStat[0].ModelName,
		Family:     cpuStat[0].Family,
		SpeedMHz:   cpuStat[0].Mhz,
		CoresUsage: percentage,
	}, nil
}

func FormatSystemInfo(info SystemInfo) string {
	html := "<div id='system-data' class='system-data'><table class='table table-striped table-hover table-sm'><tbody>"

	html += "<tr><td>Operating System:</td> <td><i class='fa fa-brands fa-linux'></i> " + info.OS + "</td></tr>"
	html += "<tr><td>Platform:</td><td> <i class='fa fa-brands fa-fedora'></i> " + info.Platform + "</td></tr>"
	html += "<tr><td>Hostname:</td><td>" + info.Hostname + "</td></tr>"
	html += "<tr><td>Number of processes running:</td><td>" + strconv.FormatUint(info.Processes, 10) + "</td></tr>"
	html += "<tr><td>Total memory:</td><td>" + strconv.FormatUint(info.TotalMemory/megabyteDiv, 10) + " MB</td></tr>"
	html += "<tr><td>Free memory:</td><td>" + strconv.FormatUint(info.FreeMemory/megabyteDiv, 10) + " MB</td></tr>"
	html += "<tr><td>Percentage used memory:</td><td>" + strconv.FormatFloat(info.UsedMemoryPct, 'f', 2, 64) + "%</td></tr></tbody></table>"
	html += "</div>"

	html += "</tbody></table></div>"
	return html
}

func FormatDiskInfo(info DiskInfo) string {
	html := "<div id='disk-data' class='disk-data'><table class='table table-striped table-hover table-sm'><tbody>"

	html += "<tr><td>Total disk space:</td><td>" + strconv.FormatUint(info.TotalSpace/gigabyteDiv, 10) + " GB</td></tr>"
	html += "<tr><td>Used disk space:</td><td>" + strconv.FormatUint(info.UsedSpace/gigabyteDiv, 10) + " GB</td></tr>"
	html += "<tr><td>Free disk space:</td><td>" + strconv.FormatUint(info.FreeSpace/gigabyteDiv, 10) + " GB</td></tr>"
	html += "<tr><td>Percentage disk space usage:</td><td>" + strconv.FormatFloat(info.UsedSpacePct, 'f', 2, 64) + "%</td></tr>"
	html += "</div>"

	html += "</tbody></table></div>"

	return html
}

func FormatCPUInfo(info CPUInfo) string {
	html := "<div id='cpu-data' class='cpu-data'><table class='table table-striped table-hover table-sm'><tbody>"

	html += "<tr><td>Model Name:</td><td>" + info.ModelName + "</td></tr>"
	html += "<tr><td>Family:</td><td>" + info.Family + "</td></tr>"
	html += "<tr><td>Speed:</td><td>" + strconv.FormatFloat(info.SpeedMHz, 'f', 2, 64) + " MHz</td></tr>"

	firstCpus := info.CoresUsage[:len(info.CoresUsage)/2]
	secondCpus := info.CoresUsage[len(info.CoresUsage)/2:]

	html += "<tr><td>Cores: </td><td><div class='row mb-4'><div class='col-md-6'><table class='table table-sm'><tbody>"
	for idx, cpupercent := range firstCpus {
		html += "<tr><td>CPU [" + strconv.Itoa(idx) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}
	html += "</tbody></table></div><div class='col-md-6'><table class='table table-sm'><tbody>"
	for idx, cpupercent := range secondCpus {
		html += "<tr><td>CPU [" + strconv.Itoa(idx+8) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}
	html += "</tbody></table></div></div></td></tr></tbody></table></div>"

	html += "</tbody></table></div>"

	return html
}
