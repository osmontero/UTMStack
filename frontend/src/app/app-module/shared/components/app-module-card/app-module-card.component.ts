import {ChangeDetectionStrategy, Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {Subject} from 'rxjs';
import {takeUntil} from 'rxjs/operators';
import {VersionType, VersionTypeService} from '../../../../shared/services/util/version-type.service';
import {UtmModulesEnum} from '../../enum/utm-module.enum';
import {UtmModuleType} from '../../type/utm-module.type';

@Component({
  selector: 'app-app-module-card',
  templateUrl: './app-module-card.component.html',
  styleUrls: ['./app-module-card.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AppModuleCardComponent implements OnInit, OnDestroy {
  @Input() module: UtmModuleType;
  @Output() showModuleIntegration = new EventEmitter<UtmModuleType>();
  versionType = VersionType;
  version: VersionType;
  modules = UtmModulesEnum;
  destroy$: Subject<void> = new Subject<void>();

  constructor(private versionTypeService: VersionTypeService) {
  }

  ngOnInit() {
    this.versionTypeService.versionType$
    .pipe(takeUntil(this.destroy$))
      .subscribe(versionType => this.version = versionType);
  }

  showIntegration() {
    this.showModuleIntegration.emit(this.module);
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
