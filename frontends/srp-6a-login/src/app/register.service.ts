import { Buffer } from 'buffer';
import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParamsOptions, HttpResponse } from '@angular/common/http';
import { Observable, from, map, tap, of, switchMap, catchError, throwError, zip, defer } from 'rxjs';
import { Injectable } from '@angular/core';

import { SRP, SrpClient } from 'fast-srp-hap'

export interface IRegisterData {
  login: string
  password: string
  mail: string
}

export interface IRegisterRequest {
  login: string
  salt: string
  verifier: string
}

export interface IRegisteredUserData {
  UserId: number
}

export interface ILoginRequestData {
  Login: string
  Secret1: string
}

export interface ILoginResponseData {
  Server: string
  Secret2: string
}

export interface ILogin2RequestData {
  Server: string
  Secret3: string
}

export interface ILogin2ResponseData {
	Secret4: string
}

export interface ILoginResult {
  SessionId: string
}


@Injectable({
  providedIn: 'root'
})
export class RegisterService {

  SALT_SIZE = 0
  KEY_SIZE = 32
  ENCODING: BufferEncoding = 'base64'

  private _client?: SrpClient

  API_ROOT = "/api/auth"
  HTTP_OPTIONS = {
    headers: new HttpHeaders({ "Content-Type": "application/json" }),
  }
  SRP_PARAMS = SRP.params[4096]

  constructor(private _http: HttpClient) { }

  private handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong.
      console.error(
        `Backend returned code ${error.status}, body was: `, error.error);
    }
    // Return an observable with a user-facing error message.
    return throwError(() => new Error('Something bad happened; please try again later.'));
  }

  private newSalt(): Observable<Buffer> {
    if(this.SALT_SIZE == 0) {
      return of(Buffer.alloc(this.SALT_SIZE))
    } else {
      return defer(() => from(SRP.genKey(this.SALT_SIZE)));
    }
  }

  private newKey(): Observable<[Buffer, Buffer]> {
    const s = this.newSalt()
    const k = defer(() => from(SRP.genKey(this.KEY_SIZE)));
    return zip(s, k)
  }

  private newClient(data: IRegisterData): Observable<Buffer> {
    return this.newKey().pipe(
      tap(([salt, a]) => console.log("[login use] salt: " + salt.toString(this.ENCODING) + " a: " + a.toString(this.ENCODING))),
      map(([salt, a]) => {
        this._client = new SrpClient(this.SRP_PARAMS, salt, Buffer.from(data.login), Buffer.from(data.password), a, false)
        return this._client.computeA()
      }),
    )
  }

  registerUser(data: IRegisterData): Observable<IRegisteredUserData> {

    console.log("[registerUser] called")

    return this.newSalt().pipe(
        tap(b => console.log("[register use] salt: " + b.toString(this.ENCODING))),
        map(salt => [ salt, SRP.computeVerifier(this.SRP_PARAMS, salt, Buffer.from(data.login), Buffer.from(data.password)) ] ),
        map(([s, v]) => ({ login: Buffer.from(data.login).toString(this.ENCODING), salt: s.toString(this.ENCODING), verifier: v.toString(this.ENCODING) } as IRegisterRequest) ),
        tap(r => console.log("[register request] " + JSON.stringify(r))),
        catchError(this.handleError),
        switchMap(request => this._http.post<IRegisteredUserData>(`${this.API_ROOT}/register`, request, { responseType: 'json' })),
      )

  }

  loginUser(data: IRegisterData): Observable<ILoginResult> {

    console.log("[loginUser] called")

    return this.newClient(data).pipe(
      map(A => ({ Login: Buffer.from(data.login).toString(this.ENCODING), Secret1: Buffer.from(A).toString(this.ENCODING) } as ILoginRequestData)),
      tap(r => console.log("[login request] " + JSON.stringify(r))),
      catchError(this.handleError),
      switchMap(request => this._http.post<ILoginResponseData>(`${this.API_ROOT}/login`, request, { responseType: 'json' })),
      map(response => {
        this._client!.setB(Buffer.from(response.Secret2, this.ENCODING))
        return { Server: response.Server, Secret3: Buffer.from(this._client!.computeM1()).toString(this.ENCODING) } as ILogin2RequestData
      }),
      switchMap(request => this._http.post<ILogin2ResponseData>(`${this.API_ROOT}/login2`, request, { responseType: 'json' })),
      map(response => {
        this._client!.checkM2(Buffer.from(response.Secret4, this.ENCODING))
        return { SessionId: Buffer.from(this._client!.computeK()).toString(this.ENCODING) } as ILoginResult
      }),
    )
  }
}
