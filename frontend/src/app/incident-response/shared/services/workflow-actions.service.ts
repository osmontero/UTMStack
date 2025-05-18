import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WorkflowActionsService {

  private actionsBehaviorSubject: BehaviorSubject<any[]> = new BehaviorSubject([]);
  actions$ = this.actionsBehaviorSubject.asObservable();

  setActions(action: any) {
    this.actionsBehaviorSubject.next([{...action}]);
  }

  clear(){
    this.actionsBehaviorSubject.next([]);
  }
}
