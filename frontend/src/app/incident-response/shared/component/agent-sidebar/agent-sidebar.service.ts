import {HttpResponse} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {BehaviorSubject, of} from 'rxjs';
import {catchError, filter, finalize, map, switchMap, tap} from 'rxjs/operators';
import {NetScanType} from '../../../../assets-discover/shared/types/net-scan.type';
import {UtmToastService} from '../../../../shared/alert/utm-toast.service';
import {UtmAgentManagerService} from '../../../../shared/services/agent/utm-agent-manager.service';
import {AgentType} from '../../../../shared/types/agent/agent.type';

@Injectable({
  providedIn: 'root'
})
export class AgentSidebarService {

  private request = new BehaviorSubject<any>(null);
  private loading = new BehaviorSubject<boolean>(false);
  private selectedAgent = new BehaviorSubject<AgentType>(null);

  request$ = this.request.asObservable();
  loading$ = this.loading.asObservable();
  selectedAgent$ = this.selectedAgent.asObservable();

  agents$ = this.request$
    .pipe(
      filter(request => !!request),
      tap(() => this.loading.next(true)),
      switchMap((request) => this.utmAgentManagerService.getAgents(request)
        .pipe(
          map((response: HttpResponse<AgentType[]>) => response.body),
          catchError(() => {
            this.toastService.showError('Error', 'Failed to load agents');
            return of([] as AgentType[]);
          }),
          finalize(() => this.loading.next(false))
        )
      ),
    );

  constructor(private utmAgentManagerService: UtmAgentManagerService,
              private toastService: UtmToastService) {}

  loadData(request: any) {
    this.request.next(request);
  }

  selectAgent(agent: AgentType) {
    this.selectedAgent.next(agent);
  }

  reset() {
    this.request.next(null);
  }

}
