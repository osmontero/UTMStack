package com.park.utmstack.service.tfa;

import com.google.zxing.BarcodeFormat;
import com.google.zxing.MultiFormatWriter;
import com.google.zxing.client.j2se.MatrixToImageWriter;
import com.google.zxing.common.BitMatrix;
import com.park.utmstack.domain.User;
import com.park.utmstack.domain.tfa.TfaMethod;
import com.park.utmstack.domain.tfa.TfaSetupState;
import com.park.utmstack.service.dto.tfa.init.Delivery;
import com.park.utmstack.service.dto.tfa.init.TfaInitResponse;
import com.park.utmstack.service.dto.tfa.verify.TfaVerifyResponse;
import com.warrenstrange.googleauth.GoogleAuthenticator;
import org.springframework.stereotype.Service;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.ByteArrayOutputStream;
import java.util.Base64;
import java.util.concurrent.TimeUnit;

import static com.park.utmstack.config.Constants.TFA_ISSUER;

@Service
public class TotpTfaService implements TfaMethodService {

    private final GoogleAuthenticator authenticator;
    private final CacheService cache;
    private final ConfigService configService;

    TotpTfaService(CacheService cache, ConfigService configService) {
        this.authenticator = new GoogleAuthenticator();
        this.cache = cache;
        this.configService = configService;
    }

    @Override
    public TfaMethod getMethod() {
        return TfaMethod.TOTP;
    }

    @Override
    public TfaInitResponse initiateSetup(User user) {
        String secret = authenticator.createCredentials().getKey();
        long expiresAt = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(300);
        TfaSetupState state = new TfaSetupState(secret, expiresAt);
        cache.storeState(user.getLogin(), TfaMethod.TOTP, state);

        String uri = String.format("otpauth://totp/%s:%s?secret=%s&issuer=%s",
                TFA_ISSUER, user.getLogin(), secret, TFA_ISSUER);

        String qrBase64 = generateQrBase64(uri);
        Delivery delivery = new Delivery(TfaMethod.TOTP, qrBase64);

        long expiresInSeconds = 300;

        return new TfaInitResponse("pending", delivery, expiresInSeconds);
    }

    @Override
    public TfaVerifyResponse verifyCode(User user, String code) {
        TfaSetupState tfaSetupState = cache.getState(user.getLogin(), TfaMethod.TOTP)
                .orElseThrow(() -> new IllegalStateException("No TFA setup found for user: " + user.getLogin()));

        boolean expired = tfaSetupState.isExpired();
        boolean valid = !expired && authenticator.authorize(tfaSetupState.getSecret(), Integer.parseInt(code));

        return new TfaVerifyResponse(
                valid,
                expired,
                tfaSetupState.getRemainingSeconds(),
                expired ? "Setup expired" : "Code verification " + (valid ? "successful" : "failed")
        );
    }


    @Override
    public void persistConfiguration(User user) {
        String secret = cache.getState(user.getLogin(), TfaMethod.TOTP)
                .orElseThrow(() -> new IllegalStateException("No TFA setup found for user: " + user.getLogin()))
                .getSecret();
        configService.enableTfa(user.getLogin(), TfaMethod.TOTP, secret);
        cache.clear(user.getLogin(), TfaMethod.TOTP);
    }

    private String generateQrBase64(String uri) {
        try {
            BitMatrix matrix = new MultiFormatWriter().encode(uri, BarcodeFormat.QR_CODE, 200, 200);
            BufferedImage image = MatrixToImageWriter.toBufferedImage(matrix);

            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            ImageIO.write(image, "png", baos);
            return Base64.getEncoder().encodeToString(baos.toByteArray());
        } catch (Exception e) {
            throw new RuntimeException("Error al generar QR", e);
        }
    }
}


