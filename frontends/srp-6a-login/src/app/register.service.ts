import { HttpClient } from '@angular/common/http';
import { Observable, from, map, tap, of } from 'rxjs';
import { Injectable } from '@angular/core';

import { SRP } from 'fast-srp-hap'

export interface RegisterData {
  login: string,
  mail: string,
  password: string,
}

export interface RegisteredUserData {
  login: string,
  salt: string,
  verifier: string,
}

@Injectable({
  providedIn: 'root'
})
export class RegisterService {

  SALT_SIZE = 32
  ENCODING = 'hex'
  //ENCODING = 'base64'

  API_ROOT = "/api"

  constructor(private _http: HttpClient) { }

  registerUser(data: RegisterData): Observable<RegisteredUserData> {

    console.log("[registerUser] called")

    const init = this.SALT_SIZE == 0 ? of(Buffer.alloc(this.SALT_SIZE)) : from(SRP.genKey(this.SALT_SIZE))

    return init.pipe(
      tap(b => console.log("[use salt] " + b.toString(this.ENCODING))),
      map(salt => {
        var v = SRP.computeVerifier(SRP.params[4096], salt, Buffer.from(data.login), Buffer.from(data.password))
        return { login: data.login, salt: salt.toString(this.ENCODING), verifier: v.toString(this.ENCODING) }
      })
    )
  }

}
