package steamcmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"rust_updater/internal/config"
)

func UpdateRustServer(workDir string) error {
	steamCmdPath := filepath.Join(workDir, config.DownloadDir, "steamcmd.exe")
	if _, err := os.Stat(steamCmdPath); os.IsNotExist(err) {
		return fmt.Errorf("steamcmd.exe not found at: %s", steamCmdPath)
	}

	fmt.Println("[" + config.AppName + "]    ┌─ SteamCMD Update Details ───────────────")
	fmt.Println("[" + config.AppName + "]    │ ◉ Running steamcmd.exe with arguments...")

	args := []string{
		"+force_install_dir", workDir,
		"+login", "anonymous",
		"+app_update", config.RustAppID,
		"+quit",
	}

	fmt.Println("["+config.AppName+"]    │ ◉ Arguments:", args)

	cmd := exec.Command(steamCmdPath, args...)
	cmd.Dir = workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start steamcmd: %w", err)
	}

	// Читаем stderr асинхронно
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("[%s]    │ [STDERR] %s\n", config.AppName, scanner.Text())
		}
	}()

	// Читаем stdout синхронно
	success := false
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("[%s]    │ %s\n", config.AppName, line)

		if strings.Contains(line, "Success! App '"+config.RustAppID+"'") || strings.Contains(line, "Update complete") {
			success = true
		}
	}

	// Проверка завершения
	err = cmd.Wait()
	if err != nil {
		if success {
			fmt.Printf("[%s]    │ ⚠ SteamCMD exited with warning: %v\n", config.AppName, err)
			fmt.Printf("[%s]    │ ✓ Rust server was updated successfully despite warning\n", config.AppName)
		} else {
			return fmt.Errorf("steamcmd exited with error: %w", err)
		}
	} else {
		fmt.Printf("[%s]    │ ✓ Rust server updated successfully\n", config.AppName)
	}

	fmt.Printf("[%s]    └─ SteamCMD update completed\n", config.AppName)
	return nil
}

func ReplaceFiles(workDir string) error {
	sourceDir := config.OxideSourceDir(workDir)
	targetDir := config.OxideTargetDir(workDir)

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory not found: %s", sourceDir)
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
		fmt.Printf("[%s]    │ ◉ Created target directory: RustDedicated_Data\n", config.AppName)
	}

	fmt.Printf("[%s]    ┌─ File Replacement Details ────────────\n", config.AppName)
	fmt.Printf("[%s]    │ ◉ Replacing RustDedicated_Data files...\n", config.AppName)

	count, err := copyFilesRecursively(sourceDir, targetDir)
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	fmt.Printf("[%s]    │ ✓ Successfully replaced %d files\n", config.AppName, count)
	fmt.Printf("[%s]    └─ File replacement completed\n", config.AppName)

	return nil
}

func copyFilesRecursively(sourceDir, targetDir string) (int, error) {
	fileCount := 0

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		targetFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, sourceFile)
		if err == nil {
			fileCount++
			fmt.Printf("[%s]    │ · Copied: %s\n", config.AppName, relPath)
		}
		return err
	})

	return fileCount, err
}
