'use client';

import { useQuery } from '@apollo/client/react/hooks';
import { GET_PRODUCTS } from '@/graphql/queries';
import { ProductCard } from '@/components/products/ProductCard';
import { Hero } from '@/components/Hero';

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  category: string;
  inventory: {
    quantity: number;
    inStock: boolean;
  };
}

export default function HomePage() {
  const { data, loading, error } = useQuery(GET_PRODUCTS, {
    variables: { limit: 8 }
  });

  if (error) {
    console.error('GraphQL Error:', error);
  }

  return (
    <div>
      {/* Hero Section */}
      <Hero />

      {/* Featured Products */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">Featured Products</h2>
          <p className="text-lg text-gray-600">Discover our most popular items</p>
        </div>

        {loading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {Array.from({ length: 8 }).map((_, index) => (
              <div key={index} className="bg-white rounded-lg shadow-sm p-4 animate-pulse">
                <div className="bg-gray-200 w-full h-48 rounded-md mb-4"></div>
                <div className="bg-gray-200 h-4 rounded mb-2"></div>
                <div className="bg-gray-200 h-4 rounded w-3/4 mb-2"></div>
                <div className="bg-gray-200 h-4 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-12">
            <p className="text-gray-500 mb-4">
              Unable to load products. Using demo data for now.
            </p>
            <DemoProducts />
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {data?.products?.map((product: Product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </section>

      {/* Categories Section */}
      <section className="bg-gray-100 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">Shop by Category</h2>
            <p className="text-lg text-gray-600">Find exactly what you're looking for</p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {categories.map((category) => (
              <div
                key={category.name}
                className="bg-white rounded-lg shadow-sm p-6 text-center hover:shadow-md transition-shadow cursor-pointer"
              >
                <div className="text-4xl mb-4">{category.icon}</div>
                <h3 className="font-medium text-gray-900">{category.name}</h3>
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  );
}

// Demo products for when GraphQL is not available
function DemoProducts() {
  const demoProducts = [
    {
      id: '1',
      name: 'Premium Wireless Headphones',
      description: 'High-quality sound with noise cancellation',
      price: 299.99,
      image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=300&fit=crop',
      category: 'Electronics',
      inventory: { quantity: 10, inStock: true }
    },
    {
      id: '2',
      name: 'Organic Cotton T-Shirt',
      description: 'Comfortable and sustainable fashion',
      price: 49.99,
      image: 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=300&h=300&fit=crop',
      category: 'Clothing',
      inventory: { quantity: 25, inStock: true }
    },
    {
      id: '3',
      name: 'Smart Watch Pro',
      description: 'Track your fitness and stay connected',
      price: 399.99,
      image: 'https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=300&h=300&fit=crop',
      category: 'Electronics',
      inventory: { quantity: 5, inStock: true }
    },
    {
      id: '4',
      name: 'Leather Backpack',
      description: 'Stylish and durable for everyday use',
      price: 129.99,
      image: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=300&h=300&fit=crop',
      category: 'Accessories',
      inventory: { quantity: 15, inStock: true }
    }
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      {demoProducts.map((product) => (
        <ProductCard key={product.id} product={product} />
      ))}
    </div>
  );
}

const categories = [
  { name: 'Electronics', icon: '📱' },
  { name: 'Clothing', icon: '👕' },
  { name: 'Home & Garden', icon: '🏠' },
  { name: 'Sports', icon: '⚽' },
];
