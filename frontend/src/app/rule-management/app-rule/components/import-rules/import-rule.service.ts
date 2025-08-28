import { Injectable } from '@angular/core';
import { Rule } from '../../../models/rule.model';

@Injectable()
export class ImportRuleService {

  isValidURL(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  private minWordsCheck(value: string, min: number, minLengthPerWord: number): boolean {
    if (!value) return false;
    const words = value.trim().split(/\s+/).filter(word => word.length >= minLengthPerWord);
    return words.length >= min;

  }

  isValidRule(obj: Rule): { isValid: boolean; errors: Record<string, string[]> } {
    const errors: Record<string, string[]> = {};

    if (!obj || typeof obj !== 'object') {
      return { isValid: false, errors: { rule: ['Rule object is missing or invalid'] } };
    }

    // dataTypes
    if (!Array.isArray(obj.dataTypes) || obj.dataTypes.length === 0) {
      errors['dataTypes'] = ['dataTypes are required'];
    }

    // name
    if (typeof obj.name !== 'string' || obj.name.trim() === '') {
      errors['name'] = ['Name is required'];
    } else if (!this.minWordsCheck(obj.name, 2, 3)) {
      errors['name'] = ['Name must contain between 2 and 3 words'];
    }

    // adversary
    if (typeof obj.adversary !== 'string' || obj.adversary.trim() === '') {
      errors['adversary'] = ['Adversary is required'];
    }

    // confidentiality
    if (typeof obj.confidentiality !== 'number') {
      errors['confidentiality'] = ['Confidentiality must be a number'];
    } else if (obj.confidentiality < 0 || obj.confidentiality > 3) {
      errors['confidentiality'] = ['Confidentiality must be between 0 and 3'];
    }

    // integrity
    if (typeof obj.integrity !== 'number') {
      errors['integrity'] = ['Integrity must be a number'];
    } else if (obj.integrity < 0 || obj.integrity > 3) {
      errors['integrity'] = ['Integrity must be between 0 and 3'];
    }

    // availability
    if (typeof obj.availability !== 'number') {
      errors['availability'] = ['Availability must be a number'];
    } else if (obj.availability < 0 || obj.availability > 3) {
      errors['availability'] = ['Availability must be between 0 and 3'];
    }

    // category
    if (typeof obj.category !== 'string' || obj.category.trim() === '') {
      errors['category'] = ['Category is required'];
    } else if (!this.minWordsCheck(obj.category, 1, 3)) {
      errors['category'] = ['Category must contain between 1 and 3 words'];
    }

    // technique
    if (typeof obj.technique !== 'string' || obj.technique.trim() === '') {
      errors['technique'] = ['Technique is required'];
    } else if (!this.minWordsCheck(obj.technique, 2, 3)) {
      errors['technique'] = ['Technique must contain between 2 and 3 words'];
    }

    // description
    if (typeof obj.description !== 'string' || obj.description.trim() === '') {
      errors['description'] = ['Description is required'];
    } else if (!this.minWordsCheck(obj.description, 2, 3)) {
      errors['description'] = ['Description must contain between 2 and 3 words'];
    }

    // definition
    if (typeof obj.definition !== 'string' || obj.definition.trim() === '') {
      errors['definition'] = ['Definition is required'];
    } else if (!this.minWordsCheck(obj.definition, 2, 3)) {
      errors['definition'] = ['Definition must contain between 2 and 3 words'];
    }

    // references
    if (!Array.isArray(obj.references)) {
      errors['references'] = ['References must be an array'];
    } else {
      const invalidRefs = obj.references.filter((ref: any) => typeof ref !== 'string' || !this.isValidURL(ref));
      if (invalidRefs.length > 0) {
        errors['references'] = ['All references must be valid URLs'];
      }
    }


    return {
      isValid: Object.keys(errors).length === 0,
      errors
    };
  }
}

