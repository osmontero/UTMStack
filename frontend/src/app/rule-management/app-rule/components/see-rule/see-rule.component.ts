
import { Component, Input } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { Rule } from '../../../models/rule.model';
import * as yaml from 'js-yaml';

@Component({
  selector: 'app-see-rule',
  templateUrl: './see-rule.component.html',
  styleUrls: ['./see-rule.component.scss'],
})
export class SeeRuleComponent {
 @Input() rowDocument: Rule;

  copied = false;

  get yamlString(): string {
    try {
      return yaml.dump(this.rowDocument, { indent: 2 });
    } catch (e) {
      console.log(e)
      return 'Error parsing YAML';
    }
  }

  copyToClipboard() {
    navigator.clipboard.writeText(this.yamlString).then(() => {
      this.copied = true;
      setTimeout(() => (this.copied = false), 1500);
    });
  }}



