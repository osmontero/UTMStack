import {Component, Input, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {NetScanType} from '../../assets-discover/shared/types/net-scan.type';
import {AgentSidebarService} from '../shared/component/agent-sidebar/agent-sidebar.service';
import {filter} from "rxjs/operators";

@Component({
  selector: 'app-console-interactive',
  templateUrl: './interactive-console.component.html',
  styleUrls: ['./interactive-console.component.scss']
})
export class InteractiveConsoleComponent implements OnInit {

  agentSelected$: Observable<NetScanType>;

  constructor(private agentSidebarService: AgentSidebarService) { }

  ngOnInit() {

    this.agentSelected$ = this.agentSidebarService.selectedAgent$
      .pipe(filter(agent => !!agent));

  }

}
