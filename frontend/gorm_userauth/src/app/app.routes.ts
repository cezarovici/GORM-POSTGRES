import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { MembersComponent } from './members/members.component';
import { HomeComponent } from './home/home.component';

export const routes: Routes = [
    { path: 'login', component: LoginComponent },
    { path: 'members', component: MembersComponent },
    { path: '', component: HomeComponent },
    { path: '**', redirectTo: '' }
];
