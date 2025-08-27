// tfa-setup.component.ts
import { Component, OnDestroy, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { interval, Observable, Subscription } from 'rxjs';
import {ThemeChangeBehavior} from '../../../behaviors/theme-change.behavior';
import {TfaMethod} from '../../../services/tfa/tfa.service';

@Component({
  selector: 'app-tfa-setup',
  templateUrl: './tfa-setup.component.html',
  styleUrls: ['./tfa-setup.component.scss']
})
export class TfaSetupComponent implements OnInit, OnDestroy {

  TfaMethod = TfaMethod;

  step: 'method-selection' | 'setup' | 'verification' | 'success' = 'method-selection';
  selectedMethod: TfaMethod | null = null;

  qrCodeUrl = '';

  code  = '';
  verifying = false;
  codeVerified = false;
  errorMessage = '';

  expiresInSeconds = 300;
  private timerSubscription: Subscription | null = null;

  resending = false;

  loginImage$: Observable<string> | undefined;

  constructor( private themeChangeBehavior: ThemeChangeBehavior,
               public sanitizer: DomSanitizer,
  ) {}

  ngOnInit(): void {
    this.loginImage$ = this.themeChangeBehavior.$themeIcon.asObservable();
  }

  ngOnDestroy(): void {
    if (this.timerSubscription) {
      this.timerSubscription.unsubscribe();
    }
  }


  selectMethod(method: TfaMethod): void {
    this.selectedMethod = method;
    this.step = 'setup';
    this.clearError();

    if (method === TfaMethod.TOTP) {
      this.generateTOTPData();
    } else if (method === TfaMethod.EMAIL) {
      this.sendEmailCode();
    }

    this.startTimer();
  }

  private generateTOTPData(): void {
    // Llamar al API para generar los datos TOTP
    // this.tfaService.generateTOTPSetup().subscribe(response => {
    //   this.qrCodeUrl = response.qrCodeUrl;
    // });

    // Mock para desarrollo - en producción remover esto
    const userEmail = 'usuario@ejemplo.com';
    const appName = 'UTMStack';
    const mockSecret = 'JBSWY3DPEHPK3PXP';
    const otpauthUrl = `otpauth://totp/${appName}:${userEmail}?secret=${mockSecret}&issuer=${appName}`;
    this.qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(otpauthUrl)}`;
  }

  private sendEmailCode(): void {
    this.resending = true;

    // En producción, llamada al backend
    setTimeout(() => {
      this.resending = false;
      console.log('Código enviado por email');
      // this.tfaService.sendEmailCode().subscribe(() => {
      //   this.resending = false;
      // });
    }, 1000);
  }


  resendCode(): void {
    if (this.selectedMethod === TfaMethod.EMAIL) {
      this.sendEmailCode();
      this.resetTimer();
      this.clearError();
    }
  }

  private startTimer(): void {
    this.expiresInSeconds = 300; // 5 minutos
    this.timerSubscription = interval(1000).subscribe(() => {
      if (this.expiresInSeconds > 0) {
        this.expiresInSeconds--;
      } else {
        this.onExpire();
      }
    });
  }

  /**
   * Reinicia el timer
   */
  private resetTimer(): void {
    if (this.timerSubscription) {
      this.timerSubscription.unsubscribe();
    }
    this.startTimer();
  }

  /**
   * Maneja la expiración del código
   */
  onExpire(): void {
    this.errorMessage = 'The code has expired. Please request a new one.';
    if (this.timerSubscription) {
      this.timerSubscription.unsubscribe();
    }

    if (this.selectedMethod === TfaMethod.EMAIL) {
      // Para email, automáticamente reenviar
      this.sendEmailCode();
      this.resetTimer();
    }
  }

  /**
   * Submits el formulario de verificación
   */
  onSubmit(): void {
    if (this.code.length < 6) {
      this.errorMessage = 'The code must be 6 digits long';
      return;
    }

    this.verifying = true;
    this.clearError();

    // Simular verificación (en producción sería llamada a API)
    setTimeout(() => {
      // Mock verification - en producción:
      // this.tfaService.verifySetupCode(this.code, this.selectedMethod, this.secretKey)
      //   .subscribe({
      //     next: (response) => {
      //       if (response.valid) {
      //         this.codeVerified = true;
      //         this.step = 'success';
      //       } else {
      //         this.errorMessage = 'Invalid code. Please try again.';
      //       }
      //       this.verifying = false;
      //     },
      //     error: (error) => {
      //       this.errorMessage = 'Verification failed. Please try again.';
      //       this.verifying = false;
      //     }
      //   });

      // Mock verification
      if (this.code === '123456' || this.code.length === 6) {
        this.codeVerified = true;
        this.step = 'success';
        if (this.timerSubscription) {
          this.timerSubscription.unsubscribe();
        }
      } else {
        this.errorMessage = 'Invalid code. Please try again.';
      }
      this.verifying = false;
    }, 1500);
  }

  /**
   * Limpia el mensaje de error
   */
  clearError(): void {
    this.errorMessage = '';
  }

  /**
   * Regresa a la selección de método
   */
  backToMethodSelection(): void {
    this.step = 'method-selection';
    this.selectedMethod = null;
    this.clearError();
    this.code = '';
    if (this.timerSubscription) {
      this.timerSubscription.unsubscribe();
    }
  }

  /**
   * Completa la configuración
   */
  completeSetup(): void {
    // En producción, guardar la configuración final
    console.log('Configuración TFA completada:', {
      method: this.selectedMethod
    });

    // this.tfaService.completeTfaSetup({
    //   method: this.selectedMethod
    // }).subscribe(() => {
    //   // Redirigir o cerrar modal
    //   this.router.navigate(['/dashboard']);
    // });

    alert('2FA setup completed successfully!');
  }

  /**
   * Omite la configuración TFA
   */
  skipSetup(): void {
    const confirmed = confirm(
      'Are you sure you want to skip Two-Factor Authentication setup?\n\n' +
      'This will make your account less secure. You can always set it up later in your security settings.'
    );

    if (confirmed) {
      console.log('User skipped TFA setup');
      // En producción, marcar en el backend que el usuario omitió la configuración
      // this.tfaService.skipTfaSetup().subscribe(() => {
      //   this.router.navigate(['/dashboard']);
      // });

      // Por ahora, simular redirección
      alert('TFA setup skipped. You can configure it later in Settings > Security.');
    }
  }
}
