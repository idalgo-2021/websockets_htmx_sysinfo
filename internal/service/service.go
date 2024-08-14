package service

import (
	"fmt"

	"websockets_htmx_sysinfo/internal/hardware"
)

type HardwareService struct{}

func NewHardwareService() *HardwareService {
	return &HardwareService{}
}

func (hs *HardwareService) GetSystemSection() (string, error) {
	systemInfo, err := hardware.GetSystemInfo()
	if err != nil {
		return "", fmt.Errorf("failed to get system info: %w", err)
	}
	return hardware.FormatSystemInfo(systemInfo), nil
}

func (hs *HardwareService) GetDiskSection() (string, error) {
	diskInfo, err := hardware.GetDiskInfo()
	if err != nil {
		return "", fmt.Errorf("failed to get disk info: %w", err)
	}
	return hardware.FormatDiskInfo(diskInfo), nil
}

func (hs *HardwareService) GetCPUSection() (string, error) {
	cpuInfo, err := hardware.GetCPUInfo()
	if err != nil {
		return "", fmt.Errorf("failed to get CPU info: %w", err)
	}
	return hardware.FormatCPUInfo(cpuInfo), nil
}
