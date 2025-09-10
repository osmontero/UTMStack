import {Component, OnInit} from '@angular/core';
import {DomSanitizer} from '@angular/platform-browser';
import {Router} from '@angular/router';
import {NgxSpinnerService} from 'ngx-spinner';
import {Observable} from 'rxjs';
import {AuthServerProvider} from '../../../../core/auth/auth-jwt.service';
import {ThemeChangeBehavior} from '../../../behaviors/theme-change.behavior';
import {TfaMethod} from '../../../services/tfa/tfa.service';
import {extractQueryParamsForNavigation} from "../../../util/query-params-to-filter.util";
import {ADMIN_DEFAULT_EMAIL, ADMIN_ROLE} from "../../../constants/global.constant";
import {StateStorageService} from "../../../../core/auth/state-storage.service";
import {AccountService} from "../../../../core/auth/account.service";
import {UtmToastService} from "../../../alert/utm-toast.service";


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
              private themeChangeBehavior: ThemeChangeBehavior,
              public sanitizer: DomSanitizer,
              private stateStorageService: StateStorageService,
              private accountService: AccountService,
              private utmToast: UtmToastService) {
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
        this.startNavigation();
      }
    }, error => {
      this.errorMessage = error.headers.get('X-UtmStack-error');
      this.isVerifying = false;
    });
  }

  backToLogin() {
    this.router.navigate(['/']);
  }

  onExpire() {

  }

  clearError() {
    if (this.verificationCode.length === 6) {
      this.onSubmit();
    }
    this.errorMessage = '';
  }

  startNavigation() {
    this.accountService.identity(true).then(account => {
      if (account) {
        const { path, queryParams } =
          extractQueryParamsForNavigation(this.stateStorageService.getUrl() ? this.stateStorageService.getUrl() : '' );
        if (path) {
          this.stateStorageService.resetPreviousUrl();
        }
        const redirectTo = (account.authorities.includes(ADMIN_ROLE) && account.email === ADMIN_DEFAULT_EMAIL)
          ? '/getting-started' : !!path ? path : '/dashboard/overview';
        console.log(redirectTo);
        this.router.navigate([redirectTo], {queryParams})
          .then(() => this.spinner.hide());
      } else {
        this.utmToast.showError('Login error', 'User without privileges.');
      }
    });
  }
}
