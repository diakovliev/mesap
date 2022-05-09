import { RegisterService, RegisteredUserData, RegisterData } from './../register.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {

  public data: RegisterData = {
    login: "Test user 1",
    mail: "test@mail.com",
    password: "1234",
  }

  constructor(private _register: RegisterService) { }

  ngOnInit(): void {
  }

  registerUser() {
    this._register.registerUser(this.data)
      .subscribe(registeredUser => console.log(`Registered user salt: ${registeredUser.salt} verifier: ${registeredUser.verifier}`))
  }

}
