'use client';

import { ProductCard } from '@/components/ProductCard';

export default function SavedItems() {
  // Mock saved items data
  const savedItems = [
    {
      id: '1',
      name: 'Premium Wireless Headphones',
      price: 89.99,
      description: 'High-quality noise-cancelling wireless headphones with premium sound.',
      image: 'https://via.placeholder.com/600x600/3B82F6/FFFFFF?text=Headphones',
    },
    {
      id: '2',
      name: 'Smart Fitness Watch',
      price: 129.99,
      description: 'Track your fitness goals with this advanced smartwatch.',
      image: 'https://via.placeholder.com/600x600/10B981/FFFFFF?text=Watch',
    },
    {
      id: '3',
      name: 'Portable Bluetooth Speaker',
      price: 45.99,
      description: 'Powerful sound in a compact, portable design.',
      image: 'https://via.placeholder.com/600x600/F59E0B/FFFFFF?text=Speaker',
    },
    {
      id: '4',
      name: 'Leather Weekend Bag',
      price: 159.99,
      description: 'Premium leather travel bag, perfect for weekend getaways.',
      image: 'https://via.placeholder.com/600x600/6366F1/FFFFFF?text=Bag',
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Saved Items</h1>
      
      {savedItems.length === 0 ? (
        <div className="text-center py-12">
          <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No saved items</h3>
          <p className="mt-1 text-sm text-gray-500">Save items you like to view them here later.</p>
          <div className="mt-6">
            <button className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
              Browse Products
            </button>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          {savedItems.map((item) => (
            <ProductCard key={item.id} product={item} />
          ))}
        </div>
      )}
    </div>
  );
}