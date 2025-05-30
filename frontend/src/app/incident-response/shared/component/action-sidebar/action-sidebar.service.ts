import {HttpResponse} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
import {map, switchMap} from 'rxjs/operators';
import {
  IncidentResponseActionTemplate,
  IncidentResponseActionTemplateService
} from '../../services/incident-response-action-template.service';

@Injectable({
  providedIn: 'root'
})
export class ActionSidebarService {

  private request = new BehaviorSubject<any>(null);
  request$ = this.request.asObservable();

  actionTemplates$ = this.request$
    .pipe(
      switchMap((request) => this.incidentResponseActionTemplateService.query(request)),
      map((response: HttpResponse<IncidentResponseActionTemplate[]>) => response.body)
    );

  constructor(private incidentResponseActionTemplateService: IncidentResponseActionTemplateService) {}

  loadData(request: any) {
    this.request.next(request);
  }

}
