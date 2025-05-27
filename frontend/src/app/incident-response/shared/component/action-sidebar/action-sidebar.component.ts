import { Component, OnInit } from '@angular/core';
import {WorkflowActionsService} from '../../services/workflow-actions.service';
import {ActionSidebarService} from './action-sidebar.service';

@Component({
  selector: 'app-action-sidebar',
  templateUrl: './action-sidebar.component.html',
  styleUrls: ['./action-sidebar.component.css']
})
export class ActionSidebarComponent implements OnInit {

  request = {
    page: 0,
    size: 25,
  };

  searching: any;

  constructor(private workFlowActionService: WorkflowActionsService,
              public actionSidebarService: ActionSidebarService) { }

  ngOnInit() {
    this.actionSidebarService.setRequest({
      ...this.request,
    });
  }

  addToWorkFlow(action: any) {
    this.workFlowActionService.addActions(action);
  }

  searchReport($event: string | number) {

  }

  openAddCustomActionModal() {

  }
}
