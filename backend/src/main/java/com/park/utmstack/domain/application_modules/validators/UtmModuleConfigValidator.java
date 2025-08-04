package com.park.utmstack.domain.application_modules.validators;

import com.park.utmstack.config.Constants;
import com.park.utmstack.domain.application_modules.UtmModule;
import com.park.utmstack.domain.application_modules.UtmModuleGroupConfiguration;
import com.park.utmstack.repository.UtmModuleGroupConfigurationRepository;
import com.park.utmstack.service.application_modules.connectors.UtmStackConnectionService;
import com.park.utmstack.service.dto.application_modules.UtmModuleGroupConfDTO;
import com.park.utmstack.service.dto.application_modules.UtmModuleGroupConfWrapperDTO;
import com.park.utmstack.util.CipherUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class UtmModuleConfigValidator {

    private final UtmModuleGroupConfigurationRepository moduleGroupConfigurationRepository;
    private final UtmStackConnectionService utmStackConnectionService;

    /**
     * Validates if the given configuration allows a successful connection to UTMStack.
     *
     * @param keys A list of configuration keys that should include at least `ipAddress` and `connectionKey`.
     * @return true if login and ping are successful; false otherwise.
     * @throws Exception If required fields are missing or the connection fails.
     */
    public boolean validate(UtmModule module, List<UtmModuleGroupConfiguration> keys) throws Exception {
        if (keys.isEmpty()) return false;

        List<UtmModuleGroupConfiguration> configurations = moduleGroupConfigurationRepository
                .findAllByGroupId(keys.get(0).getGroupId())
                .stream()
                .map(c -> {
                    if (this.containsConfigInKeys(keys, c.getConfKey()) != null) {
                        return this.containsConfigInKeys(keys, c.getConfKey());
                    } else {
                        if ("password".equals(c.getConfDataType())) {
                            c.setConfValue(CipherUtil.decrypt(
                                    c.getConfValue(),
                                    System.getenv(Constants.ENV_ENCRYPTION_KEY)
                            ));
                        }
                        return c;
                    }
                })
                .collect(Collectors.toList());

        List<UtmModuleGroupConfDTO> configDTOs = configurations.stream()
                .map(entity -> new UtmModuleGroupConfDTO(entity.getConfKey(), entity.getConfValue()))
                .collect(Collectors.toList());

        UtmModuleGroupConfWrapperDTO body = new UtmModuleGroupConfWrapperDTO(configDTOs);

        return utmStackConnectionService.testConnection(module.getModuleName().name(),body);
    }

    private UtmModuleGroupConfiguration containsConfigInKeys(List<UtmModuleGroupConfiguration> keys, String confKey) {
        return keys.stream()
                .filter(key -> key.getConfKey().equals(confKey))
                .findFirst()
                .orElse(null);
    }
}
