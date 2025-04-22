package process

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"rust_updater/internal/config"
)

func FindAndKillRustServer() error {
	dir, err := config.GetWorkDir()
	if err != nil {
		return fmt.Errorf("get cwd: %w", err)
	}
	dir = strings.ToLower(filepath.Clean(dir))

	ps := fmt.Sprintf(`Get-Process RustDedicated -ErrorAction SilentlyContinue | Where-Object { $_.Path -like "%s*" } | ForEach-Object { $_.Id }`, dir)
	out, err := exec.Command("powershell", "-Command", ps).Output()
	if err != nil {
		return fmt.Errorf("powershell error: %w", err)
	}

	pid := strings.TrimSpace(string(out))
	if pid == "" {
		fmt.Printf("[%s] RustDedicated.exe not found\n", config.AppName)
		return nil
	}

	_, err = exec.Command("taskkill", "/F", "/PID", pid).CombinedOutput()
	if err != nil {
		return fmt.Errorf("kill failed for PID %s: %w", pid, err)
	}

	fmt.Printf("[%s] RustDedicated.exe (PID %s) terminated\n", config.AppName, pid)
	return nil
}
