import {Component, Input, OnInit} from '@angular/core';
import {NetScanType} from '../../../../assets-discover/shared/types/net-scan.type';

@Component({
  selector: 'app-agent-info',
  templateUrl: './agent-info.component.html',
  styleUrls: ['./agent-info.component.css']
})
export class AgentInfoComponent implements OnInit {
  @Input() agent: NetScanType;

  constructor() { }

  ngOnInit() {
  }

  statusClass(status: string): string {
    switch (status) {
      case 'CHECK':
        return 'badge badge-warning';
      case 'OK':
        return 'badge badge-success';
      case 'DOWN':
        return 'badge badge-danger';
      default:
        return 'badge badge-secondary';
    }
  }
}
