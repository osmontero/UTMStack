import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import {INCIDENT_AUTOMATION_ALERT_FIELDS} from '../../../../shared/constants/alert/alert-field.constant';

@Component({
  selector: 'app-action-terminal',
  templateUrl: './action-terminal.component.html',
  styleUrls: ['./action-terminal.component.css']
})
export class ActionTerminalComponent implements OnInit {

  form: FormGroup;
  alertFields = INCIDENT_AUTOMATION_ALERT_FIELDS;
  command: any;

  constructor(public activeModal: NgbActiveModal,
              private fb: FormBuilder) { }

  ngOnInit() {
    this.form = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(5)]],
      description: ['', [Validators.required, Validators.minLength(5)]],
      command: ['', Validators.required],
    });
  }

  insertVariablePlaceholder($event: string) {
    this.command += `$[${$event}]`;
  }

  insertFieldPlaceholder(field: string) {
    this.command += `$(${field})`;
  }

  close() {
    this.activeModal.dismiss();
  }

  create() {
    this.activeModal.close({
      ...this.form.value
    });
  }
}
