import { Component, OnInit } from '@angular/core';
import {NgbActiveModal, NgbModal} from '@ng-bootstrap/ng-bootstrap';
import {INCIDENT_AUTOMATION_ALERT_FIELDS} from '../../../../shared/constants/alert/alert-field.constant';

@Component({
  selector: 'app-action-terminal',
  templateUrl: './action-terminal.component.html',
  styleUrls: ['./action-terminal.component.css']
})
export class ActionTerminalComponent implements OnInit {

  alertFields = INCIDENT_AUTOMATION_ALERT_FIELDS;
  command: any;

  constructor(public activeModal: NgbActiveModal, ) { }

  ngOnInit() {
  }

  insertVariablePlaceholder($event: string) {
    this.command += `$[${$event}]`;
  }

  insertFieldPlaceholder(field: string) {
    this.command += `$(${field})`;
  }

  close() {
    this.activeModal.close(true);
  }
}
