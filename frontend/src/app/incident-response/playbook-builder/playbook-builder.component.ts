import {Component, ElementRef, EventEmitter, HostListener, Input, OnInit, Output, ViewChild} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NgbActiveModal} from '@ng-bootstrap/ng-bootstrap';
import {debounceTime} from 'rxjs/operators';
import {NetScanType} from '../../assets-discover/shared/types/net-scan.type';
import {ALERT_NAME_FIELD} from '../../shared/constants/alert/alert-field.constant';
import {ElasticOperatorsEnum} from '../../shared/enums/elastic-operators.enum';
import {PrefixElementEnum} from '../../shared/enums/prefix-element.enum';
import {getValueFromPropertyPath} from '../../shared/util/get-value-object-from-property-path.util';
import {InputClassResolve} from '../../shared/util/input-class-resolve';
import {createElementPrefix, getElementPrefix} from '../../shared/util/string-util';
import {IncidentResponseRuleService} from '../shared/services/incident-response-rule.service';
import {IncidentRuleType} from '../shared/type/incident-rule.type';

@Component({
  selector: 'app-playbook-builder',
  templateUrl: './playbook-builder.component.html',
  styleUrls: ['./playbook-builder.component.css']
})
export class PlaybookBuilderComponent implements OnInit {

  @Input() alert: any;
  @Input() rule: IncidentRuleType;
  @ViewChild('autocomplete') autocomplete: ElementRef;
  @Output() ruleCreated = new EventEmitter<boolean>();
  step = 1;
  stepCompleted: number[] = [];
  creating = false;
  formRule: FormGroup;
  agents: NetScanType[];
  platforms: string[];
  command = '';
  exist = true;
  typing = true;
  rulePrefix: string = createElementPrefix(PrefixElementEnum.INCIDENT_RESPONSE_AUTOMATION);
  maxLength = 512;
  viewportHeight: number;

  workflow: any[] = [];

  constructor(private incidentResponseRuleService: IncidentResponseRuleService,
              public activeModal: NgbActiveModal,
              private fb: FormBuilder,
              public inputClass: InputClassResolve) {

    this.formRule = this.fb.group({
      id: [null],
      name: ['', Validators.required],
      description: ['', Validators.required],
      conditions: this.fb.array([]),
      command: ['', Validators.required],
      active: [true],
      agentType: [false],
      excludedAgents: [[]],
      defaultAgent: [''],
      agentPlatform: ['', Validators.required]
    });
    this.viewportHeight = window.innerHeight;
  }

  @HostListener('window:resize', ['$event'])
  onResize(event: Event) {
    this.viewportHeight = window.innerHeight;
  }

  ngOnInit() {
    if (this.rule) {
      this.exist = false;
      this.typing = false;
      this.rulePrefix = getElementPrefix(this.rule.name);
      this.formRule.patchValue(this.rule, {emitEvent: false});
      const name = this.formRule.get('name').value;
      this.formRule.get('name').setValue(this.replacePrefixInName(name));
      for (const condition of this.rule.conditions) {
        // this.getValuesForField(condition.field);
        const ruleCondition = this.fb.group({
          field: [condition.field, Validators.required],
          value: [condition.value, Validators.required],
          operator: [condition.operator]
        });
        this.command = this.rule.command;
        this.ruleConditions.push(ruleCondition);
        this.formRule.get('excludedAgents').setValue(this.rule.excludedAgents);
        this.formRule.get('agentType').setValue(this.rule.excludedAgents.length === 0 && this.rule.defaultAgent !== '');
        this.formRule.get('defaultAgent').setValue(this.rule.defaultAgent);
      }
    } else if (this.alert) {
      const alertName = this.getValueFromAlert(ALERT_NAME_FIELD);
      const ruleName = this.rulePrefix + alertName;
      this.formRule.get('name').setValue(alertName);
      this.searchRule(ruleName);
      this.addRuleCondition();
      this.ruleConditions.at(0).get('field').setValue(ALERT_NAME_FIELD);
      this.ruleConditions.at(0).get('value').setValue(alertName);
      if (alertName.toLowerCase().includes('window')) {
        this.formRule.get('agentPlatform').setValue('windows');
      }
    }
    this.formRule.get('name').valueChanges.pipe(debounceTime(1000)).subscribe(value => {
      this.searchRule(this.rulePrefix + value);
    });
  }

  get ruleConditions() {
    return this.formRule.get('conditions') as FormArray;
  }


  replacePrefixInName(name: string) {
    return name.replace(this.rulePrefix, '');
  }

  addRuleCondition() {
    const ruleCondition = this.fb.group({
      field: ['', Validators.required],
      value: ['', Validators.required],
      operator: [ElasticOperatorsEnum.IS]
    });

    this.ruleConditions.push(ruleCondition);
  }

  getValueFromAlert(field: string) {
    return getValueFromPropertyPath(this.alert, field, null);
  }

  searchRule(rule: string) {
    this.typing = true;
    this.exist = true;
    setTimeout(() => {
      const req = {
        'name.contains': rule
      };
      this.incidentResponseRuleService.query(req).subscribe(response => {
        this.exist = response.body.length > 0;
        this.typing = false;
      });
    }, 1000);
  }

  getMenuHeight() {
    return 100 - ((150 / this.viewportHeight) * 100) + 'vh';
  }

}
