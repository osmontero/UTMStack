import {Component, Input, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {filter} from 'rxjs/operators';
import {AgentType} from '../../../../shared/types/agent/agent.type';
import {AgentSidebarService} from '../agent-sidebar/agent-sidebar.service';

@Component({
  selector: 'app-interactive-console',
  templateUrl: './interactive-console.component.html',
  styleUrls: ['./interactive-console.component.scss']
})
export class InteractiveConsoleComponent implements OnInit {

  agent$: Observable<AgentType>;

  constructor(private agentSidebarService: AgentSidebarService) { }

  ngOnInit() {
    this.agent$ = this.agentSidebarService.selectedAgent$
      .pipe(filter(agent => !!agent));
  }
}
