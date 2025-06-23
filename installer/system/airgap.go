package system

import (
	"fmt"
	"strings"

	"github.com/utmstack/UTMStack/installer/config"
	"github.com/utmstack/UTMStack/installer/utils"
)

func VerifyAirGapPrerequisites() error {
	var errors []string

	if err := utils.RunCmd("docker", "--version"); err != nil {
		errors = append(errors, "Docker is not installed")
	} else {
		if err := utils.RunCmd("docker", "info"); err != nil {
			errors = append(errors, "Docker is installed but not running")
		}
	}

	if err := utils.RunCmd("nginx", "-v"); err != nil {
		errors = append(errors, "Nginx is not installed (required for reverse proxy)")
	}

	if !utils.CheckIfPathExist(config.ImagesPath) {
		errors = append(errors, fmt.Sprintf("Docker images directory not found: %s", config.ImagesPath))
	}

	if err := utils.RunCmd("systemctl", "--version"); err != nil {
		errors = append(errors, "systemctl is not available (required for service management)")
	}

	if len(errors) > 0 {
		return fmt.Errorf("AirGap prerequisites not met:\n  - %s\n\nPlease install missing components before running AirGap installation", strings.Join(errors, "\n  - "))
	}

	return nil
}
