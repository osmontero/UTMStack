import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Router} from "@angular/router";
import {NgbActiveModal} from '@ng-bootstrap/ng-bootstrap';
import {UtmToastService} from '../../../../shared/alert/utm-toast.service';
import {InputClassResolve} from '../../../../shared/util/input-class-resolve';
import {IncidentActionType} from '../../type/incident-action.type';

@Component({
  selector: 'app-incident-response-command-create',
  templateUrl: './new-playbook.component.html',
  styleUrls: ['./new-playbook.component.scss']
})
export class NewPlaybookComponent implements OnInit {
  @Input() action: IncidentActionType;
  @Output() actionCreated = new EventEmitter<IncidentActionType>();

  constructor(public activeModal: NgbActiveModal,
              public inputClass: InputClassResolve,
              public utmToastService: UtmToastService,
              private router: Router) {
  }

  ngOnInit() {}

  createNewPlaybook() {
    this.router.navigate(['/incident-response/create']);
    this.activeModal.close();
  }

  searchByRule($event: string | number) {

  }
}
