package com.park.utmstack.service.tfa;

import org.springframework.stereotype.Service;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class CacheService {
    private final Map<String, String> secretCache = new ConcurrentHashMap<>();

    public void storeSecret(String username, String secret) {
        secretCache.put(username, secret);
    }

    public String getSecret(String username) {
        return secretCache.get(username);
    }

    public void clearSecret(String username) {
        secretCache.remove(username);
    }
}

