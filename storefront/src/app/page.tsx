'use client';

import Link from 'next/link';
import { useState } from 'react';
import { useQuery } from '@apollo/client';
import { ProductCard } from '@/components/ProductCard';
import { GET_PRODUCTS } from '@/graphql/queries';
import type { Product } from '@/graphql/queries';

export default function HomePage() {
  const [selectedCategory, setSelectedCategory] = useState('all');
  
  const { loading, error, data } = useQuery(GET_PRODUCTS, {
    variables: {
      filter: {
        limit: 8
      }
    }
  });

  const categories = [
    { id: 'all', name: 'All Categories' },
    { id: 'electronics', name: 'Electronics' },
    { id: 'fashion', name: 'Fashion' },
    { id: 'home', name: 'Home & Kitchen' },
    { id: 'beauty', name: 'Beauty & Personal Care' }
  ];

  const collections = [
    {
      id: "electronics",
      name: "Electronics",
      description: "The latest gadgets and tech innovations",
      image: "https://via.placeholder.com/600x400/3B82F6/FFFFFF?text=Electronics",
      itemCount: 120
    },
    {
      id: "fashion",
      name: "Fashion",
      description: "Trending styles for every occasion",
      image: "https://via.placeholder.com/600x400/EC4899/FFFFFF?text=Fashion",
      itemCount: 250
    },
    {
      id: "home",
      name: "Home & Kitchen",
      description: "Elevate your living space",
      image: "https://via.placeholder.com/600x400/F97316/FFFFFF?text=Home",
      itemCount: 180
    }
  ];

  // Use real data if available, otherwise fall back to mock data
  const featuredProducts = data?.products || [
    {
      id: "1",
      name: "Premium Wireless Headphones",
      price: 299.99,
      rating: 4.8,
      reviews: 432,
      category: "electronics",
      description: "High-quality noise-cancelling wireless headphones with premium sound.",
      image: "https://via.placeholder.com/600x600/3B82F6/FFFFFF?text=Headphones",
      badges: ["New", "Top Rated"]
    },
    {
      id: "2", 
      name: "Smart Fitness Watch",
      price: 199.99,
      rating: 4.6,
      reviews: 251,
      category: "electronics",
      description: "Track your fitness goals with this advanced smartwatch.",
      image: "https://via.placeholder.com/600x600/10B981/FFFFFF?text=Watch",
      badges: ["Bestseller"]
    },
    {
      id: "3",
      name: "Portable Bluetooth Speaker",
      price: 79.99,
      rating: 4.5,
      reviews: 198,
      category: "electronics",
      description: "Powerful sound in a compact, portable design.",
      image: "https://via.placeholder.com/600x600/F59E0B/FFFFFF?text=Speaker",
      badges: ["Sale"]
    },
    {
      id: "4",
      name: "Leather Weekend Bag",
      price: 159.99,
      rating: 4.7,
      reviews: 86,
      category: "fashion",
      description: "Premium leather travel bag, perfect for weekend getaways.",
      image: "https://via.placeholder.com/600x600/6366F1/FFFFFF?text=Bag",
      badges: ["Handmade"]
    },
    {
      id: "5",
      name: "Smart Home Hub",
      price: 129.99,
      rating: 4.4,
      reviews: 175,
      category: "electronics",
      description: "Control all your smart home devices from one central hub.",
      image: "https://via.placeholder.com/600x600/8B5CF6/FFFFFF?text=Hub",
      badges: []
    },
    {
      id: "6",
      name: "Organic Cotton Sheets",
      price: 89.99,
      rating: 4.9,
      reviews: 324,
      category: "home",
      description: "Ultra-soft, 100% organic cotton sheets for the perfect night's sleep.",
      image: "https://via.placeholder.com/600x600/EC4899/FFFFFF?text=Sheets",
      badges: ["Eco-friendly"]
    },
    {
      id: "7",
      name: "Premium Coffee Maker",
      price: 149.99,
      rating: 4.7,
      reviews: 213,
      category: "home",
      description: "Barista-quality coffee from the comfort of your own home.",
      image: "https://via.placeholder.com/600x600/F97316/FFFFFF?text=Coffee",
      badges: ["Trending"]
    },
    {
      id: "8",
      name: "Natural Skincare Set",
      price: 69.99,
      rating: 4.8,
      reviews: 142,
      category: "beauty",
      description: "All-natural skincare products made with organic ingredients.",
      image: "https://via.placeholder.com/600x600/14B8A6/FFFFFF?text=Skincare",
      badges: ["Organic"]
    }
  ];

  const displayedProducts = selectedCategory === 'all' 
    ? featuredProducts 
    : featuredProducts.filter((product: any) => product.category === selectedCategory);

  return (
    <div className="min-h-screen">
      {/* Hero Banner */}
      <section className="relative">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-90"></div>
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32 lg:py-40 flex flex-col md:flex-row items-center">
          <div className="text-center md:text-left md:w-1/2 mb-8 md:mb-0">
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-white mb-6 leading-tight">
              Welcome to Unified Commerce OS
            </h1>
            <p className="text-lg md:text-xl text-white/90 mb-8 max-w-xl">
              Your complete unified commerce platform. Shop the latest trends in electronics, fashion, home goods and more. Free shipping on orders over $50.
            </p>
            <div className="flex flex-col sm:flex-row justify-center md:justify-start gap-4">
              <Link 
                href="/products" 
                className="px-8 py-3 rounded-lg bg-white text-blue-600 hover:bg-blue-50 transition font-semibold text-lg"
              >
                Shop Now
              </Link>
              <Link 
                href="/deals" 
                className="px-8 py-3 rounded-lg bg-transparent border-2 border-white text-white hover:bg-white/10 transition font-semibold text-lg"
              >
                View Deals
              </Link>
            </div>
          </div>
          <div className="md:w-1/2">
            <div className="relative">
              <div className="absolute -top-6 -left-6 w-32 h-32 rounded-full bg-yellow-400 opacity-50 animate-pulse"></div>
              <div className="absolute -bottom-6 -right-6 w-32 h-32 rounded-full bg-pink-400 opacity-50 animate-pulse delay-300"></div>
              <img 
                src="https://via.placeholder.com/800x600/FFFFFF/333333?text=Premium+Products" 
                alt="Premium Products" 
                className="relative z-10 rounded-lg shadow-xl w-full"
              />
            </div>
          </div>
        </div>
      </section>

      {/* Categories Selector */}
      <section className="bg-white py-8 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex overflow-x-auto gap-4 pb-2 scrollbar-hide">
            {categories.map(category => (
              <button
                key={category.id}
                className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                  selectedCategory === category.id
                    ? 'bg-blue-600 text-white shadow-md'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setSelectedCategory(category.id)}
              >
                {category.name}
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Featured Products */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Featured Products
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              Discover our handpicked selection of premium products, curated for quality and customer satisfaction.
            </p>
          </div>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {displayedProducts.map((product: any) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
          
          <div className="text-center mt-12">
            <Link href="/products" className="inline-flex items-center px-6 py-3 border border-gray-300 shadow-sm text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 transition">
              View All Products
              <svg className="ml-2 -mr-1 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
              </svg>
            </Link>
          </div>
        </div>
      </section>
      
      {/* Collections */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Shop Collections
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              Explore our curated collections and find exactly what you're looking for.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {collections.map((collection) => (
              <Link href={`/collections/${collection.id}`} key={collection.id} className="group">
                <div className="relative rounded-lg overflow-hidden shadow-md h-80">
                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent z-10"></div>
                  <img 
                    src={collection.image} 
                    alt={collection.name}
                    className="absolute inset-0 w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-500"
                  />
                  <div className="absolute inset-0 z-20 flex flex-col justify-end p-6">
                    <h3 className="text-white text-2xl font-bold mb-2">
                      {collection.name}
                    </h3>
                    <p className="text-white/80 mb-4">
                      {collection.description}
                    </p>
                    <span className="text-sm text-white/60">
                      {collection.itemCount} products
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>
      
      {/* Features Section */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-white p-6 rounded-lg shadow-sm text-center">
              <div className="inline-flex items-center justify-center h-12 w-12 rounded-md bg-blue-100 text-blue-600 mb-4">
                <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Fast Delivery</h3>
              <p className="text-gray-600">Free shipping on orders over $50. Get your items delivered within 2-3 business days.</p>
            </div>
            
            <div className="bg-white p-6 rounded-lg shadow-sm text-center">
              <div className="inline-flex items-center justify-center h-12 w-12 rounded-md bg-blue-100 text-blue-600 mb-4">
                <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Secure Payments</h3>
              <p className="text-gray-600">All transactions are secure and encrypted. We accept all major credit cards and digital wallets.</p>
            </div>
            
            <div className="bg-white p-6 rounded-lg shadow-sm text-center">
              <div className="inline-flex items-center justify-center h-12 w-12 rounded-md bg-blue-100 text-blue-600 mb-4">
                <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Easy Returns</h3>
              <p className="text-gray-600">Not satisfied? Return items within 30 days for a full refund, no questions asked.</p>
            </div>
          </div>
        </div>
      </section>
      
      {/* Testimonials Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              What Our Customers Say
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              Don't just take our word for it — hear from some of our satisfied customers!
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* Testimonial 1 */}
            <div className="bg-gray-50 p-6 rounded-lg shadow-sm relative">
              <div className="absolute top-6 left-6 text-gray-300">
                <svg className="h-12 w-12 transform -translate-x-1/4 -translate-y-1/3" fill="currentColor" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M10 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6zM22 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6z"></path>
                </svg>
              </div>
              <div className="relative">
                <p className="text-gray-700 mb-6">
                  "The quality of products I received exceeded my expectations. The customer service team was incredibly helpful when I had questions. I'll definitely be shopping here again!"
                </p>
                <div className="flex items-center">
                  <div className="h-12 w-12 rounded-full bg-gray-300 overflow-hidden">
                    <img src="https://via.placeholder.com/200x200/F59E0B/FFFFFF?text=S" alt="Sarah J." className="h-full w-full object-cover" />
                  </div>
                  <div className="ml-4">
                    <p className="text-gray-900 font-semibold">Sarah Johnson</p>
                    <p className="text-gray-500 text-sm">Verified Customer</p>
                  </div>
                </div>
                <div className="flex mt-2">
                  {[...Array(5)].map((_, i) => (
                    <svg 
                      key={i} 
                      className="h-5 w-5 text-yellow-400"
                      fill="currentColor" 
                      viewBox="0 0 20 20" 
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  ))}
                </div>
              </div>
            </div>
            
            {/* Testimonial 2 */}
            <div className="bg-gray-50 p-6 rounded-lg shadow-sm relative">
              <div className="absolute top-6 left-6 text-gray-300">
                <svg className="h-12 w-12 transform -translate-x-1/4 -translate-y-1/3" fill="currentColor" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M10 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6zM22 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6z"></path>
                </svg>
              </div>
              <div className="relative">
                <p className="text-gray-700 mb-6">
                  "Fast delivery and the products were exactly as described. The website made it easy to find exactly what I was looking for. Highly recommend!"
                </p>
                <div className="flex items-center">
                  <div className="h-12 w-12 rounded-full bg-gray-300 overflow-hidden">
                    <img src="https://via.placeholder.com/200x200/3B82F6/FFFFFF?text=M" alt="Michael T." className="h-full w-full object-cover" />
                  </div>
                  <div className="ml-4">
                    <p className="text-gray-900 font-semibold">Michael Thomas</p>
                    <p className="text-gray-500 text-sm">Verified Customer</p>
                  </div>
                </div>
                <div className="flex mt-2">
                  {[...Array(5)].map((_, i) => (
                    <svg 
                      key={i} 
                      className="h-5 w-5 text-yellow-400"
                      fill="currentColor" 
                      viewBox="0 0 20 20" 
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  ))}
                </div>
              </div>
            </div>
            
            {/* Testimonial 3 */}
            <div className="bg-gray-50 p-6 rounded-lg shadow-sm relative">
              <div className="absolute top-6 left-6 text-gray-300">
                <svg className="h-12 w-12 transform -translate-x-1/4 -translate-y-1/3" fill="currentColor" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M10 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6zM22 8c-4.4 0-8 3.6-8 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 14c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6z"></path>
                </svg>
              </div>
              <div className="relative">
                <p className="text-gray-700 mb-6">
                  "I was skeptical about ordering online, but the return policy gave me confidence. In the end, I loved everything I ordered! The quality is outstanding."
                </p>
                <div className="flex items-center">
                  <div className="h-12 w-12 rounded-full bg-gray-300 overflow-hidden">
                    <img src="https://via.placeholder.com/200x200/EC4899/FFFFFF?text=A" alt="Alicia R." className="h-full w-full object-cover" />
                  </div>
                  <div className="ml-4">
                    <p className="text-gray-900 font-semibold">Alicia Rodriguez</p>
                    <p className="text-gray-500 text-sm">Verified Customer</p>
                  </div>
                </div>
                <div className="flex mt-2">
                  {[...Array(5)].map((_, i) => (
                    <svg 
                      key={i} 
                      className={`h-5 w-5 ${i < 4 ? 'text-yellow-400' : 'text-gray-300'}`}
                      fill="currentColor" 
                      viewBox="0 0 20 20" 
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
      
      {/* Trending Products Section */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <span className="text-blue-600 font-semibold uppercase tracking-wide">Hot right now</span>
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mt-2 mb-4">
              Trending Products
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              The most popular items that everyone's talking about this season.
            </p>
          </div>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
            <div className="bg-white rounded-lg shadow-sm overflow-hidden group">
              <div className="relative h-48 md:h-64">
                <img 
                  src="https://via.placeholder.com/400x400/F59E0B/FFFFFF?text=Trending+1" 
                  alt="Trending Product 1" 
                  className="h-full w-full object-cover group-hover:opacity-90 transition"
                />
                <div className="absolute top-2 right-2 bg-white rounded-full p-1 shadow-md">
                  <svg className="h-5 w-5 text-gray-500 hover:text-red-500 cursor-pointer" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
                  </svg>
                </div>
              </div>
              <div className="p-4">
                <h3 className="font-medium text-gray-900 mb-1">Smart Watch Pro</h3>
                <div className="flex items-center">
                  <span className="text-lg font-bold text-gray-900 mr-2">$129.99</span>
                  <span className="text-sm line-through text-gray-500">$159.99</span>
                </div>
              </div>
            </div>
            
            <div className="bg-white rounded-lg shadow-sm overflow-hidden group">
              <div className="relative h-48 md:h-64">
                <img 
                  src="https://via.placeholder.com/400x400/3B82F6/FFFFFF?text=Trending+2" 
                  alt="Trending Product 2" 
                  className="h-full w-full object-cover group-hover:opacity-90 transition"
                />
                <div className="absolute top-2 right-2 bg-white rounded-full p-1 shadow-md">
                  <svg className="h-5 w-5 text-gray-500 hover:text-red-500 cursor-pointer" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
                  </svg>
                </div>
              </div>
              <div className="p-4">
                <h3 className="font-medium text-gray-900 mb-1">Wireless Earbuds</h3>
                <div className="flex items-center">
                  <span className="text-lg font-bold text-gray-900 mr-2">$89.99</span>
                  <span className="text-sm line-through text-gray-500">$109.99</span>
                </div>
              </div>
            </div>
            
            <div className="bg-white rounded-lg shadow-sm overflow-hidden group">
              <div className="relative h-48 md:h-64">
                <img 
                  src="https://via.placeholder.com/400x400/10B981/FFFFFF?text=Trending+3" 
                  alt="Trending Product 3" 
                  className="h-full w-full object-cover group-hover:opacity-90 transition"
                />
                <div className="absolute top-2 right-2 bg-white rounded-full p-1 shadow-md">
                  <svg className="h-5 w-5 text-gray-500 hover:text-red-500 cursor-pointer" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
                  </svg>
                </div>
              </div>
              <div className="p-4">
                <h3 className="font-medium text-gray-900 mb-1">Smart Home Camera</h3>
                <div className="flex items-center">
                  <span className="text-lg font-bold text-gray-900 mr-2">$79.99</span>
                  <span className="text-sm line-through text-gray-500">$99.99</span>
                </div>
              </div>
            </div>
            
            <div className="bg-white rounded-lg shadow-sm overflow-hidden group">
              <div className="relative h-48 md:h-64">
                <img 
                  src="https://via.placeholder.com/400x400/6366F1/FFFFFF?text=Trending+4" 
                  alt="Trending Product 4" 
                  className="h-full w-full object-cover group-hover:opacity-90 transition"
                />
                <div className="absolute top-2 right-2 bg-white rounded-full p-1 shadow-md">
                  <svg className="h-5 w-5 text-gray-500 hover:text-red-500 cursor-pointer" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
                  </svg>
                </div>
              </div>
              <div className="p-4">
                <h3 className="font-medium text-gray-900 mb-1">Portable Power Bank</h3>
                <div className="flex items-center">
                  <span className="text-lg font-bold text-gray-900 mr-2">$49.99</span>
                  <span className="text-sm line-through text-gray-500">$69.99</span>
                </div>
              </div>
            </div>
          </div>
          
          <div className="text-center mt-8">
            <Link href="/trending" className="inline-flex items-center text-blue-600 hover:text-blue-800 font-medium">
              View all trending products
              <svg className="ml-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
              </svg>
            </Link>
          </div>
        </div>
      </section>
      
      {/* Brands Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Trusted by Top Brands
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              We partner with the best brands to bring you quality products you can trust.
            </p>
          </div>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 items-center">
            {[1, 2, 3, 4, 5, 6, 7, 8].map((brand) => (
              <div key={brand} className="flex justify-center">
                <img 
                  src={`https://via.placeholder.com/180x60/E5E7EB/6B7280?text=Brand+${brand}`} 
                  alt={`Brand Partner ${brand}`} 
                  className="h-12 object-contain grayscale hover:grayscale-0 opacity-70 hover:opacity-100 transition-all"
                />
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Newsletter */}
      <section className="py-16 bg-blue-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold text-white mb-4">
            Join Our Newsletter
          </h2>
          <p className="text-lg text-white/80 mb-8 max-w-2xl mx-auto">
            Subscribe to get special offers, free giveaways, and once-in-a-lifetime deals.
          </p>
          <form className="max-w-md mx-auto flex flex-col sm:flex-row">
            <input 
              type="email" 
              placeholder="Your email address" 
              className="w-full sm:flex-1 px-4 py-3 rounded-lg sm:rounded-r-none mb-3 sm:mb-0 focus:outline-none focus:ring-2 focus:ring-blue-300"
              required
            />
            <button 
              type="submit" 
              className="w-full sm:w-auto bg-white text-blue-600 px-6 py-3 sm:rounded-l-none rounded-lg font-medium hover:bg-blue-50 transition"
            >
              Subscribe
            </button>
          </form>
          <p className="text-white/70 text-sm mt-4">
            By subscribing you agree to our Terms of Service and Privacy Policy.
          </p>
        </div>
      </section>

      {/* Download App CTA */}
      <section className="py-16 bg-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row items-center justify-between">
            <div className="md:w-1/2 mb-8 md:mb-0">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                Download Our Mobile App
              </h2>
              <p className="text-lg text-gray-600 mb-6">
                Shop on the go and get exclusive mobile-only offers. Our app makes shopping even easier!
              </p>
              <div className="flex flex-wrap gap-4">
                <a href="#" className="flex items-center bg-black text-white px-4 py-3 rounded-lg hover:bg-gray-800 transition">
                  <svg className="h-6 w-6 mr-2" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path d="M17.5649 12.3664C17.5306 9.52926 19.8648 8.10656 19.9785 8.03646C18.5491 6.02199 16.4118 5.82959 15.6376 5.81145C13.8137 5.63718 12.0474 6.92667 11.1211 6.92667C10.1948 6.92667 8.74964 5.82959 7.20677 5.86588C5.24035 5.90216 3.41409 7.03099 2.42961 8.77455C0.396768 12.3301 1.92964 17.5777 3.87792 20.3602C4.85953 21.7283 6.01472 23.2568 7.51563 23.1842C8.97882 23.1116 9.54285 22.2409 11.3122 22.2409C13.0816 22.2409 13.6093 23.1842 15.1284 23.1479C16.6839 23.1116 17.6837 21.765 18.629 20.3784C19.7478 18.8135 20.2027 17.2849 20.221 17.2123C20.1866 17.1942 17.6019 16.127 17.5649 12.3664Z"/>
                    <path d="M14.6023 3.6975C15.4146 2.69834 15.9602 1.34061 15.8102 0C14.6387 0.0514888 13.2094 0.774079 12.3607 1.75508C11.6047 2.62392 10.9491 4.01795 11.118 5.32196C12.4391 5.41018 13.7537 4.68313 14.6023 3.6975Z"/>
                  </svg>
                  <div>
                    <div className="text-xs">Download on the</div>
                    <div className="text-lg font-semibold font-sans">App Store</div>
                  </div>
                </a>
                <a href="#" className="flex items-center bg-black text-white px-4 py-3 rounded-lg hover:bg-gray-800 transition">
                  <svg className="h-6 w-6 mr-2" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path d="M2.86567 0.5C2.26682 0.5 1.73682 0.64 1.27567 0.92C0.81567 1.20667 0.464017 1.58667 0.224017 2.06C-0.0159834 2.53333 -0.0693167 3.06 0.0773833 3.64C0.213883 4.21333 0.484017 4.71333 0.884017 5.14C1.28402 5.56667 1.76068 5.85333 2.31401 6C2.86068 6.13333 3.38068 6.1 3.87401 5.9C4.36068 5.7 4.77401 5.38 5.11401 4.94C5.45401 4.5 5.62401 4.01333 5.62401 3.48C5.62401 2.76667 5.37401 2.15333 4.87401 1.64C4.36735 1.12667 3.69901 0.873333 2.86567 0.5Z" />
                    <path d="M0.0127692 23.5V7.12H5.76277V23.5H0.0127692Z" />
                    <path d="M17.6752 7.12C19.8219 7.12 21.4952 7.73333 22.6952 8.96C23.9085 10.1867 24.5152 11.9533 24.5152 14.26V23.5H18.7652V15.14C18.7652 14.22 18.5685 13.54 18.1752 13.1C17.7819 12.66 17.1885 12.44 16.3952 12.44C15.6019 12.44 14.9952 12.66 14.5752 13.1C14.1685 13.54 13.9652 14.22 13.9652 15.14V23.5H8.21519V7.12H13.9652V9.02C14.4819 8.34 15.1685 7.82 16.0219 7.46C16.8752 7.10667 17.8085 6.93 18.8219 7.12H17.6752Z" />
                  </svg>
                  <div>
                    <div className="text-xs">GET IT ON</div>
                    <div className="text-lg font-semibold font-sans">Google Play</div>
                  </div>
                </a>
              </div>
            </div>
            <div className="md:w-1/3">
              <img 
                src="https://via.placeholder.com/500x600/3B82F6/FFFFFF?text=Mobile+App" 
                alt="Mobile App" 
                className="w-full max-w-xs mx-auto rounded-3xl shadow-xl"
              />
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}