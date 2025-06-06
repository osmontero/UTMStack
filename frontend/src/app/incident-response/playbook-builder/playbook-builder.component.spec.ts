import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PlaybookBuilderComponent } from './playbook-builder.component';

describe('PlaybookBuilderComponent', () => {
  let component: PlaybookBuilderComponent;
  let fixture: ComponentFixture<PlaybookBuilderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PlaybookBuilderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PlaybookBuilderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
