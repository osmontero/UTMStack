import {HttpResponse} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {BehaviorSubject, of} from 'rxjs';
import {catchError, filter, finalize, map, switchMap, tap} from 'rxjs/operators';
import {UtmToastService} from '../../../../shared/alert/utm-toast.service';
import {
  IncidentResponseActionTemplate,
  IncidentResponseActionTemplateService
} from '../../services/incident-response-action-template.service';
import {UtmNetScanService} from "../../../../assets-discover/shared/services/utm-net-scan.service";
import {NetScanType} from "../../../../assets-discover/shared/types/net-scan.type";

@Injectable({
  providedIn: 'root'
})
export class AgentSidebarService {

  private request = new BehaviorSubject<any>(null);
  private loading = new BehaviorSubject<boolean>(false);

  request$ = this.request.asObservable();
  loading$ = this.loading.asObservable();

  agents$ = this.request$
    .pipe(
      filter(request => !!request),
      tap(() => this.loading.next(true)),
      switchMap((request) => this.utmNetScanService.query(request)
        .pipe(
          map((response: HttpResponse<NetScanType[]>) => response.body),
          catchError(() => {
            this.toastService.showError('Error', 'Failed to load agents');
            return of([]);
          }),
          finalize(() => this.loading.next(false))
        )
      ),
    );

  constructor(private utmNetScanService: UtmNetScanService,
              private toastService: UtmToastService) {}

  loadData(request: any) {
    this.request.next(request);
  }

  reset() {
    this.request.next(null);
  }

}
