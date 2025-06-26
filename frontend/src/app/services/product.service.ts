import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

interface CreateProductDTO {
  name: string;
}

interface ProductResponseDTO {
  id: number;
  name: string;
  barcode?: string;
  createdAt: string;
  updatedAt: string;
}

@Injectable({
  providedIn: 'root'
})
export class ProductService {
  private readonly apiUrl = '/products'; // URL relativa para usar o interceptor

  constructor(private http: HttpClient) {}

  getAllProducts(): Observable<ProductResponseDTO[]> {
    return this.http.get<ProductResponseDTO[]>(`${this.apiUrl}/all`);
  }

  createProduct(product: CreateProductDTO): Observable<ProductResponseDTO> {
    return this.http.post<{ product: ProductResponseDTO }>(`${this.apiUrl}/create`, product)
      .pipe(
        // Extrai o objeto product do response
        // @ts-ignore
        map(response => response.product)
      );
  }

  deleteProduct(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/delete/${id}`);
  }
}
