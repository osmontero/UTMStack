import { Component, OnInit } from '@angular/core';
import {WorkflowActionsService} from '../../services/workflow-actions.service';

@Component({
  selector: 'app-action-sidebar',
  templateUrl: './action-sidebar.component.html',
  styleUrls: ['./action-sidebar.component.css']
})
export class ActionSidebarComponent implements OnInit {

  predefinedActions = [
    { icon: '📁', label: 'Create Incident', description: 'Creates a new incident' },
    { icon: '✅', label: 'Change Status to "under_review"', description: 'Marks alert as under review' },
    { icon: '📧', label: 'Send Email', description: 'Send a notification email' },
  ];
  searching: any;

  constructor(private workFlowActionService: WorkflowActionsService) { }

  ngOnInit() {
  }

  addToWorkFlow(action: any) {
    this.workFlowActionService.setActions(action);
  }

  searchReport($event: string | number) {

  }

  openAddCustomActionModal() {

  }
}
