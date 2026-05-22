import { create } from 'zustand';
import type { Product, OrderReq, OrderResponse } from '@/types/api';
import { placeOrder, validateCoupon } from '@/services/api';

interface CartItem extends Product {
  quantity: number;
}

interface CartState {
  items: CartItem[];
  orderResponse: OrderResponse | null;
  isOrdering: boolean;
  error: string | null;
  // Coupon state
  couponCode: string;
  isCouponValid: boolean | null;
  isValidatingCoupon: boolean;
  couponMessage: string | null;
  
  addToCart: (product: Product) => void;
  removeFromCart: (productId: number) => void;
  updateQuantity: (productId: number, quantity: number) => void;
  clearCart: () => void;
  confirmOrder: () => Promise<void>;
  resetOrder: () => void;
  // Coupon actions
  setCouponCode: (code: string) => void;
  validateCouponAction: () => Promise<void>;
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  orderResponse: null,
  isOrdering: false,
  error: null,
  // Initial coupon state
  couponCode: '',
  isCouponValid: null,
  isValidatingCoupon: false,
  couponMessage: null,

  addToCart: (product) => {
    const { items } = get();
    const existingItem = items.find((item) => item.id === product.id);

    if (existingItem) {
      set({
        items: items.map((item) =>
          item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        ),
      });
    } else {
      set({ items: [...items, { ...product, quantity: 1 }] });
    }
  },

  removeFromCart: (productId) => {
    set({ items: get().items.filter((item) => item.id !== productId) });
  },

  updateQuantity: (productId, quantity) => {
    if (quantity <= 0) {
      get().removeFromCart(productId);
      return;
    }
    set({
      items: get().items.map((item) =>
        item.id === productId ? { ...item, quantity } : item
      ),
    });
  },

  clearCart: () => {
    set({ items: [], couponCode: '', isCouponValid: null, couponMessage: null });
  },

  setCouponCode: (code) => {
    set({ couponCode: code, isCouponValid: null, couponMessage: null });
  },

  validateCouponAction: async () => {
    const { couponCode } = get();
    if (!couponCode) return;

    set({ isValidatingCoupon: true, couponMessage: null });
    try {
      const response = await validateCoupon({ couponCode });
      set({ 
        isCouponValid: response.isValid, 
        couponMessage: response.isValid ? 'Coupon applied successfully!' : 'Invalid coupon code. Order will proceed without discount.',
        isValidatingCoupon: false 
      });
    } catch (err: any) {
      set({ 
        isCouponValid: false, 
        couponMessage: 'Error validating coupon. Proceeding without discount.',
        isValidatingCoupon: false 
      });
    }
  },

  confirmOrder: async () => {
    const { items, couponCode, isCouponValid } = get();
    if (items.length === 0) return;

    set({ isOrdering: true, error: null });
    try {
      const orderReq: OrderReq = {
        couponCode: isCouponValid ? couponCode : undefined,
        items: items.map((item) => ({
          productId: item.id,
          quantity: item.quantity,
        })),
      };
      const response = await placeOrder(orderReq);
      set({ orderResponse: response, isOrdering: false });
    } catch (err: any) {
      set({ error: err.message || 'Failed to place order', isOrdering: false });
    }
  },

  resetOrder: () => {
    set({ orderResponse: null, error: null, couponCode: '', isCouponValid: null, couponMessage: null });
  },
}));

export const selectCartTotal = (state: CartState) =>
  state.items.reduce((total, item) => total + item.price * item.quantity, 0);

export const selectCartCount = (state: CartState) =>
  state.items.reduce((count, item) => count + item.quantity, 0);
