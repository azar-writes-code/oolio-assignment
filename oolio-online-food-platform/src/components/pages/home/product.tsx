import { Button } from "@/components/ui/button";
import { ShoppingCart, Plus, Minus } from "lucide-react";
import type { Product as ProductType } from "@/types/api";
import { useCartStore } from "@/store/useCartStore";

interface ProductProps {
  product: ProductType;
}

export function Product({ product }: ProductProps) {
  const addToCart = useCartStore((state) => state.addToCart);
  const updateQuantity = useCartStore((state) => state.updateQuantity);
  const cartItem = useCartStore((state) =>
    state.items.find((item) => item.id === product.id)
  );

  const handleAddToCart = () => {
    addToCart(product);
  };

  return (
    <div className="flex flex-col w-full">
      <div className="relative group">
        <img
          className={`w-full aspect-square object-cover rounded-xl border-2 transition-all duration-300 ${
            cartItem ? "border-secondary-button" : "border-transparent"
          }`}
          src={product.image.desktop}
          alt={product.name}
          onError={(e)=>(e.currentTarget.src = product.image.tablet)}
        />
        {cartItem ? (
          <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 w-[85%] sm:w-2/3 flex items-center justify-between border border-secondary-button bg-secondary-button text-white rounded-full py-3 px-4 shadow-lg">
            <button
              onClick={() => updateQuantity(product.id, cartItem.quantity - 1)}
              className="p-1 hover:bg-white/20 rounded-full border border-white/50 transition-colors"
            >
              <Minus className="size-4" />
            </button>
            <span className="font-bold text-sm sm:text-base">{cartItem.quantity}</span>
            <button
              onClick={() => updateQuantity(product.id, cartItem.quantity + 1)}
              className="p-1 hover:bg-white/20 rounded-full border border-white/50 transition-colors"
            >
              <Plus className="size-4" />
            </button>
          </div>
        ) : (
          <Button
            variant="default"
            size="lg"
            onClick={handleAddToCart}
            className="absolute -bottom-6 left-1/2 -translate-x-1/2 w-[85%] sm:w-2/3 border border-[#B8A1A2] bg-white hover:bg-font-secondary-10 hover:border-secondary-button text-font-primary-80 rounded-full py-6 flex items-center justify-center gap-2 shadow-sm transition-all"
          >
            <ShoppingCart color="#C66C50" strokeWidth={2} className="size-5" />
            <span className="font-bold text-sm">Add to cart</span>
          </Button>
        )}
      </div>
      <div className="flex flex-col gap-1 mt-10">
        <h6 className="text-xs text-font-primary-10 font-medium">
          {product.category[0]}
        </h6>
        <p className="text-base text-font-primary-80 font-bold leading-tight truncate">
          {product.name}
        </p>
        <p className="text-base text-secondary-button font-bold">
          ${product.price.toFixed(2)}
        </p>
      </div>
    </div>
  );
}
