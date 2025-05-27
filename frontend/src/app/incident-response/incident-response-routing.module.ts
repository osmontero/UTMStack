import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {UserRouteAccessService} from '../core/auth/user-route-access-service';
import {ADMIN_ROLE} from '../shared/constants/global.constant';
import {IncidentResponseAutomationComponent} from './incident-response-automation/incident-response-automation.component';
import {IncidentResponseViewComponent} from './incident-response-view/incident-response-view.component';
import {PlaybookBuilderComponent} from './playbook-builder/playbook-builder.component';
import {PlaybooksComponent} from "./playbooks/playbooks.component";

const routes: Routes = [
  {path: '', redirectTo: 'audit'},
  {
    path: 'audit',
    component: IncidentResponseViewComponent,
    canActivate: [UserRouteAccessService],
    data: {authorities: [ADMIN_ROLE]}
  },
  {
    path: 'automation',
    component: IncidentResponseAutomationComponent,
    canActivate: [UserRouteAccessService],
    data: {authorities: [ADMIN_ROLE]}
  },
  {
    path: 'create',
    component: PlaybookBuilderComponent,
    canActivate: [UserRouteAccessService],
    data: {authorities: [ADMIN_ROLE]}
  },
  {
    path: 'playbooks',
    component: PlaybooksComponent,
    canActivate: [UserRouteAccessService],
    data: {authorities: [ADMIN_ROLE]}
  },

];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IncidentResponseRoutingModule {
}

