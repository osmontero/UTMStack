import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PlaybooksComponent } from './playbooks.component';

describe('PlaybooksComponent', () => {
  let component: PlaybooksComponent;
  let fixture: ComponentFixture<PlaybooksComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PlaybooksComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PlaybooksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
