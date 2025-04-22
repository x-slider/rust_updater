package config

import (
	"os"
	"path/filepath"
)

const (
	AppName     = "rust-updater"
	AppVersion  = "1.0.0"
	DownloadDir = "download"
)

var (
	SteamCMDPath = filepath.Join(DownloadDir, "steamcmd.exe")
	SteamCMDURL  = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
	RustAppID    = "258550"
)

var (
	OxideRepoAPI   = "https://api.github.com/repos/OxideMod/Oxide.Rust/releases/latest"
	OxideAssetName = "Oxide.Rust.zip"

	OxideSourceDir = func(workDir string) string {
		return filepath.Join(workDir, DownloadDir, "RustDedicated_Data")
	}

	OxideTargetDir = func(workDir string) string {
		return filepath.Join(workDir, "RustDedicated_Data")
	}
)

var ServerInstallPath string

func InitWorkDir() error {
	var err error
	ServerInstallPath, err = os.Getwd()
	return err
}

func GetWorkDir() (string, error) {
	if ServerInstallPath == "" {
		return os.Getwd()
	}
	return ServerInstallPath, nil
}
