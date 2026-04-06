export interface ProductImage {
  thumbnail: string;
  mobile: string;
  tablet: string;
  desktop: string;
}

export interface Product {
  id: number;
  name: string;
  description?: string;
  price: number;
  category: string[];
  image: ProductImage;
  stock: number;
}

export interface OrderItemReq {
  productId: number;
  quantity: number;
}

export interface OrderReq {
  couponCode?: string;
  items: OrderItemReq[];
}

export interface OrderItem {
  productId: number;
  quantity: number;
}

export interface OrderResponse {
  id: string;
  status: string;
  items: OrderItem[];
  products: Product[];
  couponValid: boolean;
  totalPrice: number;
  discountedPrice: number;
}

export interface ValidateCouponReq {
  couponCode: string;
}

export interface ValidateCouponResponse {
  isValid: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}
