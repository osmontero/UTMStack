package com.park.utmstack.service.dto.application_modules;

import com.park.utmstack.domain.application_modules.UtmModuleGroupConfiguration;
import com.park.utmstack.domain.application_modules.validators.ValidModuleConfiguration;
import lombok.Data;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
@ValidModuleConfiguration
public class GroupConfigurationDTO {
    @NotNull
    private Long moduleId;
    @NotEmpty
    private List<UtmModuleGroupConfiguration> keys;
}
