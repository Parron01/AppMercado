import { Component } from '@angular/core';
import { RouterOutlet, RouterModule } from '@angular/router';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [RouterOutlet, RouterModule],
  template: `
    <div class="layout">
      <aside class="sidebar">
        <h2 class="logo">AppMercado</h2>
        <nav class="nav">
          <a routerLink="/" routerLinkActive="active" class="nav-item">🛒 Produtos</a>
          <a routerLink="/historico" routerLinkActive="active" class="nav-item">📜 Histórico</a>
          <a routerLink="/configuracoes" routerLinkActive="active" class="nav-item">⚙️ Configurações</a>
          <a routerLink="/perfil" routerLinkActive="active" class="nav-item">👤 Perfil</a>
        </nav>
      </aside>
      <main class="content">
        <router-outlet></router-outlet>
      </main>
    </div>
  `,
  styles: `
    .layout {
      display: flex;
      height: 100vh;
      font-family: 'Segoe UI', Arial, sans-serif;
    }
    .sidebar {
      width: 250px;
      background-color: #2d7a4b;
      color: #fff;
      display: flex;
      flex-direction: column;
      padding: 1rem;
    }
    .logo {
      font-size: 1.5rem;
      font-weight: bold;
      text-align: center;
      margin-bottom: 2rem;
    }
    .nav {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
    .nav-item {
      text-decoration: none;
      color: #fff;
      font-size: 1.1rem;
      padding: 0.5rem 1rem;
      border-radius: 6px;
      transition: background 0.3s;
    }
    .nav-item:hover {
      background-color: #25613a;
    }
    .nav-item.active {
      background-color: #1e4d2e;
    }
    .content {
      flex: 1;
      padding: 2rem;
      background-color: #f5f5f5;
      overflow-y: auto;
    }
  `
})
export class LayoutComponent {}
