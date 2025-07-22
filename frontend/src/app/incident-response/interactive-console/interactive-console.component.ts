import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-interactive-console',
  templateUrl: './interactive-console.component.html',
  styleUrls: ['./interactive-console.component.scss']
})
export class InteractiveConsoleComponent implements OnInit {
  selectedAgent: any;

  constructor() { }

  ngOnInit() {
  }

}
