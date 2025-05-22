import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export function containsVariable(variables: string[]): ValidatorFn {
  console.log('contains variable', variables);
  return (control: AbstractControl): ValidationErrors | null => {
    const value = control.value || '';
    const hasVariable = variables.some(v => value.includes(v));
    return hasVariable ? null : { noVariableUsed: true };
  };
}
