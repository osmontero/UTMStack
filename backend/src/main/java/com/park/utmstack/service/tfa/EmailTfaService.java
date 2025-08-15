package com.park.utmstack.service.tfa;

import com.park.utmstack.domain.User;
import com.park.utmstack.domain.tfa.TfaMethod;
import com.park.utmstack.service.MailService;
import com.park.utmstack.service.dto.tfa.init.Delivery;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.time.Duration;

@Service
@RequiredArgsConstructor
public class EmailTfaService implements TfaMethodService {

    private final CacheService cache;
    private final ConfigService configService;
    private TfaService tfaService;
    private MailService mailService;

    private static final int CODE_LENGTH = 6;
    private static final long EXPIRES_IN_SECONDS = 300;

    @Override
    public TfaMethod getMethod() {
        return TfaMethod.EMAIL;
    }

    @Override
    public TfaInitResponse initiateSetup(User user) {
        String secret = tfaService.generateSecret();
        String code = tfaService.generateCode(secret);

        cache.storeSecret(user.getLogin(), secret);
        mailService.sendTfaVerificationCode(user, code);

        Delivery delivery = new Delivery(TfaMethod.EMAIL, "Código enviado por correo.");
        return new TfaInitResponse("pending", delivery, 300);
    }


    @Override
    public TfaVerifyResponse verifyCode(User user, String code) {
        boolean valid = tfaService.validateCode(cache.getSecret(user.getLogin()), code);
        return new TfaVerifyResponse(valid, valid ? "Código verificado." : "Código inválido.");
    }

    @Override
    public void persistConfiguration(User user) {
        configService.enableTfa(user.getLogin(), TfaMethod.EMAIL, cache.getSecret(user.getLogin()));
        cache.clearSecret(user.getLogin());
    }
}

