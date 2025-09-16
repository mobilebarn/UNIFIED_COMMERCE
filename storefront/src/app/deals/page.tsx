'use client';

import { useQuery } from '@apollo/client';
import { useState, useEffect } from 'react';
import { ProductCard } from '@/components/ProductCard';
import { GET_ACTIVE_PROMOTIONS, GET_CAMPAIGNS, GET_PRODUCTS } from '@/graphql/queries';
import type { Promotion, Campaign, Product } from '@/graphql/queries';

export default function DealsPage() {
  const [timeLeft, setTimeLeft] = useState({
    days: 2,
    hours: 14,
    minutes: 36,
    seconds: 22
  });

  // Fetch active promotions
  const { loading: promotionsLoading, error: promotionsError, data: promotionsData } = useQuery(GET_ACTIVE_PROMOTIONS, {
    variables: { merchantId: 'default-merchant-id' } // In a real app, this would come from context or props
  });

  // Fetch campaigns
  const { loading: campaignsLoading, error: campaignsError, data: campaignsData } = useQuery(GET_CAMPAIGNS, {
    variables: { filter: { merchantId: 'default-merchant-id' } }
  });

  // Fetch products on sale (for demo purposes, we'll fetch all products)
  const { loading: productsLoading, error: productsError, data: productsData } = useQuery(GET_PRODUCTS, {
    variables: { filter: { search: 'sale' } } // This would be adjusted based on actual promotion logic
  });

  // Countdown timer effect
  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(prev => {
        const newSeconds = prev.seconds - 1;
        if (newSeconds >= 0) {
          return { ...prev, seconds: newSeconds };
        }
        
        const newMinutes = prev.minutes - 1;
        if (newMinutes >= 0) {
          return { ...prev, minutes: newMinutes, seconds: 59 };
        }
        
        const newHours = prev.hours - 1;
        if (newHours >= 0) {
          return { ...prev, hours: newHours, minutes: 59, seconds: 59 };
        }
        
        const newDays = prev.days - 1;
        if (newDays >= 0) {
          return { days: newDays, hours: 23, minutes: 59, seconds: 59 };
        }
        
        // Reset timer when it reaches zero
        return { days: 2, hours: 14, minutes: 36, seconds: 22 };
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  // Format time values to always have two digits
  const formatTime = (value: number) => value.toString().padStart(2, '0');

  if (promotionsLoading || campaignsLoading || productsLoading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-3xl font-bold text-gray-900 mb-8">Loading Deals...</h1>
          </div>
        </div>
      </div>
    );
  }

  if (promotionsError || campaignsError || productsError) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-3xl font-bold text-red-600 mb-8">Error Loading Deals</h1>
            <p className="text-gray-600">
              {promotionsError?.message || campaignsError?.message || productsError?.message}
            </p>
          </div>
        </div>
      </div>
    );
  }

  // Extract data from GraphQL responses
  const promotions: Promotion[] = promotionsData?.activePromotions || [];
  const campaigns: Campaign[] = campaignsData?.campaigns || [];
  const products: Product[] = productsData?.products || [];

  // For demo purposes, we'll create mock deal products with real data structure
  // In a real implementation, these would be derived from actual promotions
  const dealProducts = products.map((product, index) => ({
    id: product.id,
    name: product.title,
    price: product.price,
    originalPrice: product.price * 1.3, // Mock original price
    discount: Math.floor(Math.random() * 40) + 10, // Mock discount percentage
    description: product.description || 'Premium product with great features.',
    image: product.imageUrl || `https://via.placeholder.com/600x600/3B82F6/FFFFFF?text=Product+${index + 1}`,
    rating: (Math.random() * 2 + 3).toFixed(1), // Mock rating between 3-5
    reviews: Math.floor(Math.random() * 500) + 50, // Mock review count
  }));

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
                <span className="text-3xl font-bold text-red-600">{formatTime(timeLeft.days)}</span>
                <span className="text-sm text-gray-600">Days</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">{formatTime(timeLeft.hours)}</span>
                <span className="text-sm text-gray-600">Hours</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">{formatTime(timeLeft.minutes)}</span>
                <span className="text-sm text-gray-600">Minutes</span>
              </div>
              <div className="flex flex-col items-center">
                <span className="text-3xl font-bold text-red-600">{formatTime(timeLeft.seconds)}</span>
                <span className="text-sm text-gray-600">Seconds</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Active Promotions Section */}
      {promotions.length > 0 && (
        <section className="py-12 bg-white">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-12">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                Current Promotions
              </h2>
              <p className="text-lg text-gray-600 max-w-3xl mx-auto">
                Take advantage of our active deals and save on your favorite products.
              </p>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {promotions.map((promotion) => (
                <div key={promotion.id} className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
                  <div className="flex justify-between items-start mb-4">
                    <h3 className="text-xl font-bold text-gray-900">{promotion.name}</h3>
                    {promotion.discountValue && (
                      <span className="bg-red-100 text-red-800 text-sm font-bold px-2 py-1 rounded">
                        {promotion.discountType === 'PERCENTAGE' ? `${promotion.discountValue}% OFF` : `$${promotion.discountValue} OFF`}
                      </span>
                    )}
                  </div>
                  <p className="text-gray-600 mb-4">{promotion.description}</p>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-500">
                      Valid: {new Date(promotion.startDate).toLocaleDateString()} - {promotion.endDate ? new Date(promotion.endDate).toLocaleDateString() : 'Ongoing'}
                    </span>
                    <button className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 transition">
                      Shop Now
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

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

      {/* Campaigns Section */}
      {campaigns.length > 0 && (
        <section className="py-16 bg-gray-100">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-12">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                Special Campaigns
              </h2>
              <p className="text-lg text-gray-600 max-w-3xl mx-auto">
                Discover our featured campaigns with exclusive offers and deals.
              </p>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {campaigns.map((campaign) => (
                <div key={campaign.id} className="bg-white rounded-lg p-8 shadow-sm">
                  <h3 className="text-2xl font-bold text-gray-900 mb-4">{campaign.name}</h3>
                  <p className="text-gray-600 mb-6">{campaign.description}</p>
                  <div className="flex justify-between items-center mb-6">
                    <div>
                      <span className="text-sm text-gray-500">Campaign Period:</span>
                      <p className="font-medium">
                        {new Date(campaign.startDate).toLocaleDateString()} - {campaign.endDate ? new Date(campaign.endDate).toLocaleDateString() : 'Ongoing'}
                      </p>
                    </div>
                    {campaign.budget && (
                      <div>
                        <span className="text-sm text-gray-500">Budget:</span>
                        <p className="font-medium">${campaign.budget.toLocaleString()}</p>
                      </div>
                    )}
                  </div>
                  <button className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-md font-medium hover:opacity-90 transition">
                    Explore Campaign
                  </button>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

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