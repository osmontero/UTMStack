import { Component, OnInit } from '@angular/core';
import {WorkflowActionsService} from '../../services/workflow-actions.service';

@Component({
  selector: 'app-action-sidebar',
  templateUrl: './action-sidebar.component.html',
  styleUrls: ['./action-sidebar.component.css']
})
export class ActionSidebarComponent implements OnInit {

  predefinedActions = [
    {
      id: 1,
      icon: '📁',
      label: 'Create Incident',
      description: 'Creates a new incident',
      command: 'Invoke-Incident -Create -Title "New Security Incident"'
    },
    {
      id: 2,
      icon: '✅',
      label: 'Change Status to "under_review"',
      description: 'Marks alert as under review',
      command: 'Update-AlertStatus -Id $AlertId -Status "under_review"'
    },
    {
      id: 3,
      icon: '📧',
      label: 'Send Email',
      description: 'Send a notification email',
      command: 'Send-Mail -To "secops@example.com" -Subject "Alert Under Review" -Body "An alert has been flagged for review."'
    },
    {
      id: 4,
      icon: '🚫',
      label: 'Block IP',
      description: 'Blocks the source IP address at the firewall',
      command: 'Block-IP -Address $SourceIP'
    },
    {
      id: 5,
      icon: '🔄',
      label: 'Restart Service',
      description: 'Restarts a critical service on affected asset',
      command: 'Restart-Service -Name "nginx" -ComputerName $TargetHost'
    },
    {
      id: 6,
      icon: '📝',
      label: 'Log to SIEM',
      description: 'Sends a log entry to SIEM system',
      command: 'Write-SIEMLog -Message "Incident created for $AlertId"'
    }
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
