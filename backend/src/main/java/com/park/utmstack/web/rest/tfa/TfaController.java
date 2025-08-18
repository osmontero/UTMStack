package com.park.utmstack.web.rest.tfa;

import com.park.utmstack.domain.User;
import com.park.utmstack.domain.application_events.enums.ApplicationEventType;
import com.park.utmstack.service.UserService;
import com.park.utmstack.service.application_events.ApplicationEventService;
import com.park.utmstack.service.dto.tfa.init.TfaInitRequest;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyRequest;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyResponse;
import com.park.utmstack.service.tfa.TfaService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
@RequestMapping("api/tfa")
public class TfaController {

    private final Logger log = LoggerFactory.getLogger(TfaController.class);
    private static final String CLASSNAME = "TfaController";

    private final TfaService tfaService;
    private final UserService userService;
    private final ApplicationEventService applicationEventService;

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

}
