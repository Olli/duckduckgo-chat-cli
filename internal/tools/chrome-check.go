package tools

import (
	"fmt"
	"os/exec"
	"os"
	"runtime"
	"strconv"
	"strings"
	"github.com/fatih/color"

)
const (
	MinChromeVersion = "115.0.5790.110"

)

func CheckChromeVersion() {
	version, err := getChromeVersion()
	if err != nil {
		color.Yellow("Warning: %v", err)
		color.Yellow("Continuing without Chrome version check...")
		return // Continue instead of panic
	}

	result, err := compareVersions(version, MinChromeVersion)
	if err != nil {
		color.Yellow("Warning: Chrome version check failed: %v", err)
		color.Yellow("Continuing without version verification...")
		return // Continue instead of panic
	}

	if result < 0 {
		color.Yellow("Warning: Chrome %s+ recommended, found %s", MinChromeVersion, version)
		color.Yellow("The application might still work, but it's recommended to upgrade Chrome")
	}
}

func compareVersions(v1, v2 string) (int, error) {
	// Clean version strings
	v1 = strings.TrimSpace(v1)
	v1 = strings.TrimPrefix(v1, "Google Chrome ")
	v1 = strings.TrimPrefix(v1, "Chromium ")

	// Remove "snap" suffix if present
	if idx := strings.Index(v1, " snap"); idx != -1 {
		v1 = v1[:idx]
	}

	v1parts := strings.Split(v1, ".")
	v2parts := strings.Split(v2, ".")

	if len(v1parts) == 0 || len(v2parts) == 0 {
		return 0, fmt.Errorf("invalid version format")
	}

	// Compare major version numbers
	v1num, err := strconv.Atoi(strings.TrimSpace(v1parts[0]))
	if err != nil {
		return 0, fmt.Errorf("invalid version: %s", v1)
	}

	v2num, err := strconv.Atoi(strings.TrimSpace(v2parts[0]))
	if err != nil {
		return 0, fmt.Errorf("invalid version: %s", v2)
	}

	if v1num > v2num {
		return 1, nil
	} else if v1num < v2num {
		return -1, nil
	}
	return 0, nil
}

func getChromeVersion() (string, error) {
	switch runtime.GOOS {
	case "windows":
		paths := []string{
			"reg query \"HKEY_CURRENT_USER\\Software\\Google\\Chrome\\BLBeacon\" /v version",
			"reg query \"HKLM\\SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Google Chrome\" /v Version",
		}

		for _, cmd := range paths {
			if output, err := exec.Command("cmd", "/C", cmd).Output(); err == nil {
				if version := extractWindowsVersion(string(output)); version != "" {
					return version, nil
				}
			}
		}
		return "", fmt.Errorf("chrome not found in Windows registry")

	case "linux":
		browsers := []string{
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
		}

		for _, browser := range browsers {
			if output, err := exec.Command(browser, "--version").Output(); err == nil {
				return strings.TrimSpace(string(output)), nil
			}
		}

		// Vérifier les chemins communs
		paths := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/snap/bin/chromium",
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				if output, err := exec.Command(path, "--version").Output(); err == nil {
					return strings.TrimSpace(string(output)), nil
				}
			}
		}

		return "", fmt.Errorf("neither Chrome nor Chromium found on system")

	case "darwin":
		paths := []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}

		for _, path := range paths {
			if output, err := exec.Command(path, "--version").Output(); err == nil {
				return strings.TrimSpace(string(output)), nil
			}
		}
		return "", fmt.Errorf("chrome not found on macOS")

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func extractWindowsVersion(output string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "REG_SZ") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[len(fields)-1]
			}
		}
	}
	return ""
}



