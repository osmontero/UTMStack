import {Component, Input, OnInit} from '@angular/core';
import {FormGroup} from '@angular/forms';
import {Observable, of} from 'rxjs';
import {catchError, finalize, map, tap} from 'rxjs/operators';
import {UtmNetScanService} from '../../../../assets-discover/shared/services/utm-net-scan.service';
import {NetScanType} from '../../../../assets-discover/shared/types/net-scan.type';
import {UtmToastService} from '../../../../shared/alert/utm-toast.service';
import {InputClassResolve} from '../../../../shared/util/input-class-resolve';

@Component({
  selector: 'app-action-builder',
  templateUrl: './action-builder.component.html',
  styleUrls: ['./action-builder.component.css']
})
export class ActionBuilderComponent implements OnInit {

  @Input() group: FormGroup;
  agents: any[];

  platforms$: Observable<string[]>;
  loadingPlatforms = false;
  agents$: Observable<NetScanType[]>;
  loadingAgents = false;
  noPlatforms = false;

  predefinedActions = [
    { icon: '📁', label: 'Create Incident', description: 'Creates a new incident' },
    { icon: '✅', label: 'Change Status to "under_review"', description: 'Marks alert as under review' },
    { icon: '📧', label: 'Send Email', description: 'Send a notification email' },
  ];

  workflow: any[] = [];

  constructor(private utmNetScanService: UtmNetScanService,
              public inputClass: InputClassResolve,
              private utmToastService: UtmToastService) { }

  ngOnInit() {
    this.platforms$ = this.getPlatforms();
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


  addToWorkFlow(action: any) {
    this.workflow.push(action);
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

}
