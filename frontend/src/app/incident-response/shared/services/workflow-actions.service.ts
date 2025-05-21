import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
import {
  ActionConditionalEnum
} from '../component/action-conditional/action-conditional.component';

export interface WorkflowAction {
  label: string;
  description: string;
  conditional?: { key: ActionConditionalEnum, value: string};
}

@Injectable({
  providedIn: 'root'
})
export class WorkflowActionsService {

  private actionsBehaviorSubject: BehaviorSubject<WorkflowAction[]> = new BehaviorSubject([]);
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
