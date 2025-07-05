package validations

import (
	"fmt"

	"github.com/utmstack/UTMStack/plugins/modules-config/config"
)

func ValidateModuleConfig(moduleName string, config *config.ModuleGroup) error {
	switch moduleName {
	case "SOPHOS":
		if err := ValidateSophosConfig(config); err != nil {
			return fmt.Errorf("%v", err)
		}
	default:
		return fmt.Errorf("unsupported module: %s", moduleName)
	}

	return nil
}
