import { RegisterComponent } from './register/register.component';
import { LoginComponent } from './login/login.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  { path:"login", component: LoginComponent },
  { path:"register", component: RegisterComponent },
  //{ path: '', redirectTo: "login", pathMatch: "full" },
  { path: '', redirectTo: "login", pathMatch: "full" },
  //{ path: '**', redirectTo: "login", pathMatch: "full" },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
