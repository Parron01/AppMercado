import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ProductService } from '../../services/product.service';

interface Product {
  id: number;
  name: string;
}

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  providers: [ProductService]
})
export class HomeComponent implements OnInit {
  products: Product[] = [];
  newProduct: Product = { id: 0, name: '' };

  constructor(private productService: ProductService) {}

  ngOnInit() {
    this.loadProducts();
  }

  loadProducts() {
    this.productService.getAllProducts().subscribe({
      next: (data) => {
        this.products = data.map((product) => ({
          id: product.id,
          name: product.name
        }));
      },
      error: (err) => console.error('Erro ao carregar produtos:', err)
    });
  }

  addProduct() {
    if (!this.newProduct.name.trim()) {
      return;
    }

    const productToCreate = {
      name: this.newProduct.name.trim()
    };

    this.productService.createProduct(productToCreate).subscribe({
      next: (createdProduct) => {
        this.products.push({
          id: createdProduct.id,
          name: createdProduct.name
        });
        this.newProduct = { id: 0, name: '' };
      },
      error: (err) => console.error('Erro ao adicionar produto:', err)
    });
  }

  removeProduct(id: number) {
    if (!confirm('Tem certeza que deseja remover este produto da lista?')) {
      return;
    }
    this.productService.deleteProduct(id).subscribe({
      next: () => {
        this.products = this.products.filter((p) => p.id !== id);
      },
      error: (err) => console.error('Erro ao remover produto:', err)
    });
  }
}