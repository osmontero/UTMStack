import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormGroup} from '@angular/forms';
import {Observable} from 'rxjs';
import {map, tap} from 'rxjs/operators';
import {INCIDENT_AUTOMATION_ALERT_FIELDS} from '../../../../shared/constants/alert/alert-field.constant';
import {ALERT_INDEX_PATTERN} from '../../../../shared/constants/main-index-pattern.constant';
import {ElasticSearchIndexService} from '../../../../shared/services/elasticsearch/elasticsearch-index.service';

@Component({
  selector: 'app-condition-item',
  templateUrl: './condition-item.component.html',
  styleUrls: ['./condition-item.component.css']
})
export class ConditionItemComponent implements OnInit {
  @Input() group: FormGroup;
  @Input() index: number;
  @Output() delete: EventEmitter<number> = new EventEmitter();

  alertFields = INCIDENT_AUTOMATION_ALERT_FIELDS;
  req = {
    page: 0,
    size: 10,
    indexPattern: ALERT_INDEX_PATTERN,
    keyword: ''
  };
  values$: Observable<string[]>;
  loading = false;

  constructor(private elasticSearchIndexService: ElasticSearchIndexService) { }

  ngOnInit() {
  }

  onScroll() {
    this.req = {
      ...this.req,
      size: this.req.size + 10
    };

    this.getValues();
  }

  getValuesForField(key: string | null) {

    this.req = {
      ...this.req,
      keyword: key ? key + '.keyword' : this.req.keyword
    };
    this.getValues();
  }

  getValues() {
    this.loading = true;
    this.values$ = this.elasticSearchIndexService.getElasticFieldValues(this.req)
      .pipe(
        tap((res) => this.loading = !this.loading),
        map(res => res.body));
  }

  removeRuleCondition(i: any) {
    this.delete.emit(this.index);
  }
}
