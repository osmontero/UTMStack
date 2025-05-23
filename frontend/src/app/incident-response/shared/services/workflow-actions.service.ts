import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import {ActionConditionalEnum} from '../component/action-conditional/action-conditional.component';
import {map} from "rxjs/operators";

export interface WorkflowAction {
  id: number;
  label: string;
  description: string;
  command: string;
  conditional?: { key: ActionConditionalEnum, value: string };
}

@Injectable({
  providedIn: 'root'
})
export class WorkflowActionsService {

  private actionsBehaviorSubject: BehaviorSubject<WorkflowAction[]> = new BehaviorSubject([]);
  actions$ = this.actionsBehaviorSubject.asObservable();

  readonly command$: Observable<string> = this.actions$.pipe(
    map(actions => {
      if (actions.length === 1) {
        return actions[0].command;
      }

      return actions.map((action, index) => {
        const operator = index === 0 ? ''
          : action.conditional.key === ActionConditionalEnum.SUCCESS ? '&&'
            : action.conditional.key === ActionConditionalEnum.FAILURE ? '||'
              : ';';

        return `${operator} ${action.command}`.trim();
      }).join(' ').trim();
    })
  );

  setActions(action: any) {
    const actions = this.actionsBehaviorSubject.value ? this.actionsBehaviorSubject.value : [];

    this.actionsBehaviorSubject.next([...actions, {
      ...action,
      conditional: { key: ActionConditionalEnum.ALWAYS, value: ';'},
    }]);
  }

  updateAction(action: any) {
    const actions = this.actionsBehaviorSubject.value ? this.actionsBehaviorSubject.value : [];

    const index = actions.findIndex((act: any) => act.id === action.id);

    const newActions = [...actions];
    newActions[index] = {
      ...action,
    };

    this.actionsBehaviorSubject.next(newActions);

  }

  deleteAction(action: any) {
    const actions = this.actionsBehaviorSubject.value ? this.actionsBehaviorSubject.value : [];
    this.actionsBehaviorSubject.next(actions.filter(act => act !== action));
  }

  clear() {
    this.actionsBehaviorSubject.next([]);
  }
}
