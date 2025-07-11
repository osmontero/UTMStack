package com.park.utmstack.service.application_modules.connectors;

import com.park.utmstack.config.Constants;
import com.park.utmstack.service.dto.application_modules.UtmModuleGroupConfWrapperDTO;
import com.park.utmstack.service.web_clients.rest_template.RestTemplateService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.*;
import org.springframework.stereotype.Service;


@Service
@RequiredArgsConstructor
public class UtmStackConnectionService {

    private final Logger log = LoggerFactory.getLogger(UtmStackConnectionService.class);
    private final RestTemplateService restTemplateService;
    private static final String CLASSNAME = "UtmStackConnectionService";


    public boolean testConnection(String module, UtmModuleGroupConfWrapperDTO configurations) throws Exception {
        final String ctx = CLASSNAME + ".testConnection";
        HttpHeaders headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        headers.add("Accept", "*/*");
        headers.set(Constants.EVENT_PROCESSOR_INTERNAL_KEY_HEADER, System.getenv(Constants.ENV_INTERNAL_KEY));

        String baseUrl = "http://" + System.getenv(Constants.ENV_EVENT_PROCESSOR_HOST)  + ":" + System.getenv(Constants.ENV_EVENT_PROCESSOR_PORT);
        String endPoint = baseUrl + "/api/v1/modules-config/validate?nameShort=" + module;
        try{
            ResponseEntity<String> response = restTemplateService.post(
                    endPoint,
                    configurations,
                    String.class,
                    headers
            );
            return response.getStatusCode().is2xxSuccessful();
        } catch (Exception e) {
            String msg = ctx + ": " + e.getLocalizedMessage();
            log.error(msg);
            throw new RuntimeException(ctx + ": " + e.getMessage());
        }
    }
}

