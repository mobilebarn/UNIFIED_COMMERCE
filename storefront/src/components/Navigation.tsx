'use client';

import Link from 'next/link';
import { useState, useRef, useEffect } from 'react';
import { useCart } from './CartProvider';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'next/navigation';
import { useMutation } from '@apollo/client';
import { LOGOUT } from '@/graphql/queries';
import { useSearchSuggestions } from '@/hooks/useSearchSuggestions';

export function Navigation() {
  const { getCartCount } = useCart();
  const { user, isAuthenticated, logout } = useAuthStore();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [cartOpen, setCartOpen] = useState(false);
  const [accountMenuOpen, setAccountMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [showSuggestions, setShowSuggestions] = useState(false);
  
  const router = useRouter();
  const searchInputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);
  
  const { suggestions, loading: suggestionsLoading } = useSearchSuggestions(searchQuery);
  
  const [logoutMutation] = useMutation(LOGOUT);

  const toggleMobileMenu = () => setMobileMenuOpen(!mobileMenuOpen);
  const toggleSearch = () => setSearchOpen(!searchOpen);
  const toggleCart = () => setCartOpen(!cartOpen);
  const toggleAccountMenu = () => setAccountMenuOpen(!accountMenuOpen);
  
  const handleLogout = async () => {
    try {
      await logoutMutation();
      logout();
      router.push('/');
      router.refresh();
    } catch (error) {
      console.error('Logout error:', error);
      // Still logout locally even if server call fails
      logout();
      router.push('/');
      router.refresh();
    }
  };
  
  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery('');
      setShowSuggestions(false);
    }
  };
  
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
    setShowSuggestions(e.target.value.length > 1);
  };
  
  const handleSuggestionClick = (suggestion: any) => {
    setSearchQuery(suggestion.title);
    setShowSuggestions(false);
    router.push(`/search?q=${encodeURIComponent(suggestion.title)}`);
  };
  
  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionsRef.current && 
        !suggestionsRef.current.contains(event.target as Node) &&
        searchInputRef.current && 
        !searchInputRef.current.contains(event.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };
    
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);
  
  const cartCount = getCartCount();
  
  return (
    <nav className="bg-white shadow-sm sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo and navigation */}
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Link href="/" className="flex items-center">
                <svg className="h-8 w-8 text-blue-600" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 4L4 8L12 12L20 8L12 4Z" />
                  <path d="M4 12L12 16L20 12" />
                  <path d="M4 16L12 20L20 16" />
                </svg>
                <span className="ml-2 text-2xl font-bold text-blue-600">
                  Unified Commerce
                </span>
              </Link>
            </div>
            <div className="hidden md:block ml-10">
              <div className="flex items-baseline space-x-4">
                <Link href="/products" className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
                  All Products
                </Link>
                <div className="relative group">
                  <button className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center">
                    Categories
                    <svg className="ml-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path>
                    </svg>
                  </button>
                  <div className="absolute left-0 mt-2 w-48 bg-white rounded-md shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                    <div className="py-1">
                      <Link href="/categories/electronics" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Electronics</Link>
                      <Link href="/categories/fashion" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Fashion</Link>
                      <Link href="/categories/home" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Home & Kitchen</Link>
                      <Link href="/categories/beauty" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Beauty & Personal Care</Link>
                      <Link href="/categories" className="block px-4 py-2 text-sm text-blue-600 border-t border-gray-100">View All Categories</Link>
                    </div>
                  </div>
                </div>
                <Link href="/deals" className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
                  Deals
                </Link>
                <Link href="/new-arrivals" className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
                  New Arrivals
                </Link>
              </div>
            </div>
          </div>
          
          {/* Search, account and cart */}
          <div className="flex items-center space-x-4">
            <div className="hidden md:block relative" ref={suggestionsRef}>
              <form onSubmit={handleSearchSubmit} className="relative">
                <input 
                  ref={searchInputRef}
                  type="text" 
                  placeholder="Search products..." 
                  value={searchQuery}
                  onChange={handleSearchChange}
                  className="w-64 pl-10 pr-10 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <svg className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <button 
                  type="submit"
                  className="absolute right-3 top-2.5 h-5 w-5 text-gray-400 hover:text-blue-600"
                >
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                  </svg>
                </button>
              </form>
              
              {/* Search suggestions dropdown */}
              {showSuggestions && (
                <div className="absolute left-0 right-0 mt-1 bg-white rounded-md shadow-lg z-50 border border-gray-200 max-h-80 overflow-y-auto">
                  {suggestionsLoading ? (
                    <div className="px-4 py-3 text-sm text-gray-500">Loading suggestions...</div>
                  ) : suggestions.length > 0 ? (
                    suggestions.map((suggestion: any) => (
                      <button
                        key={suggestion.id}
                        type="button"
                        onClick={() => handleSuggestionClick(suggestion)}
                        className="w-full text-left px-4 py-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 flex items-center"
                      >
                        {suggestion.imageUrl && (
                          <img 
                            src={suggestion.imageUrl} 
                            alt={suggestion.title} 
                            className="w-8 h-8 object-cover rounded mr-3"
                          />
                        )}
                        <div>
                          <div className="font-medium text-gray-900">{suggestion.title}</div>
                          <div className="text-xs text-gray-500 capitalize">{suggestion.type.toLowerCase()}</div>
                        </div>
                      </button>
                    ))
                  ) : (
                    <div className="px-4 py-3 text-sm text-gray-500">No suggestions found</div>
                  )}
                </div>
              )}
            </div>
            
            {isAuthenticated ? (
              <div className="relative group">
                <button 
                  onClick={toggleAccountMenu}
                  className="text-gray-600 hover:text-blue-600 p-2 flex items-center"
                >
                  <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span className="ml-1 text-sm hidden md:inline">{user?.firstName}</span>
                </button>
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                  <div className="py-1">
                    <Link href="/account/profile" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">My Profile</Link>
                    <Link href="/account/orders" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">My Orders</Link>
                    <Link href="/account/wishlist" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Wishlist</Link>
                    <button 
                      onClick={handleLogout}
                      className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 border-t border-gray-100"
                    >
                      Sign Out
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <Link href="/login" className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium">
                  Sign In
                </Link>
                <Link href="/register" className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 transition-colors">
                  Register
                </Link>
              </div>
            )}
            
            <div className="relative">
              <button 
                onClick={toggleCart}
                className="relative bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors flex items-center"
              >
                <svg className="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                Cart ({cartCount})
                {cartCount > 0 && (
                  <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                    {cartCount > 9 ? '9+' : cartCount}
                  </span>
                )}
              </button>
              
              {/* Mini Cart Dropdown */}
              {cartOpen && (
                <div className="absolute right-0 mt-2 w-80 bg-white rounded-md shadow-lg z-50">
                  <div className="p-4">
                    <div className="flex items-center justify-between border-b border-gray-200 pb-3">
                      <h3 className="text-lg font-medium text-gray-900">Shopping Cart</h3>
                      <span className="text-sm text-gray-500">{cartCount} items</span>
                    </div>
                    
                    {cartCount === 0 ? (
                      <div className="py-6 text-center">
                        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"></path>
                        </svg>
                        <p className="mt-4 text-sm text-gray-500">Your cart is empty</p>
                        <button 
                          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md text-sm hover:bg-blue-700 transition"
                          onClick={toggleCart}
                        >
                          Continue Shopping
                        </button>
                      </div>
                    ) : (
                      <>
                        <div className="max-h-60 overflow-y-auto py-3">
                          {/* This would be populated with actual cart items */}
                          <div className="flex py-3 border-b border-gray-100">
                            <div className="h-16 w-16 flex-shrink-0 overflow-hidden rounded-md border border-gray-200">
                              <img src="https://via.placeholder.com/80x80/F59E0B/FFFFFF?text=Demo" alt="Product" className="h-full w-full object-cover object-center" />
                            </div>
                            <div className="ml-4 flex flex-1 flex-col">
                              <div className="flex justify-between text-base font-medium text-gray-900">
                                <h3 className="text-sm">Demo Product</h3>
                                <p className="ml-4 text-sm">$59.99</p>
                              </div>
                              <div className="flex flex-1 items-end justify-between text-sm">
                                <p className="text-gray-500">Qty 1</p>
                                <button className="text-red-500 hover:text-red-700">Remove</button>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div className="border-t border-gray-100 pt-4">
                          <div className="flex justify-between text-base font-medium text-gray-900 mb-4">
                            <p>Subtotal</p>
                            <p>$59.99</p>
                          </div>
                          <div className="flex flex-col space-y-2">
                            <Link 
                              href="/cart" 
                              onClick={toggleCart}
                              className="w-full bg-blue-600 text-white px-6 py-2 rounded-md text-center hover:bg-blue-700 transition"
                            >
                              View Cart
                            </Link>
                            <Link 
                              href="/checkout" 
                              onClick={toggleCart}
                              className="w-full bg-gray-100 text-gray-900 px-6 py-2 rounded-md text-center hover:bg-gray-200 transition"
                            >
                              Checkout
                            </Link>
                          </div>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
      
      {/* Mobile menu button */}
      <div className="md:hidden border-t border-gray-200 py-2 px-4">
        <div className="flex justify-between">
          <button 
            onClick={toggleMobileMenu}
            className="text-gray-600 hover:text-blue-600"
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <button 
            onClick={toggleSearch}
            className="text-gray-600 hover:text-blue-600"
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
        </div>
      </div>
      
      {/* Mobile search */}
      {searchOpen && (
        <div className="md:hidden p-4 border-t border-gray-200 bg-white">
          <div className="relative" ref={suggestionsRef}>
            <form onSubmit={handleSearchSubmit}>
              <input 
                ref={searchInputRef}
                type="text" 
                placeholder="Search products..." 
                value={searchQuery}
                onChange={handleSearchChange}
                className="w-full pl-10 pr-10 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <svg className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <button 
                type="submit"
                className="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600"
              >
                <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </button>
            </form>
            
            {/* Mobile search suggestions dropdown */}
            {showSuggestions && (
              <div className="absolute left-0 right-0 mt-1 bg-white rounded-md shadow-lg z-50 border border-gray-200 max-h-80 overflow-y-auto">
                {suggestionsLoading ? (
                  <div className="px-4 py-3 text-sm text-gray-500">Loading suggestions...</div>
                ) : suggestions.length > 0 ? (
                  suggestions.map((suggestion: any) => (
                    <button
                      key={suggestion.id}
                      type="button"
                      onClick={() => handleSuggestionClick(suggestion)}
                      className="w-full text-left px-4 py-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 flex items-center"
                    >
                      {suggestion.imageUrl && (
                        <img 
                          src={suggestion.imageUrl} 
                          alt={suggestion.title} 
                          className="w-8 h-8 object-cover rounded mr-3"
                        />
                      )}
                      <div>
                        <div className="font-medium text-gray-900">{suggestion.title}</div>
                        <div className="text-xs text-gray-500 capitalize">{suggestion.type.toLowerCase()}</div>
                      </div>
                    </button>
                  ))
                ) : (
                  <div className="px-4 py-3 text-sm text-gray-500">No suggestions found</div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
      
      {/* Mobile menu */}
      {mobileMenuOpen && (
        <div className="md:hidden border-t border-gray-200 bg-white">
          <div className="px-2 pt-2 pb-3 space-y-1">
            <Link href="/products" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
              All Products
            </Link>
            <button className="flex items-center justify-between w-full px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
              <span>Categories</span>
              <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </button>
            <Link href="/deals" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
              Deals
            </Link>
            <Link href="/new-arrivals" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
              New Arrivals
            </Link>
            <div className="border-t border-gray-200 pt-4 pb-2">
              {isAuthenticated ? (
                <>
                  <Link href="/account" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
                    My Account
                  </Link>
                  <Link href="/account/orders" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
                    My Orders
                  </Link>
                  <button 
                    onClick={handleLogout}
                    className="block w-full text-left px-3 py-2 rounded-md text-base font-medium text-red-600 hover:bg-gray-50"
                  >
                    Sign Out
                  </button>
                </>
              ) : (
                <>
                  <Link href="/login" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
                    Sign In
                  </Link>
                  <Link href="/register" className="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:bg-gray-50">
                    Register
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </nav>
  );
}