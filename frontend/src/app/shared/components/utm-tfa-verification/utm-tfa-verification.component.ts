import {Component, Input} from '@angular/core';
import {TfaMethod} from '../../services/tfa/tfa.service';

@Component({
  selector: 'app-utm-tfa-verification',
  templateUrl: './utm-tfa-verification.component.html',
  styleUrls: ['./utm-tfa-verification.component.css']
})
export class UtmTfaVerificationComponent {

  @Input() method: TfaMethod;
  @Input() qrCodeUrl?: string;
  @Input() expiresInSeconds = 300;

  code = '';
  verifying = false;
  isLoginFailed = false;
  errorMessage = '';
  protected readonly TfaMethod = TfaMethod;

  onSubmit() {
    this.verifying = true;
    setTimeout(() => {
      this.verifying = false;
      this.isLoginFailed = this.code !== '123456';
      this.errorMessage = this.isLoginFailed ? 'Invalid code' : '';
    }, 1000);
  }

  onExpire() {
    this.errorMessage = 'Verification expired. Please try again.';
  }
}
