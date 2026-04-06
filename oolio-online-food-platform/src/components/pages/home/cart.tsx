import { Button } from "@/components/ui/button";
import { Dialog, DialogTrigger, DialogContent } from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { CircleX, Trees, Loader2 } from "lucide-react";
import { OrderConfirm } from "./order-confirm";
import {
  useCartStore,
  selectCartTotal,
  selectCartCount,
} from "@/store/useCartStore";
import { useShallow } from "zustand/react/shallow";

export function Cart() {
  const {
    items,
    removeFromCart,
    confirmOrder,
    isOrdering,
    couponCode,
    setCouponCode,
    validateCouponAction,
    isCouponValid,
    isValidatingCoupon,
    couponMessage,
    total,
    count,
  } = useCartStore(
    useShallow((state) => ({
      items: state.items,
      removeFromCart: state.removeFromCart,
      confirmOrder: state.confirmOrder,
      isOrdering: state.isOrdering,
      couponCode: state.couponCode,
      setCouponCode: state.setCouponCode,
      validateCouponAction: state.validateCouponAction,
      isCouponValid: state.isCouponValid,
      isValidatingCoupon: state.isValidatingCoupon,
      couponMessage: state.couponMessage,
      total: selectCartTotal(state),
      count: selectCartCount(state),
    })),
  );

  if (items.length === 0) {
    return <EmptyCart />;
  }

  const isConfirmDisabled =
    isOrdering || (couponCode !== "" && isCouponValid === null);

  return (
    <div className="w-full bg-white h-fit px-6 py-10 rounded-2xl shadow-sm">
      <h1 className="text-2xl font-bold mb-6 text-[#C66C50]">
        Your cart ({count})
      </h1>
      <div className="space-y-6">
        <div className="flex flex-col gap-4 max-h-[60vh] overflow-y-auto pr-2 custom-scrollbar">
          {items.map((item) => (
            <div key={item.id} className="space-y-4">
              <div className="flex justify-between items-center gap-2">
                <div className="space-y-1">
                  <p className="text-sm text-font-primary-80 font-bold truncate max-w-[150px] sm:max-w-none">
                    {item.name}
                  </p>
                  <p className="flex items-center gap-3">
                    <span className="text-[#C66C50] font-bold">
                      {item.quantity}x
                    </span>
                    <span className="text-font-primary-10 font-medium text-xs">
                      @ ${item.price.toFixed(2)}
                    </span>
                    <span className="text-font-primary-50 font-bold text-xs">
                      ${(item.price * item.quantity).toFixed(2)}
                    </span>
                  </p>
                </div>
                <Button
                  size="icon"
                  variant="ghost"
                  onClick={() => removeFromCart(item.id)}
                  className="hover:bg-red-50 group shrink-0"
                >
                  <CircleX
                    className="size-5 transition-colors group-hover:text-red-600"
                    strokeWidth={1.5}
                    color="#B8A1A2"
                  />
                </Button>
              </div>
              <Separator className="bg-gray-100" />
            </div>
          ))}
        </div>

        <div className="pt-2 space-y-6">
          {/* Coupon Section */}
          <div className="space-y-3">
            <p className="text-sm font-bold text-font-primary-80">
              Have a coupon?
            </p>
            <div className="flex gap-2">
              <input
                type="text"
                value={couponCode}
                onChange={(e) => setCouponCode(e.target.value)}
                placeholder="Enter code"
                className="flex-1 px-4 py-2 rounded-lg border border-gray-200 text-sm focus:outline-none focus:border-[#C66C50] transition-colors"
              />
              <Button
                variant="outline"
                size="sm"
                onClick={() => validateCouponAction()}
                disabled={isValidatingCoupon || !couponCode}
                className="border-secondary-button text-secondary-button hover:bg-secondary-button hover:text-white px-4 h-10 font-bold transition-all"
              >
                {isValidatingCoupon ? (
                  <Loader2 className="size-4 animate-spin" />
                ) : (
                  "Verify"
                )}
              </Button>
            </div>
            {couponMessage && (
              <p
                className={`text-xs font-bold leading-relaxed ${isCouponValid ? "text-green-600" : "text-amber-600"}`}
              >
                {couponMessage}
              </p>
            )}
          </div>

          <div className="flex justify-between items-center">
            <p className="text-font-primary-80 text-sm font-medium">
              Order Total
            </p>
            <p className="text-font-primary-80 text-2xl font-bold">
              ${total.toFixed(2)}
            </p>
          </div>

          <div className="flex w-full mx-auto bg-font-secondary-10 p-4 justify-center space-x-2 items-center rounded-xl">
            <Trees className="size-5 text-green-700" />
            <p className="text-font-primary-80 text-xs sm:text-sm font-medium">
              This is a <span className="font-bold">carbon-neutral</span>{" "}
              delivery
            </p>
          </div>

          <Dialog>
            <DialogTrigger asChild>
              <Button
                disabled={isConfirmDisabled}
                className="bg-secondary-button hover:bg-secondary-button/90 w-full py-7 rounded-full text-base font-bold text-white transition-all shadow-md active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={() => confirmOrder()}
              >
                {isOrdering ? (
                  <>
                    <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                    Confirming...
                  </>
                ) : (
                  "Confirm Order"
                )}
              </Button>
            </DialogTrigger>
            <DialogContent
              showCloseButton={false}
              className="w-[90vw] md:max-w-md p-0 overflow-hidden rounded-3xl border-none"
            >
              <OrderConfirm />
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  );
}

function EmptyCart() {
  return (
    <div className="w-full bg-white h-fit px-6 py-10 rounded-2xl shadow-sm">
      <h1 className="text-2xl font-bold mb-6 text-[#C66C50]">Your cart (0)</h1>
      <div className="flex flex-col items-center justify-center py-8">
        <img
          className="w-32 h-32 object-contain mb-4 opacity-80"
          src="empty-cart.png"
          alt="empty cart"
        />
        <p className="text-sm text-font-primary-50 font-bold">
          Your added items will appear here
        </p>
      </div>
    </div>
  );
}
