import { Routes } from '@angular/router';
import { LayoutComponent } from './components/layout/layout.component';
import { HomeComponent } from './components/home/home.component';
import { LoginComponent } from './components/login/login.component';
import { inject } from '@angular/core';
import { AuthService } from './services/auth.service';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  {
    path: '',
    canActivate: [() => {
      const auth = inject(AuthService);
      if (!auth.isLoggedIn()) {
        auth.logout();
        return false;
      }
      return true;
    }],
    component: LayoutComponent,
    children: [
      { path: '', component: HomeComponent },
      { path: 'historico', component: HomeComponent },
      { path: 'configuracoes', component: HomeComponent },
      { path: 'perfil', component: HomeComponent }
    ]
  }
];
