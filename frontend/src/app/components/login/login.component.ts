import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="login-container">
      <h1>Login</h1>
      <form (ngSubmit)="onSubmit()" #loginForm="ngForm" autocomplete="off">
        <div class="form-group">
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            [(ngModel)]="loginData.email"
            required
            placeholder="Digite seu email"
          />
        </div>
        <div class="form-group">
          <label for="password">Senha</label>
          <input
            type="password"
            id="password"
            name="password"
            [(ngModel)]="loginData.password"
            required
            placeholder="Digite sua senha"
          />
        </div>
        <button type="submit" [disabled]="loginForm.invalid">Entrar</button>
        <p *ngIf="errorMessage" class="error-message">{{ errorMessage }}</p>
      </form>
    </div>
  `,
  styles: `
    .login-container {
      max-width: 400px;
      margin: 2rem auto;
      padding: 2rem;
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
      font-family: 'Segoe UI', Arial, sans-serif;
    }
    h1 {
      text-align: center;
      color: #2d7a4b;
      margin-bottom: 1.5rem;
    }
    .form-group {
      margin-bottom: 1rem;
    }
    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: bold;
    }
    input {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid #ccc;
      border-radius: 4px;
      font-size: 1rem;
    }
    button {
      width: 100%;
      padding: 0.7rem;
      background: #2d7a4b;
      color: #fff;
      border: none;
      border-radius: 4px;
      font-size: 1rem;
      cursor: pointer;
      transition: background 0.3s;
    }
    button:disabled {
      background: #ccc;
      cursor: not-allowed;
    }
    .error-message {
      color: #e53935;
      margin-top: 1rem;
      text-align: center;
    }
  `
})
export class LoginComponent {
  loginData = { email: '', password: '' };
  errorMessage: string | null = null;

  constructor(private http: HttpClient, private router: Router, private authService: AuthService) {}

  onSubmit() {
    console.log('Tentando realizar login com os dados:', this.loginData);

    this.http.post<{ token: string; user: any }>('/auth/login', this.loginData).subscribe({
      next: (response) => {
        console.log('Login bem-sucedido! Dados retornados:', response);
        this.authService.setAuth(response.token, response.user);
        this.router.navigate(['/']);
      },
      error: (err) => {
        console.error('Erro ao realizar login:', err);
        this.errorMessage = err.error?.error || 'Erro ao realizar login. Tente novamente.';
      }
    });
  }
}
