package updater

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"rust_updater/internal/config"
	"rust_updater/internal/downloader"
	"rust_updater/internal/extractor"
	"rust_updater/internal/process"
	"rust_updater/internal/steamcmd"
)

func Run() error {
	log.SetFlags(0)
	log.SetPrefix("[" + config.AppName + "] ")

	if err := config.InitWorkDir(); err != nil {
		return fmt.Errorf("init workdir: %w", err)
	}

	log.Println("╔══════════════════════════════════════╗")
	log.Println("║       RUST SERVER UPDATE PROCESS     ║")
	log.Println("╚══════════════════════════════════════╝")

	if err := prepare(); err != nil {
		return err
	}

	if err := stopRustServer(); err != nil {
		log.Printf("Warning: failed to stop server: %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := steamcmd.UpdateRustServer(config.ServerInstallPath); err != nil {
		return fmt.Errorf("steam update: %w", err)
	}

	if err := steamcmd.ReplaceFiles(config.ServerInstallPath); err != nil {
		return fmt.Errorf("oxide update: %w", err)
	}

	log.Println("╔══════════════════════════════════════╗")
	log.Println("║       UPDATE COMPLETED SUCCESSFULLY  ║")
	log.Println("╚══════════════════════════════════════╝")
	return nil
}

func prepare() error {
	log.Println("[1/5] Preparing environment...")
	if err := os.RemoveAll(config.DownloadDir); err != nil {
		return fmt.Errorf("cleanup download dir: %w", err)
	}
	if err := os.Mkdir(config.DownloadDir, 0755); err != nil {
		return fmt.Errorf("mkdir download dir: %w", err)
	}

	oxideURL, err := downloader.GetGitHubReleaseAssetURL(config.OxideRepoAPI, config.OxideAssetName)
	if err != nil {
		return fmt.Errorf("fetch oxide url: %w", err)
	}

	resources := map[string]string{
		"oxide.zip":    oxideURL,
		"steamcmd.zip": config.SteamCMDURL,
	}

	log.Println("[2/5] Downloading resources...")
	for name, url := range resources {
		path := filepath.Join(config.DownloadDir, name)
		log.Printf("   ◉ Downloading: %s", name)
		if err := downloader.DownloadFile(url, path); err != nil {
			return fmt.Errorf("download %s: %w", name, err)
		}
	}

	log.Println("[2/5] Extracting resources...")
	var errs []error
	for name := range resources {
		path := filepath.Join(config.DownloadDir, name)
		log.Printf("   ◉ Extracting: %s", name)
		if err := extractor.UnzipFile(path, config.DownloadDir); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}

	// Clean up downloaded zip files
	for name := range resources {
		path := filepath.Join(config.DownloadDir, name)
		_ = os.Remove(path)
	}

	return errors.Join(errs...)
}

func stopRustServer() error {
	log.Println("[3/5] Stopping Rust server...")
	return process.FindAndKillRustServer()
}
