import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {SERVER_API_URL} from '../../../app.constants';

export enum TfaMethod {
  EMAIL = 'EMAIL',
  TOTP = 'TOTP'
}

export interface TfaInitRequest {
  method: TfaMethod;
}

export interface TfaInitResponse {
  method: string;
  delivery: {
    qrBase64?: string;
    code?: string;
    expiresInSeconds: number;
  };
}

export interface TfaVerifyRequest {
  method: TfaMethod;
  code: string;
}

export interface TfaVerifyResponse {
  valid: boolean;
  expired: boolean;
  message?: string;
  expiresInSeconds?: number;
}

@Injectable({
  providedIn: 'root'
})
export class TfaService {
  private readonly baseUrl = `${SERVER_API_URL}api/tfa`;

  constructor(private http: HttpClient) {}

  initTfa(request: TfaInitRequest): Observable<TfaInitResponse> {
    return this.http.post<TfaInitResponse>(`${this.baseUrl}/init`, request);
  }

  verifyTfa(request: TfaVerifyRequest): Observable<TfaVerifyResponse> {
    return this.http.post<TfaVerifyResponse>(`${this.baseUrl}/verify`, request);
  }
}
