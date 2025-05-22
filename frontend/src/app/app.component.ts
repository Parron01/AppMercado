import { Component } from '@angular/core';
// import { RouterOutlet } from '@angular/router'; // Removido se não usar rotas no template principal

@Component({
  selector: 'app-root',
  imports: [], // Removido se não usar rotas no template principal
  standalone: true, // Adicionado para componentes standalone sem módulos dedicados
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = ''; 
}