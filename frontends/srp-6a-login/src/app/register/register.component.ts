import { RegisterService, IRegisteredUserData, IRegisterData } from './../register.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {

  public data: IRegisterData = {
    login: "Test user 1",
    mail: "test@mail.com",
    password: "1234",
  }

  constructor(private _register: RegisterService) { }

  ngOnInit(): void {
  }

  registerUser() {
    this._register.registerUser(this.data)
      .subscribe(registeredUser => console.log(`Registered user id: ${registeredUser.UserId}`))
  }

  loginUser() {
    this._register.loginUser(this.data)
      .subscribe(loginData => console.log(`Login data: ${JSON.stringify(loginData)}`))
  }

}
