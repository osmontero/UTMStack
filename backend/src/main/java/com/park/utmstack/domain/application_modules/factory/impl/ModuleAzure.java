package com.park.utmstack.domain.application_modules.factory.impl;

import com.park.utmstack.domain.application_modules.UtmModule;
import com.park.utmstack.domain.application_modules.enums.ModuleName;
import com.park.utmstack.domain.application_modules.factory.IModule;
import com.park.utmstack.domain.application_modules.types.ModuleConfigurationKey;
import com.park.utmstack.domain.application_modules.types.ModuleRequirement;
import com.park.utmstack.service.application_modules.UtmModuleService;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

@Component
public class ModuleAzure implements IModule {
    private static final String CLASSNAME = "ModuleAzure";

    private final UtmModuleService moduleService;

    public ModuleAzure(UtmModuleService moduleService) {
        this.moduleService = moduleService;
    }

    @Override
    public UtmModule getDetails(Long serverId) throws Exception {
        final String ctx = CLASSNAME + ".getDetails";
        try {
            return moduleService.findByServerIdAndModuleName(serverId, ModuleName.AZURE);
        } catch (Exception e) {
            throw new Exception(ctx + ": " + e.getMessage());
        }
    }

    @Override
    public List<ModuleRequirement> checkRequirements(Long serverId) throws Exception {
        return Collections.emptyList();
    }

    @Override
    public List<ModuleConfigurationKey> getConfigurationKeys(Long groupId) throws Exception {
        List<ModuleConfigurationKey> keys = new ArrayList<>();

        // workspaceId
        keys.add(ModuleConfigurationKey.builder()
                .withGroupId(groupId)
                .withConfKey("workspaceId")
                .withConfName("Workspace ID")
                .withConfDescription("Azure Log Analytics Workspace ID")
                .withConfDataType("text")
                .withConfRequired(true)
                .build());

        // clientId
        keys.add(ModuleConfigurationKey.builder()
                .withGroupId(groupId)
                .withConfKey("clientId")
                .withConfName("Application (client) ID")
                .withConfDescription("Azure AD Application (client) ID")
                .withConfDataType("text")
                .withConfRequired(true)
                .build());

        // tenantId
        keys.add(ModuleConfigurationKey.builder()
                .withGroupId(groupId)
                .withConfKey("tenantId")
                .withConfName("Directory (tenant) ID")
                .withConfDescription("Azure Active Directory (tenant) ID")
                .withConfDataType("text")
                .withConfRequired(true)
                .build());

        // clientSecret
        keys.add(ModuleConfigurationKey.builder()
                .withGroupId(groupId)
                .withConfKey("clientSecret")
                .withConfName("Client Secret Value")
                .withConfDescription("Azure AD Application Client Secret")
                .withConfDataType("text")
                .withConfRequired(true)
                .build());

        return keys;
    }
}
