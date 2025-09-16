'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function NotFound() {
  const pathname = usePathname();
  
  // Suggested navigation links based on common paths
  const suggestedLinks = [
    { href: '/', label: 'Home' },
    { href: '/products', label: 'All Products' },
    { href: '/categories', label: 'Categories' },
    { href: '/deals', label: 'Deals & Promotions' },
    { href: '/account', label: 'My Account' },
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto text-center">
          {/* Error Icon */}
          <div className="flex justify-center mb-8">
            <div className="flex items-center justify-center h-20 w-20 rounded-full bg-red-100">
              <svg className="h-12 w-12 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
              </svg>
            </div>
          </div>
          
          {/* Error Title and Message */}
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-6">Page Not Found</h1>
          <p className="text-lg md:text-xl text-gray-600 mb-2">
            Sorry, we couldn't find the page you're looking for.
          </p>
          <p className="text-md text-gray-500 mb-8">
            The requested URL <span className="font-mono bg-gray-100 px-2 py-1 rounded">{pathname}</span> was not found on this server.
          </p>
          
          {/* Search Form */}
          <div className="mb-12 max-w-2xl mx-auto">
            <form action="/search" method="GET" className="flex gap-2">
              <input
                type="text"
                name="q"
                placeholder="Search products..."
                className="flex-grow px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
              <button
                type="submit"
                className="px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Search
              </button>
            </form>
          </div>
          
          {/* Quick Navigation */}
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Quick Navigation</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
              {suggestedLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  className="flex flex-col items-center justify-center p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
                >
                  <span className="text-blue-600 font-medium">{link.label}</span>
                </Link>
              ))}
            </div>
          </div>
          
          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row justify-center gap-4 mb-12">
            <Link 
              href="/" 
              className="inline-flex items-center justify-center px-8 py-4 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Go back home
            </Link>
            <Link 
              href="/products" 
              className="inline-flex items-center justify-center px-8 py-4 border border-gray-300 shadow-sm text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Continue Shopping
            </Link>
            <Link 
              href="/contact" 
              className="inline-flex items-center justify-center px-8 py-4 border border-gray-300 shadow-sm text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Contact Support
            </Link>
          </div>
          
          {/* Help Section */}
          <div className="bg-blue-50 rounded-lg p-6 max-w-2xl mx-auto">
            <h3 className="text-lg font-bold text-gray-900 mb-2">Need Help?</h3>
            <p className="text-gray-600 mb-4">
              If you're sure this page should exist, please contact our support team or try searching for what you need.
            </p>
            <div className="flex flex-col sm:flex-row justify-center gap-3">
              <a 
                href="mailto:support@unifiedcommerce.com" 
                className="text-blue-600 hover:text-blue-800 font-medium"
              >
                support@unifiedcommerce.com
              </a>
              <span className="hidden sm:inline text-gray-400">|</span>
              <a 
                href="tel:+18001234567" 
                className="text-blue-600 hover:text-blue-800 font-medium"
              >
                1-800-123-4567
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}