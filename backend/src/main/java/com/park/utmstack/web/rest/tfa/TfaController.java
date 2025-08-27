package com.park.utmstack.web.rest.tfa;

import com.park.utmstack.config.Constants;
import com.park.utmstack.domain.User;
import com.park.utmstack.domain.UtmConfigurationParameter;
import com.park.utmstack.domain.application_events.enums.ApplicationEventType;
import com.park.utmstack.service.UserService;
import com.park.utmstack.service.UtmConfigurationParameterService;
import com.park.utmstack.service.application_events.ApplicationEventService;
import com.park.utmstack.service.dto.tfa.init.TfaInitRequest;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import com.park.utmstack.service.dto.tfa.save.TfaSaveRequest;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyRequest;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyResponse;
import com.park.utmstack.service.tfa.TfaService;
import com.park.utmstack.util.UtilResponse;
import com.park.utmstack.util.exceptions.UtmMailException;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

import static com.park.utmstack.config.Constants.PROP_TFA_METHOD;

@RestController
@RequiredArgsConstructor
@RequestMapping("api/tfa")
public class TfaController {

    private final Logger log = LoggerFactory.getLogger(TfaController.class);
    private static final String CLASSNAME = "TfaController";

    private final TfaService tfaService;
    private final UserService userService;
    private final ApplicationEventService applicationEventService;
    private final UtmConfigurationParameterService utmConfigurationParameterService;

    @PostMapping("/init")
    public ResponseEntity<TfaInitResponse> initTfa(@RequestBody TfaInitRequest request) {
        final String ctx = CLASSNAME + ".initTfa";
        try {
            User user = userService.getCurrentUserLogin();
            TfaInitResponse response = tfaService.initiateSetup(user, request.getMethod());
            return ResponseEntity.ok(response);
        } catch (Exception e){
            String msg = ctx + ": " + e.getMessage();
            log.error(msg);
            applicationEventService.createEvent(msg, ApplicationEventType.ERROR);
            throw e;
        }

    }

    @PostMapping("/verify")
    public ResponseEntity<TfaVerifyResponse> verifyTfa(@RequestBody TfaVerifyRequest request) {
        final String ctx = CLASSNAME + ".verifyTfa";
        try {
            User user = userService.getCurrentUserLogin();
            TfaVerifyResponse response = tfaService.verifyCode(user, request);
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            String msg = ctx + ": " + e.getMessage();
            log.error(msg);
            applicationEventService.createEvent(msg, ApplicationEventType.ERROR);
            throw e;
        }
    }

    @PostMapping("/complete")
    public ResponseEntity<Void> completeTfa(@RequestBody TfaSaveRequest request) {
        final String ctx = CLASSNAME + ".completeTfa";
        try {

            List<UtmConfigurationParameter> tfaParams = utmConfigurationParameterService.getConfigParameterBySectionId(Constants.TFA_SETTING_ID);

            for (UtmConfigurationParameter param : tfaParams) {
                switch (param.getConfParamShort()) {
                    case PROP_TFA_METHOD:
                        param.setConfParamValue(String.valueOf(request.getMethod()));
                        break;
                    case Constants.PROP_TFA_ENABLE:
                        param.setConfParamValue("true");
                        break;
                }
            }

            utmConfigurationParameterService.saveAll(tfaParams);
            return ResponseEntity.ok().build();
        } catch (UtmMailException e) {
            String msg = ctx + ": " + e.getMessage();
            log.error(msg);
            applicationEventService.createEvent(msg, ApplicationEventType.ERROR);
            return UtilResponse.buildPreconditionFailedResponse(msg);
        } catch (IllegalArgumentException e) {
            String msg = ctx + ": " + e.getMessage();
            log.error(msg);
            applicationEventService.createEvent(msg, ApplicationEventType.ERROR);
            return UtilResponse.buildBadRequestResponse(msg);
        } catch (Exception e) {
            String msg = ctx + ": " + e.getMessage();
            log.error(msg);
            applicationEventService.createEvent(msg, ApplicationEventType.ERROR);
            return UtilResponse.buildInternalServerErrorResponse(msg);
        }
    }

}
