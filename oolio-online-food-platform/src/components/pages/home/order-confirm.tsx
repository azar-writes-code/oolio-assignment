import { Button } from "@/components/ui/button";
import {
  DialogClose,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { CheckCircle } from "lucide-react";
import { useCartStore, selectCartTotal } from "@/store/useCartStore";

export function OrderConfirm() {
  const orderResponse = useCartStore((state) => state.orderResponse);
  const items = useCartStore((state) => state.items);
  const clearCart = useCartStore((state) => state.clearCart);
  const resetOrder = useCartStore((state) => state.resetOrder);
  const total = useCartStore(selectCartTotal);

  const handleStartNewOrder = () => {
    clearCart();
    resetOrder();
  };

  // If we have an order response, use it. Otherwise fall back to current cart items (for preview)
  const displayItems = orderResponse ? orderResponse.products.map(p => ({
    ...p,
    quantity: orderResponse.items.find(i => i.productId === p.id)?.quantity || 0
  })) : items;

  const displayTotal = orderResponse ? orderResponse.discountedPrice : total;

  return (
    <div className="bg-white p-8 rounded-3xl w-full max-w-md">
      <DialogHeader className="items-start text-left mb-6">
        <CheckCircle className="size-10 text-green-600 mb-4" strokeWidth={2.5} />
        <DialogTitle className="text-4xl font-bold text-font-primary-80 leading-tight">
          Order Confirmed
        </DialogTitle>
        <DialogDescription className="text-font-primary-50 font-medium text-base">
          We hope you enjoy your food!
        </DialogDescription>
      </DialogHeader>

      <div className="bg-font-secondary-10 rounded-xl p-6 mb-8 space-y-4">
        <div className="max-h-[300px] overflow-y-auto pr-2 custom-scrollbar space-y-4">
          {displayItems.map((item) => (
            <div key={item.id} className="space-y-4">
              <div className="flex justify-between items-center gap-4">
                <div className="flex gap-4 items-center">
                  <img
                    src={item.image.thumbnail}
                    alt={item.name}
                    className="size-12 rounded-lg object-cover"
                  />
                  <div className="space-y-1">
                    <p className="text-sm text-font-primary-80 font-bold truncate max-w-[150px]">
                      {item.name}
                    </p>
                    <p className="flex items-center gap-3">
                      <span className="text-[#C66C50] font-bold">{item.quantity}x</span>
                      <span className="text-font-primary-10 font-medium text-xs">
                        @ ${item.price.toFixed(2)}
                      </span>
                    </p>
                  </div>
                </div>
                <p className="text-font-primary-80 font-bold text-base">
                  ${(item.price * item.quantity).toFixed(2)}
                </p>
              </div>
              <Separator className="bg-gray-200/50" />
            </div>
          ))}
        </div>

        <div className="flex justify-between items-center pt-2">
          <p className="text-font-primary-80 text-sm font-medium">
            Order Total
          </p>
          <p className="text-font-primary-80 text-2xl font-bold">
            ${displayTotal.toFixed(2)}
          </p>
        </div>
      </div>

      <DialogClose asChild>
        <Button
          onClick={handleStartNewOrder}
          className="bg-secondary-button hover:bg-secondary-button/90 w-full py-7 rounded-full text-base font-bold text-white transition-all shadow-md active:scale-[0.98]"
        >
          Start New Order
        </Button>
      </DialogClose>
    </div>
  );
}
