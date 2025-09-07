'use client';

import { useQuery } from '@apollo/client';
import { useState } from 'react';
import { useParams } from 'next/navigation';
import { ProductCard } from '@/components/ProductCard';
import { GET_PRODUCTS, GET_CATEGORY } from '@/graphql/queries';
import type { Product, Category } from '@/graphql/queries';

export default function CategoryPage() {
  const { id } = useParams();
  const [selectedCategory, setSelectedCategory] = useState(id as string);
  
  const { loading: productsLoading, error: productsError, data: productsData } = useQuery(GET_PRODUCTS, {
    variables: {
      filter: { category: selectedCategory }
    }
  });

  const { loading: categoryLoading, error: categoryError, data: categoryData } = useQuery(GET_CATEGORY, {
    variables: {
      id: id
    }
  });

  // Map category ID to display name
  const getCategoryName = (categoryId: string) => {
    const categoryMap: Record<string, string> = {
      'electronics': 'Electronics',
      'fashion': 'Fashion',
      'home': 'Home & Kitchen',
      'beauty': 'Beauty & Personal Care',
      'sports': 'Sports & Outdoors',
      'books': 'Books & Media',
    };
    
    return categoryMap[categoryId] || categoryId;
  };

  if (productsLoading || categoryLoading) return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-8">Loading {getCategoryName(id as string)}...</h1>
        </div>
      </div>
    </div>
  );

  if (productsError || categoryError) return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-red-600 mb-8">Error Loading Data</h1>
          <p className="text-gray-600">{(productsError || categoryError)?.message}</p>
        </div>
      </div>
    </div>
  );

  const products: Product[] = productsData?.products || [];
  const category: Category | null = categoryData?.category || null;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Category Header */}
      <section className="bg-gradient-to-r from-blue-600 to-purple-600 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
            {category?.name || getCategoryName(id as string)}
          </h1>
          {category?.description && (
            <p className="text-xl text-white/90 max-w-3xl mx-auto">
              {category.description}
            </p>
          )}
        </div>
      </section>

      {/* Breadcrumb Navigation */}
      {category?.parent && (
        <section className="bg-white py-4 shadow-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <nav className="text-sm">
              <ol className="flex items-center space-x-2">
                <li>
                  <a href="/" className="text-blue-600 hover:text-blue-800">Home</a>
                </li>
                <li className="text-gray-400">/</li>
                <li>
                  <a href="/categories" className="text-blue-600 hover:text-blue-800">Categories</a>
                </li>
                {category.parent && (
                  <>
                    <li className="text-gray-400">/</li>
                    <li>
                      <a href={`/categories/${category.parent.id}`} className="text-blue-600 hover:text-blue-800">
                        {category.parent.name}
                      </a>
                    </li>
                  </>
                )}
                <li className="text-gray-400">/</li>
                <li className="text-gray-600">{category.name}</li>
              </ol>
            </nav>
          </div>
        </section>
      )}

      {/* Subcategories */}
      {category?.children && category.children.length > 0 && (
        <section className="bg-white py-8 shadow-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Subcategories</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {category.children.map((child) => (
                <a 
                  key={child.id}
                  href={`/categories/${child.id}`}
                  className="block p-4 border border-gray-200 rounded-lg hover:border-blue-500 hover:shadow-md transition-all"
                >
                  <h3 className="font-medium text-gray-900">{child.name}</h3>
                  {child.description && (
                    <p className="text-sm text-gray-500 mt-1">{child.description}</p>
                  )}
                </a>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* Category Filter */}
      <section className="bg-white py-6 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex overflow-x-auto gap-4 pb-2 scrollbar-hide">
            <button
              className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                selectedCategory === 'all'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
              onClick={() => setSelectedCategory('all')}
            >
              All Categories
            </button>
            <button
              className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                selectedCategory === 'electronics'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
              onClick={() => setSelectedCategory('electronics')}
            >
              Electronics
            </button>
            <button
              className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                selectedCategory === 'fashion'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
              onClick={() => setSelectedCategory('fashion')}
            >
              Fashion
            </button>
            <button
              className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                selectedCategory === 'home'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
              onClick={() => setSelectedCategory('home')}
            >
              Home & Kitchen
            </button>
            <button
              className={`whitespace-nowrap px-5 py-2 rounded-full font-medium text-sm transition-colors ${
                selectedCategory === 'beauty'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
              onClick={() => setSelectedCategory('beauty')}
            >
              Beauty & Personal Care
            </button>
          </div>
        </div>
      </section>

      {/* Products Grid */}
      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {products.length === 0 ? (
            <div className="text-center py-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">No products found</h2>
              <p className="text-gray-600">Try selecting a different category or check back later.</p>
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
                {products.map((product) => (
                  <ProductCard key={product.id} product={product} />
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
            </>
          )}
        </div>
      </section>
    </div>
  );
}