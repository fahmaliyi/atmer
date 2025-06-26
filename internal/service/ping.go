package service

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func Ping(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}
	err := cmd.Run()
	return err == nil
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func GetModemIP(ip string) string {
	if strings.HasPrefix(ip, "172.") {
		return ip
	}

	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ip
	}

	last := atoi(parts[3])
	third := atoi(parts[2])

	if last > 0 {
		parts[3] = itoa(last - 1)
	} else if third > 0 {
		parts[2] = itoa(third - 1)
		parts[3] = "0"
	}
	return strings.Join(parts, ".")
}
