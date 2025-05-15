import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConditionBuilderComponent } from './condition-builder.component';

describe('ConditionBuilderComponent', () => {
  let component: ConditionBuilderComponent;
  let fixture: ComponentFixture<ConditionBuilderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConditionBuilderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConditionBuilderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
