package com.park.utmstack.security.tfa;

public interface TfaService {
    boolean sendChallenge(String userId);
    boolean verifyResponse(String userId, String response);
    String getMethod();
}

