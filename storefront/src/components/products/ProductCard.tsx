'use client';

import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { ShoppingCartIcon, HeartIcon } from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import { useCartStore } from '@/stores/cart';
import { useMutation, useQuery } from '@apollo/client';
import { ADD_TO_WISHLIST, REMOVE_FROM_WISHLIST, GET_WISHLIST } from '@/graphql/queries';

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

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { addItem, openCart } = useCartStore();
  
  // Check if product is in wishlist
  const { data: wishlistData } = useQuery(GET_WISHLIST);
  const isInWishlist = wishlistData?.wishlist?.some((item: any) => item.productId === product.id);
  
  const [addToWishlist] = useMutation(ADD_TO_WISHLIST, {
    refetchQueries: [{ query: GET_WISHLIST }]
  });
  
  const [removeFromWishlist] = useMutation(REMOVE_FROM_WISHLIST, {
    refetchQueries: [{ query: GET_WISHLIST }]
  });

  const handleAddToCart = async () => {
    setIsLoading(true);
    
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 300));
    
    addItem({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.image,
    });
    
    setIsLoading(false);
    openCart();
  };

  const toggleWishlist = async () => {
    try {
      if (isInWishlist) {
        await removeFromWishlist({ variables: { productId: product.id } });
      } else {
        await addToWishlist({ variables: { productId: product.id } });
      }
    } catch (err) {
      console.error('Error updating wishlist:', err);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm hover:shadow-lg transition-shadow group">
      {/* Product Image */}
      <div className="relative aspect-square overflow-hidden rounded-t-lg">
        <Image
          src={product.image}
          alt={product.name}
          fill
          className="object-cover group-hover:scale-105 transition-transform duration-300"
        />
        
        {/* Favorite Button */}
        <button
          onClick={toggleWishlist}
          className="absolute top-3 right-3 p-2 bg-white rounded-full shadow-md hover:shadow-lg transition-shadow"
        >
          {isInWishlist ? (
            <HeartSolidIcon className="h-5 w-5 text-red-500" />
          ) : (
            <HeartIcon className="h-5 w-5 text-gray-400" />
          )}
        </button>

        {/* Stock Badge */}
        {!product.inventory.inStock && (
          <div className="absolute top-3 left-3 bg-red-500 text-white px-2 py-1 rounded text-xs font-medium">
            Out of Stock
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="p-4">
        <div className="mb-2">
          <span className="text-xs text-gray-500 uppercase tracking-wide">
            {product.category}
          </span>
        </div>
        
        <Link href={`/products/${product.id}`} className="block">
          <h3 className="font-medium text-gray-900 group-hover:text-blue-600 line-clamp-1">
            {product.name}
          </h3>
        </Link>
        
        <p className="text-sm text-gray-500 mt-1 line-clamp-2">
          {product.description}
        </p>
        
        <div className="mt-3 flex items-center justify-between">
          <span className="text-lg font-medium text-gray-900">
            ${product.price.toFixed(2)}
          </span>
          
          <button
            onClick={handleAddToCart}
            disabled={isLoading || !product.inventory.inStock}
            className={`flex items-center gap-1 px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
              isLoading || !product.inventory.inStock
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                : 'bg-blue-600 text-white hover:bg-blue-700'
            }`}
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-1 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Adding...
              </>
            ) : (
              <>
                <ShoppingCartIcon className="h-4 w-4" />
                Add to Cart
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}