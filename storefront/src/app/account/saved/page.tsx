'use client';

import { useQuery, useMutation } from '@apollo/client';
import { GET_WISHLIST, REMOVE_FROM_WISHLIST } from '@/graphql/queries';
import type { WishlistItem } from '@/graphql/queries';

export default function SavedItems() {
  const { data, loading, error, refetch } = useQuery(GET_WISHLIST);
  const [removeFromWishlist] = useMutation(REMOVE_FROM_WISHLIST, {
    onCompleted: () => {
      refetch();
    }
  });
  
  const wishlistItems = data?.wishlist || [];

  const handleRemoveItem = async (productId: string) => {
    try {
      await removeFromWishlist({ variables: { productId } });
    } catch (err) {
      console.error('Error removing item from wishlist:', err);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-700">Error loading saved items: {error.message}</p>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Saved Items</h1>
      
      {wishlistItems.length === 0 ? (
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
          {wishlistItems.map((item: WishlistItem) => (
            <div key={item.id} className="bg-white rounded-lg shadow-sm hover:shadow-lg transition-shadow group relative">
              {item.product && (
                <>
                  <div className="relative aspect-square overflow-hidden rounded-t-lg">
                    <img
                      src={item.product.imageUrl || 'https://via.placeholder.com/600x600/3B82F6/FFFFFF?text=Product'}
                      alt={item.product.title}
                      className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                  
                  <div className="p-4">
                    <h3 className="font-medium text-gray-900 line-clamp-1">{item.product.title}</h3>
                    <p className="text-sm text-gray-500 mt-1 line-clamp-2">{item.product.description}</p>
                    <div className="mt-2 flex items-center justify-between">
                      <span className="text-lg font-medium text-gray-900">
                        ${item.product.priceRange?.minVariantPrice?.toFixed(2) || '0.00'}
                      </span>
                      <button 
                        onClick={() => handleRemoveItem(item.productId)}
                        className="text-sm font-medium text-red-600 hover:text-red-800"
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}