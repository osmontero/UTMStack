import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionSidebarComponent } from './action-sidebar.component';

describe('ActionSidebarComponent', () => {
  let component: ActionSidebarComponent;
  let fixture: ComponentFixture<ActionSidebarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ActionSidebarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionSidebarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
