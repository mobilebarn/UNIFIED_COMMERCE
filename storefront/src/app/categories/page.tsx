'use client';

import Link from 'next/link';
import { useQuery } from '@apollo/client';
import { GET_CATEGORIES } from '@/graphql/queries';
import type { Category } from '@/graphql/queries';

export default function CategoriesPage() {
  // For now, we'll keep the static categories but in a real implementation,
  // we would fetch them from the GraphQL API
  const categories = [
    {
      id: 'electronics',
      name: 'Electronics',
      description: 'The latest gadgets and tech innovations',
      image: 'https://via.placeholder.com/600x400/3B82F6/FFFFFF?text=Electronics',
      productCount: 120,
    },
    {
      id: 'fashion',
      name: 'Fashion',
      description: 'Trending styles for every occasion',
      image: 'https://via.placeholder.com/600x400/EC4899/FFFFFF?text=Fashion',
      productCount: 250,
    },
    {
      id: 'home',
      name: 'Home & Kitchen',
      description: 'Elevate your living space',
      image: 'https://via.placeholder.com/600x400/F97316/FFFFFF?text=Home',
      productCount: 180,
    },
    {
      id: 'beauty',
      name: 'Beauty & Personal Care',
      description: 'Premium beauty and personal care products',
      image: 'https://via.placeholder.com/600x400/14B8A6/FFFFFF?text=Beauty',
      productCount: 95,
    },
    {
      id: 'sports',
      name: 'Sports & Outdoors',
      description: 'Gear up for your next adventure',
      image: 'https://via.placeholder.com/600x400/8B5CF6/FFFFFF?text=Sports',
      productCount: 75,
    },
    {
      id: 'books',
      name: 'Books & Media',
      description: 'Expand your knowledge and entertainment',
      image: 'https://via.placeholder.com/600x400/F59E0B/FFFFFF?text=Books',
      productCount: 210,
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Banner */}
      <section className="bg-gradient-to-r from-blue-600 to-purple-600 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
            Shop by Category
          </h1>
          <p className="text-xl text-white/90 max-w-3xl mx-auto">
            Discover our wide range of product categories, carefully curated for quality and style.
          </p>
        </div>
      </section>

      {/* Categories Grid */}
      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
            {categories.map((category) => (
              <Link 
                key={category.id} 
                href={`/categories/${category.id}`}
                className="group bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
              >
                <div className="relative h-48">
                  <img 
                    src={category.image} 
                    alt={category.name}
                    className="w-full h-full object-cover"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent"></div>
                  <div className="absolute bottom-4 left-4 right-4">
                    <h3 className="text-xl font-bold text-white">{category.name}</h3>
                    <p className="text-white/80 text-sm">{category.productCount} products</p>
                  </div>
                </div>
                <div className="p-6">
                  <p className="text-gray-600">{category.description}</p>
                  <div className="mt-4 flex items-center text-blue-600 font-medium group-hover:text-blue-700">
                    Shop now
                    <svg className="ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Featured Collections */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Featured Collections
            </h2>
            <p className="text-lg text-gray-600 max-w-3xl mx-auto">
              Explore our curated collections and find exactly what you're looking for.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="relative rounded-lg overflow-hidden h-80 group">
              <img 
                src="https://via.placeholder.com/800x400/3B82F6/FFFFFF?text=Summer+Collection" 
                alt="Summer Collection" 
                className="w-full h-full object-cover"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent"></div>
              <div className="absolute inset-0 flex flex-col justify-end p-8">
                <h3 className="text-2xl font-bold text-white mb-2">Summer Collection</h3>
                <p className="text-white/80 mb-4">Lightweight and breathable essentials for the warmer months</p>
                <Link href="/collections/summer" className="inline-flex items-center text-white font-medium hover:text-blue-200">
                  Shop Collection
                  <svg className="ml-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
              </div>
            </div>
            
            <div className="relative rounded-lg overflow-hidden h-80 group">
              <img 
                src="https://via.placeholder.com/800x400/10B981/FFFFFF?text=Tech+Deals" 
                alt="Tech Deals" 
                className="w-full h-full object-cover"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent"></div>
              <div className="absolute inset-0 flex flex-col justify-end p-8">
                <h3 className="text-2xl font-bold text-white mb-2">Tech Deals</h3>
                <p className="text-white/80 mb-4">Save on the latest gadgets and electronics</p>
                <Link href="/collections/tech-deals" className="inline-flex items-center text-white font-medium hover:text-blue-200">
                  Shop Collection
                  <svg className="ml-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}