import { Product } from "./product";
import { useProducts } from "@/hooks/useProducts";

export function ProductList() {
  const { data: products, isLoading, error } = useProducts();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 flex-2">
        {Array.from({ length: 6 }).map((_, index) => (
          <div key={index} className="animate-pulse bg-gray-100 rounded-lg h-64 w-full" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500 font-semibold p-4 border border-red-200 rounded-lg bg-red-50">
        Error loading products: {(error as Error).message}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-y-12 gap-x-6">
      {products?.data?.map((product) => (
        <Product key={product.id} product={product} />
      ))}
    </div>
  );
}
