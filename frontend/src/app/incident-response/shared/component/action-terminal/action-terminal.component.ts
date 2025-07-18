import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import {ALERT_FIELDS, INCIDENT_AUTOMATION_ALERT_FIELDS} from '../../../../shared/constants/alert/alert-field.constant';

@Component({
  selector: 'app-action-terminal',
  templateUrl: './action-terminal.component.html',
  styleUrls: ['./action-terminal.component.scss']
})
export class ActionTerminalComponent implements OnInit {

  form: FormGroup;
  alertFields = ALERT_FIELDS;
  command: any;

  constructor(public activeModal: NgbActiveModal,
              private fb: FormBuilder) {
    this.alertFields = this.alertFields.reduce((acc: any[], field) => {
      if (typeof field === 'object' && field !== null && 'fields' in field) {
        return acc.concat(field.fields);
      }

      return acc.concat(field);
    }, []);
  }

  ngOnInit() {
    this.form = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(5)]],
      description: ['', [Validators.required, Validators.minLength(5)]],
      command: ['', Validators.required],
    });
  }

  insertVariablePlaceholder($event: string) {
    this.form.get('command') .setValue(this.form.get('command').value + `$(${ $event })`);
  }

  insertFieldPlaceholder(field: string) {
    this.form.get('command') .setValue(this.form.get('command').value + `$(${ field })`);
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
