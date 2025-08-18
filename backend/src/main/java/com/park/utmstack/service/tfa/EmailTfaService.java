package com.park.utmstack.service.tfa;

import com.park.utmstack.domain.User;
import com.park.utmstack.domain.tfa.TfaMethod;
import com.park.utmstack.domain.tfa.TfaSetupState;
import com.park.utmstack.service.MailService;
import com.park.utmstack.service.dto.tfa.init.Delivery;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyResponse;
import com.park.utmstack.util.exceptions.UtmMailException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.concurrent.TimeUnit;

@Service
@RequiredArgsConstructor
public class EmailTfaService implements TfaMethodService {

    private static final String CLASSNAME = "EmailTfaService";
    private final CacheService cache;
    private final ConfigService configService;
    private final EmailTotpService tfaService;
    private final MailService mailService;

    private static final int CODE_LENGTH = 6;
    private static final long EXPIRES_IN_SECONDS = 300;

    @Override
    public TfaMethod getMethod() {
        return TfaMethod.EMAIL;
    }

    @Override
    public TfaInitResponse initiateSetup(User user) {
        final String ctx = CLASSNAME + ".initiateSetup";
        try {
            mailService.sendCheckEmail(List.of(user.getEmail()));

            String secret = tfaService.generateSecret();
            String code = tfaService.generateCode(secret);

            long expiresAt = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(300);
            TfaSetupState state = new TfaSetupState(secret, expiresAt);
            cache.storeState(user.getLogin(), TfaMethod.TOTP, state);

            mailService.sendTfaVerificationCode(user, code);

            Delivery delivery = new Delivery(TfaMethod.EMAIL, "Código enviado por correo.");
            return new TfaInitResponse("pending", delivery, 300);
        }
        catch (Exception e) {
            throw new UtmMailException(ctx + ": " + e.getLocalizedMessage());
        }
    }


    @Override
    public TfaVerifyResponse verifyCode(User user, String code) {

        TfaSetupState tfaSetupState = cache.getState(user.getLogin(), TfaMethod.EMAIL)
                .orElseThrow(() -> new IllegalStateException("No TFA setup found for user: " + user.getLogin()));

        boolean expired = tfaSetupState.isExpired();
        boolean valid = !expired && tfaService.validateCode(tfaSetupState.getSecret(), code);

        return new TfaVerifyResponse(
                valid,
                expired,
                tfaSetupState.getRemainingSeconds(),
                expired ? "Setup expired" : "Code verification " + (valid ? "successful" : "failed")
        );
    }

    @Override
    public void persistConfiguration(User user) {
        String secret = cache.getState(user.getLogin(), TfaMethod.EMAIL)
                .orElseThrow(() -> new IllegalStateException("No TFA setup found for user: " + user.getLogin()))
                .getSecret();
        configService.enableTfa(user.getLogin(), TfaMethod.EMAIL, secret);
        cache.clear(user.getLogin(), TfaMethod.EMAIL);
    }
}

