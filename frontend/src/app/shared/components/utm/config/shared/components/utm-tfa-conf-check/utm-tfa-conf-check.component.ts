import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {AccountService} from '../../../../../../../core/auth/account.service';
import {UtmToastService} from '../../../../../../alert/utm-toast.service';
import {UtmConfigEmailCheckService} from '../../../../../../services/config/utm-config-email-check.service';
import {TfaService} from '../../../../../../services/tfa/tfa.service';
import {SectionConfigParamType} from '../../../../../../types/configuration/section-config-param.type';

@Component({
  selector: 'app-utm-tfa-conf-check',
  templateUrl: './utm-tfa-conf-check.component.html',
  styleUrls: ['./utm-tfa-conf-check.component.css']
})
export class UtmTfaConfCheckComponent implements OnInit {
  @Input() validConfig: boolean;
  @Input() config: SectionConfigParamType[] = [];
  @Input() configToSave: SectionConfigParamType[] = [];
  @Output() isChecked = new EventEmitter<boolean>();
  checking: any;
  email: string;

  constructor(private utmConfigEmailCheckService: UtmConfigEmailCheckService,
              private accountService: AccountService,
              private utmToastService: UtmToastService,
              private tfaService: TfaService) {
  }

  ngOnInit() {
  }

  initTfa() {
    const tfaMethod = this.config.find(conf => conf.confParamShort === 'utmstack.tfa.method');
    this.checking = true;
    this.accountService.identity().then(account => {
      this.tfaService.initTfa({
        method: tfaMethod.confParamValue
      }).subscribe(() => {
        this.checking = false;
        this.isChecked.next(true);
      }, (error) => {
        this.checking = false;
        if (error.status === 400) {
          this.utmToastService.showError('Error initializing TFA',
            'An error occurred while initializing two-factor authentication, please check configuration and try again');
        } else {
          this.utmToastService.showError('Error initializing TFA',
            'An error occurred while initializing two-factor authentication, please contact with the support team');
        }
        this.isChecked.next(false);
      });
    });
  }
}
