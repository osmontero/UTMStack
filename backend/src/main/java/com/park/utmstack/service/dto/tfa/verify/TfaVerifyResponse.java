package com.park.utmstack.service.dto.tfa.verify;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class TfaVerifyResponse {
    private Boolean status;
    private String message;
}
