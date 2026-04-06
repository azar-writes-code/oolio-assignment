import { Cart, ProductList } from "@/components/pages";

export default function Home() {
  return (
    <main className="min-h-screen bg-font-secondary-10 py-8 px-4 sm:px-8 lg:px-12 xl:px-24">
      <div className="container mx-auto flex flex-col lg:flex-row gap-8 items-start">
        <div className="flex-1 w-full">
          <h1 className="text-4xl font-bold text-font-primary-80 mb-8">Desserts</h1>
          <ProductList />
        </div>
        <aside className="w-full lg:w-[380px] xl:w-[420px] lg:sticky lg:top-8">
          <Cart />
        </aside>
      </div>
    </main>
  );
}
