'use client';

import { ProductCard } from '@/components/ProductCard';

export default function DealsPage() {
  // Mock deal products data
  const dealProducts = [
    {
      id: '1',
      name: 'Premium Wireless Headphones',
      price: 89.99,
      originalPrice: 129.99,
      discount: 30,
      description: 'High-quality noise-cancelling wireless headphones with premium sound.',
      image: 'https://via.placeholder.com/600x600/3B82F6/FFFFFF?text=Headphones',
      rating: 4.8,
      reviews: 432,
    },
    {
      id: '2',
      name: 'Smart Fitness Watch',
      price: 129.99,
      originalPrice: 199.99,
      discount: 35,
      description: 'Track your fitness goals with this advanced smartwatch.',
      image: 'https://via.placeholder.com/600x600/10B981/FFFFFF?text=Watch',
      rating: 4.6,
      reviews: 251,
    },
    {
      id: '3',
      name: 'Portable Bluetooth Speaker',
      price: 45.99,
      originalPrice: 69.99,
      discount: 34,
      description: 'Powerful sound in a compact, portable design.',
      image: 'https://via.placeholder.com/600x600/F59E0B/FFFFFF?text=Speaker',
      rating: 4.5,
      reviews: 198,
    },
    {
      id: '4',
      name: 'Leather Weekend Bag',
      price: 159.99,
      originalPrice: 249.99,
      discount: 36,
      description: 'Premium leather travel bag, perfect for weekend getaways.',
      image: 'https://via.placeholder.com/600x600/6366F1/FFFFFF?text=Bag',
      rating: 4.7,
      reviews: 86,
    },
    {
      id: '5',
      name: 'Smart Home Hub',
      price: 89.99,
      originalPrice: 129.99,
      discount: 31,
      description: 'Control all your smart home devices from one central hub.',
      image: 'https://via.placeholder.com/600x600/8B5CF6/FFFFFF?text=Hub',
      rating: 4.4,
      reviews: 175,
    },
    {
      id: '6',
      name: 'Organic Cotton Sheets',
      price: 69.99,
      originalPrice: 99.99,
      discount: 30,
      description: 'Ultra-soft, 100% organic cotton sheets for the perfect night\'s sleep.',
      image: 'https://via.placeholder.com/600x600/EC4899/FFFFFF?text=Sheets',
      rating: 4.9,
      reviews: 324,
    },
    {
      id: '7',
      name: 'Premium Coffee Maker',
      price: 119.99,
      originalPrice: 179.99,
      discount: 33,
      description: 'Barista-quality coffee from the comfort of your own home.',
      image: 'https://via.placeholder.com/600x600/F97316/FFFFFF?text=Coffee',
      rating: 4.7,
      reviews: 213,
    },
    {
      id: '8',
      name: 'Natural Skincare Set',
      price: 49.99,
      originalPrice: 79.99,
      discount: 38,
      description: 'All-natural skincare products made with organic ingredients.',
      image: 'https://via.placeholder.com/600x600/14B8A6/FFFFFF?text=Skincare',
      rating: 4.8,
      reviews: 142,
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Banner */}
      <section className="bg-gradient-to-r from-red-600 to-orange-500 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
            Exclusive Deals & Discounts
          </h1>
          <p className="text-xl text-white/90 max-w-3xl mx-auto">
            Save big on our premium products with limited-time offers. Don't miss out!
          </p>
        </div>
      </section>

      {/* Countdown Timer */}
      <section className="py-8 bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Deal Ends In</h2>
            <div className="flex justify-center space-x-4">
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">02</span>
                <span className="text-sm text-gray-600">Days</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">14</span>
                <span className="text-sm text-gray-600">Hours</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">36</span>
                <span className="text-sm text-gray-600">Minutes</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">22</span>
                <span className="text-sm text-gray-600">Seconds</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Deals Grid */}
      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {dealProducts.map((product) => (
              <div key={product.id} className="relative">
                {product.discount && (
                  <div className="absolute top-2 left-2 bg-red-500 text-white text-sm font-bold px-2 py-1 rounded z-10">
                    {product.discount}% OFF
                  </div>
                )}
                <ProductCard product={product} />
              </div>
            ))}
          </div>
          
          {/* Pagination */}
          <div className="mt-12 flex justify-center">
            <nav className="flex items-center space-x-2">
              <button className="px-4 py-2 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200">
                Previous
              </button>
              <button className="px-4 py-2 rounded-md bg-blue-600 text-white">
                1
              </button>
              <button className="px-4 py-2 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200">
                2
              </button>
              <button className="px-4 py-2 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200">
                3
              </button>
              <button className="px-4 py-2 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200">
                Next
              </button>
            </nav>
          </div>
        </div>
      </section>

      {/* Additional Deals Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              More Great Deals
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              Discover additional savings on our most popular products.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg p-8 text-white">
              <h3 className="text-2xl font-bold mb-4">Free Shipping</h3>
              <p className="mb-6">On all orders over $50. Limited time offer!</p>
              <button className="bg-white text-blue-600 px-6 py-3 rounded-md font-medium hover:bg-gray-100 transition">
                Shop Now
              </button>
            </div>
            
            <div className="bg-gradient-to-r from-green-500 to-teal-600 rounded-lg p-8 text-white">
              <h3 className="text-2xl font-bold mb-4">Buy 2, Get 1 Free</h3>
              <p className="mb-6">On select fashion items. Mix and match your favorites.</p>
              <button className="bg-white text-green-600 px-6 py-3 rounded-md font-medium hover:bg-gray-100 transition">
                Shop Now
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}