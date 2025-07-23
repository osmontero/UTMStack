import { Component, OnInit } from '@angular/core';
import {AgentSidebarService} from './agent-sidebar.service';

@Component({
  selector: 'app-agent-sidebar',
  templateUrl: './agent-sidebar.component.html',
  styleUrls: ['./agent-sidebar.component.scss']
})
export class AgentSidebarComponent implements OnInit {
  searching: any;
  request = {
    page: 0,
    size: 5,
    agent: true
  };

  constructor(public agentSidebarService: AgentSidebarService) { }

  ngOnInit() {
    this.agentSidebarService.loadData({
      ...this.request
    });
  }

  searchAgent($event: string | number) {
    this.agentSidebarService.loadData({
      ...this.request,
      page: 0,
      assetIpMacName: $event.toString()
    });
  }

  onScroll() {

  }

  trackByFn(index: number, item: any) {
    return item.id || index;
  }

  agentDetail(action: any) {

  }
}
