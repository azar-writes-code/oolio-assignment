import api from '@/lib/axios';
import type { Product, OrderReq, OrderResponse, ValidateCouponReq, ValidateCouponResponse, PaginatedResponse } from '@/types/api';

export const getProducts = async (): Promise<PaginatedResponse<Product>> => {
  const { data } = await api.get('/products/');
  return data;
};

export const placeOrder = async (order: OrderReq): Promise<OrderResponse> => {
  const { data } = await api.post('/order/', order);
  return data;
};

export const validateCoupon = async (req: ValidateCouponReq): Promise<ValidateCouponResponse> => {
  const { data } = await api.post('/order/validate-coupon', req);
  return data;
};
