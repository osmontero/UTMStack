import {Component, Input, OnInit} from '@angular/core';
import {AgentType} from '../../../../shared/types/agent/agent.type';

@Component({
  selector: 'app-interactive-console',
  templateUrl: './interactive-console.component.html',
  styleUrls: ['./interactive-console.component.css']
})
export class InteractiveConsoleComponent implements OnInit {

  @Input() agent: AgentType;

  constructor() { }

  ngOnInit() {
  }

}
