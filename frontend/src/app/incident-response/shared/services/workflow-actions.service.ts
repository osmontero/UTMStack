import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WorkflowActionsService {

  private actionsBehaviorSubject: BehaviorSubject<any[]> = new BehaviorSubject([]);
  actions$ = this.actionsBehaviorSubject.asObservable();

  setActions(action: any) {
    const actions = this.actionsBehaviorSubject.value ? this.actionsBehaviorSubject.value : [];
    this.actionsBehaviorSubject.next([...actions, action]);
  }

  deleteAction(action: any) {
    const actions = this.actionsBehaviorSubject.value ? this.actionsBehaviorSubject.value : [];
    this.actionsBehaviorSubject.next(actions.filter(act => act !== action));
  }

  clear(){
    this.actionsBehaviorSubject.next([]);
  }
}
