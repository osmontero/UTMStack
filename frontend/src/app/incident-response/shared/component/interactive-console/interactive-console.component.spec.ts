import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InteractiveConsoleComponent } from './interactive-console.component';

describe('InteractiveConsoleComponent', () => {
  let component: InteractiveConsoleComponent;
  let fixture: ComponentFixture<InteractiveConsoleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InteractiveConsoleComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InteractiveConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
