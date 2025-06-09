package com.park.utmstack.service.alert_response_rule;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import com.jayway.jsonpath.Criteria;
import com.jayway.jsonpath.Filter;
import com.jayway.jsonpath.Predicate;
import com.park.utmstack.config.Constants;
import com.park.utmstack.domain.alert_response_rule.UtmAlertResponseActionTemplate;
import com.park.utmstack.domain.alert_response_rule.UtmAlertResponseRule;
import com.park.utmstack.domain.alert_response_rule.UtmAlertResponseRuleExecution;
import com.park.utmstack.domain.alert_response_rule.UtmAlertResponseRuleHistory;
import com.park.utmstack.domain.alert_response_rule.enums.RuleExecutionStatus;
import com.park.utmstack.domain.alert_response_rule.enums.RuleNonExecutionCause;
import com.park.utmstack.domain.application_events.enums.ApplicationEventType;
import com.park.utmstack.domain.chart_builder.types.query.FilterType;
import com.park.utmstack.domain.chart_builder.types.query.OperatorType;
import com.park.utmstack.domain.compliance.UtmComplianceStandardSection;
import com.park.utmstack.domain.shared_types.AlertType;
import com.park.utmstack.repository.alert_response_rule.UtmAlertResponseActionTemplateRepository;
import com.park.utmstack.repository.alert_response_rule.UtmAlertResponseRuleExecutionRepository;
import com.park.utmstack.repository.alert_response_rule.UtmAlertResponseRuleHistoryRepository;
import com.park.utmstack.repository.alert_response_rule.UtmAlertResponseRuleRepository;
import com.park.utmstack.repository.network_scan.UtmNetworkScanRepository;
import com.park.utmstack.service.UtmStackService;
import com.park.utmstack.service.agent_manager.AgentService;
import com.park.utmstack.service.application_events.ApplicationEventService;
import com.park.utmstack.service.dto.UtmAlertResponseActionTemplateDTO;
import com.park.utmstack.service.dto.UtmAlertResponseRuleDTO;
import com.park.utmstack.service.dto.agent_manager.AgentDTO;
import com.park.utmstack.service.dto.agent_manager.AgentStatusEnum;
import com.park.utmstack.service.grpc.CommandResult;
import com.park.utmstack.service.incident_response.UtmIncidentVariableService;
import com.park.utmstack.service.incident_response.grpc_impl.IncidentResponseCommandService;
import com.park.utmstack.util.UtilJson;
import com.park.utmstack.util.exceptions.UtmNotImplementedException;
import io.grpc.stub.StreamObserver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Async;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StringUtils;

import java.util.*;
import java.util.concurrent.TimeUnit;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

@Service
@Transactional
public class UtmAlertResponseActionTemplateService {

    private static final String CLASSNAME = "UtmAlertResponseActionTemplateService";
    private final Logger log = LoggerFactory.getLogger(UtmAlertResponseActionTemplateService.class);

    private final UtmAlertResponseActionTemplateRepository utmAlertResponseActionTemplateRepository;

    private final UtmStackService utmStackService;


    public UtmAlertResponseActionTemplateService(UtmAlertResponseActionTemplateRepository utmAlertResponseActionTemplateRepository,
                                                 UtmStackService utmStackService) {
        this.utmAlertResponseActionTemplateRepository = utmAlertResponseActionTemplateRepository;
        this.utmStackService = utmStackService;
    }

    public UtmAlertResponseActionTemplate save(UtmAlertResponseActionTemplate alertResponseActionTemplate) {
        final String ctx = CLASSNAME + ".save";
        try {
            if (!utmStackService.isInDevelop()) {
                alertResponseActionTemplate.setId(this.getSystemSequenceNextValue());
                alertResponseActionTemplate.setSystemOwner(true);
            } else {
                alertResponseActionTemplate.setSystemOwner(false);
            }

            return utmAlertResponseActionTemplateRepository.save(alertResponseActionTemplate);
        } catch (Exception e) {
            throw new RuntimeException(ctx + ": " + e.getLocalizedMessage());
        }
    }


    public Long getSystemSequenceNextValue() {
        return utmAlertResponseActionTemplateRepository.findFirstBySystemOwnerIsTrueOrderByIdDesc()
                .map(rule -> rule.getId() + 1)
                .orElse(1L);
    }

}
