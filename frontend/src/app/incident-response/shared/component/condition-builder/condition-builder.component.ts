import {Component, Input, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {ElasticOperatorsEnum} from '../../../../shared/enums/elastic-operators.enum';
import {IncidentRuleType} from '../../type/incident-rule.type';

@Component({
  selector: 'app-condition-builder',
  templateUrl: './condition-builder.component.html',
  styleUrls: ['./condition-builder.component.css']
})
export class ConditionBuilderComponent implements OnInit {

  @Input() alert: any;
  @Input() group: FormGroup;
  @Input() rule: IncidentRuleType;

  constructor(private fb: FormBuilder) { }

  ngOnInit() {
    if (this.rule) {
      for (const condition of this.rule.conditions) {
        const ruleCondition = this.fb.group({
          field: [condition.field, Validators.required],
          value: [condition.value, Validators.required],
          operator: [condition.operator]
        });
        this.ruleConditions.push(ruleCondition);
      }
    } else {
      this.addRuleCondition();
    }
  }

  get ruleConditions() {
    return this.group.get('conditions') as FormArray;
  }

  addRuleCondition() {
    const ruleCondition = this.fb.group({
      field: ['', Validators.required],
      value: ['', Validators.required],
      operator: [ElasticOperatorsEnum.IS]
    });

    this.ruleConditions.push(ruleCondition);
  }

  removeRuleCondition(index: number) {
    this.ruleConditions.removeAt(index);
  }

}
