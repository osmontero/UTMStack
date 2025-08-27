package com.park.utmstack.domain.tfa;

import lombok.Getter;
import lombok.ToString;

@Getter
@ToString
public class TfaSetupState {
    private final String secret;
    private final long expiresAt;
    private final long setupStartedAt;

    public TfaSetupState(String secret, long expiresAt) {
        this.secret = secret;
        this.expiresAt = expiresAt;
        this.setupStartedAt = System.currentTimeMillis();
    }

    public boolean isExpired() {
        return System.currentTimeMillis() > expiresAt;
    }

    public long getRemainingSeconds() {
        long remaining = expiresAt - System.currentTimeMillis();
        return Math.max(remaining / 1000, 0);
    }
}


