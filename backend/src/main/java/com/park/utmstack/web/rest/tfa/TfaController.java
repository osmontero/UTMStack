package com.park.utmstack.web.rest.tfa;

import com.park.utmstack.service.dto.tfa.init.TfaInitRequest;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
@RequestMapping("/tfa")
public class TfaController {


    @PostMapping("/init")
    public ResponseEntity<TfaInitResponse> initTfa(@RequestBody TfaInitRequest request) {

        TfaInitResponse response = tfaService.initiateSetup(username, request.getMethod());
        return ResponseEntity.ok(response);
    }
}
