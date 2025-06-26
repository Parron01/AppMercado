import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideZoneChangeDetection } from '@angular/core';
import { AppComponent } from './app/app.component';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';
import { inject } from '@angular/core';
import { AuthService } from './app/services/auth.service';
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';

const API_BASE_URL = 'http://localhost:8080'; // Variável global para a URL base da API

bootstrapApplication(AppComponent, {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideHttpClient(
      withInterceptors([
        (req, next) => {
          // Verifica se a URL é relativa antes de adicionar a URL base
          const isRelativeUrl = !req.url.startsWith('http://') && !req.url.startsWith('https://');
          const apiReq = isRelativeUrl
            ? req.clone({ url: `${API_BASE_URL}${req.url}` })
            : req;
          // Adiciona o token se existir
          const authService = inject(AuthService);
          const token = authService.getToken();
          const authReq = token
            ? apiReq.clone({ setHeaders: { Authorization: `Bearer ${token}` } })
            : apiReq;
          return next(authReq).pipe(
            catchError((err) => {
              if (err.status === 401) {
                authService.logout();
              }
              return throwError(() => err);
            })
          );
        }
      ])
    )
  ]
}).catch((err) => console.error(err));
