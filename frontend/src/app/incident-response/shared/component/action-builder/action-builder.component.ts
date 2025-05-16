import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {FormGroup} from '@angular/forms';
import {Observable, of, Subject} from 'rxjs';
import {catchError, finalize, map, takeUntil, tap} from 'rxjs/operators';
import {UtmNetScanService} from '../../../../assets-discover/shared/services/utm-net-scan.service';
import {NetScanType} from '../../../../assets-discover/shared/types/net-scan.type';
import {UtmToastService} from '../../../../shared/alert/utm-toast.service';
import {InputClassResolve} from '../../../../shared/util/input-class-resolve';
import {WorkflowActionsService} from '../../services/workflow-actions.service';

@Component({
  selector: 'app-action-builder',
  templateUrl: './action-builder.component.html',
  styleUrls: ['./action-builder.component.css']
})
export class ActionBuilderComponent implements OnInit, OnDestroy {

  @Input() group: FormGroup;
  agents: any[];

  platforms$: Observable<string[]>;
  loadingPlatforms = false;
  agents$: Observable<NetScanType[]>;
  loadingAgents = false;
  noPlatforms = false;

  workflow$: Observable<any[]>;
  destroy$: Subject<void> = new Subject<void>();

  constructor(private utmNetScanService: UtmNetScanService,
              public inputClass: InputClassResolve,
              private utmToastService: UtmToastService,
              private workflowActionsService: WorkflowActionsService) { }

  ngOnInit() {
    this.platforms$ = this.getPlatforms();

    this.workflow$ = this.workflowActionsService.actions$
      .pipe(takeUntil(this.destroy$));

    this.workflow$.subscribe( w => console.log(w));
  }

  getPlatforms() {
    this.loadingPlatforms = true;
    return this.utmNetScanService.getAssetsPlatforms().pipe(
      map(res => res.body || []),
      tap(platforms => {
        this.noPlatforms = platforms.length === 0;
      }),
      catchError(err => {
        this.utmToastService.showError('Error fetching', 'An error has occurred while fetching platforms');
        return of([]);
      }),
      finalize(() => {
        this.loadingPlatforms = false;
      })
    );
  }

  getAgents(platform: any) {
    this.group.get('excludedAgents').reset();
    this.group.get('defaultAgent').reset();

    this.loadingAgents = true;
    this.agents$ = this.fetchAgents(platform);
  }

  onChangeToggle($event) {
    if ($event) {
      this.group.get('excludedAgents').setValue([]);
    } else {
      this.group.get('defaultAgent').setValue('');
    }
    this.group.get('agentType').setValue($event);
  }

  onChangeExclude($event: any) {
    const hostnames = $event.map(value => value.assetName);
    this.group.get('excludedAgents').setValue(hostnames);
  }

  fetchAgents(platform: string) {
    return this.utmNetScanService.query({page: 0, size: 10000, agent: true, osPlatform: platform})
      .pipe(
        map(res => res.body || []),
        tap(agents => {
          if (agents.length === 1) {
            this.group.get('excludedAgents').disable();
          } else {
            this.group.get('excludedAgents').enable();
          }
        }),
        catchError(err => {
          this.utmToastService.showError('Error fetching', 'An error has occurred while fetching agents');
          return of([]);
        }),
        finalize(() => {
          this.loadingAgents = false;
        })
      );
  }

  openActionSidebar() {

  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  removeAction() {

  }
}
