import {Component, OnInit} from '@angular/core';
import {DomSanitizer} from '@angular/platform-browser';
import {Router} from '@angular/router';
import {NgxSpinnerService} from 'ngx-spinner';
import {Observable} from "rxjs";
import {AuthServerProvider} from '../../../../core/auth/auth-jwt.service';
import {UtmToastService} from '../../../alert/utm-toast.service';
import {ThemeChangeBehavior} from '../../../behaviors/theme-change.behavior';
import {TfaMethod} from '../../../services/tfa/tfa.service';


@Component({
  selector: 'app-totp',
  templateUrl: './totp.component.html',
  styleUrls: ['./totp.component.scss']
})
export class TotpComponent implements OnInit {
  form: any = {};
  isLoggedIn = false;
  isLoginFailed = false;
  errorMessage = '';
  currentUser: any;
  isVerifying = false;
  loginImage$: Observable<string>;
  TfaMethod = TfaMethod;
  method: TfaMethod;
  isVerified = false;
  verificationCode = '';

  constructor(private authService: AuthServerProvider,
              private router: Router,
              private spinner: NgxSpinnerService,
              private utmToast: UtmToastService,
              private themeChangeBehavior: ThemeChangeBehavior,
              public sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
    this.method = this.authService.tfaMethod;
    this.loginImage$ = this.themeChangeBehavior.$themeIcon.asObservable();
  }


  onSubmit() {
    this.isVerifying = true;
    this.authService
      .verifyCode(this.verificationCode).subscribe((auth) => {
      if (auth) {
        this.isVerified = true;
        this.spinner.show();
        this.router.navigate(['/dashboard/overview'])
          .then(() => this.spinner.hide());
      }
    }, error => {
      this.errorMessage = error.headers.get('X-UtmStack-error');
      console.log(error.headers.get('X-UtmStack-error'));
      this.isVerifying = false;
    });
  }

  backToLogin() {
    this.router.navigate(['/']);
  }

  onExpire() {

  }

  clearError() {
    this.errorMessage = '';
  }
}
