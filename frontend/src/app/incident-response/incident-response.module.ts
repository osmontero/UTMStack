import {CommonModule} from '@angular/common';
import {NgModule} from '@angular/core';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {NgSelectModule} from '@ng-select/ng-select';
import {TranslateModule} from '@ngx-translate/core';
import {AlertManagementSharedModule} from '../data-management/alert-management/shared/alert-management-shared.module';
import {UtmSharedModule} from '../shared/utm-shared.module';
import {CreateWorkflowComponent} from './create-workflow/create-workflow.component';
import {
  IncidentResponseAutomationComponent
} from './incident-response-automation/incident-response-automation.component';
import {IncidentResponseCommandComponent} from './incident-response-command/incident-response-command.component';
import {IncidentResponseRoutingModule} from './incident-response-routing.module';
import {IncidentResponseViewComponent} from './incident-response-view/incident-response-view.component';
import { PlaybookBuilderComponent } from './playbook-builder/playbook-builder.component';
import {IncidentResponseSharedModule} from './shared/incident-response-shared.module';
import {WorkflowActionsService} from './shared/services/workflow-actions.service';

@NgModule({
  declarations:
    [
      IncidentResponseViewComponent, IncidentResponseCommandComponent,
      IncidentResponseAutomationComponent,
      CreateWorkflowComponent,
      PlaybookBuilderComponent
    ],
  imports:
    [
      CommonModule,
      IncidentResponseRoutingModule,
      IncidentResponseSharedModule,
      UtmSharedModule,
      NgbModule,
      NgSelectModule,
      FormsModule,
      AlertManagementSharedModule,
      TranslateModule,
      ReactiveFormsModule
    ],
  entryComponents: [],
})
export class IncidentResponseModule {
}
