import { Component, OnInit } from '@angular/core';
import {WorkflowActionsService} from '../../services/workflow-actions.service';
import {ActionSidebarService} from './action-sidebar.service';
import {ActionTerminalComponent} from "../action-terminal/action-terminal.component";
import {ModalService} from "../../../../core/modal/modal.service";
import {NgbModal} from "@ng-bootstrap/ng-bootstrap";

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
              public actionSidebarService: ActionSidebarService,
              private modalService: NgbModal) { }

  ngOnInit() {
    this.actionSidebarService.loadData({
      ...this.request,
    });
  }

  addToWorkFlow(action: any) {
    this.workFlowActionService.addActions(action);
  }

  searchReport($event: string ) {
    this.actionSidebarService.loadData({
      ...this.request,
      'label.contains': $event.toString().toString()
    });
  }

  openActionSidebar() {
    const dialogRef = this.modalService.open(ActionTerminalComponent, {centered: true, size: 'lg'});

    dialogRef.result.then(
      result => {
        if (result) {
          this.workFlowActionService.addActions({
            ...result
          });
        }
      },
      reason => {
        console.log('Modal cerrado por:', reason);
      }
    );
  }
}
