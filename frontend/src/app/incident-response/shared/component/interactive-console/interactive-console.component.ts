import {Component, Input, OnInit} from '@angular/core';
import {NetScanType} from '../../../../assets-discover/shared/types/net-scan.type';

@Component({
  selector: 'app-interactive-console',
  templateUrl: './interactive-console.component.html',
  styleUrls: ['./interactive-console.component.css']
})
export class InteractiveConsoleComponent implements OnInit {

  @Input() agent: NetScanType;

  constructor() { }

  ngOnInit() {
  }

}
