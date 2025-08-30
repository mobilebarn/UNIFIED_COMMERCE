'use client';

import { useCart } from './CartProvider';

interface ProductCardProps {
  product: {
    id: string;
    name: string;
    price: number;
    rating: number;
    reviews: number;
    category: string;
    description: string;
    image: string;
    badges?: string[];
  };
}

export function ProductCard({ product }: ProductCardProps) {
  const { addToCart } = useCart();

  // Render star rating
  const renderRating = (rating: number) => {
    return (
      <div className="flex items-center">
        {[...Array(5)].map((_, i) => (
          <svg 
            key={i}
            className={`w-4 h-4 ${i < Math.floor(rating) ? 'text-yellow-400' : 'text-gray-300'}`}
            fill="currentColor" 
            viewBox="0 0 20 20" 
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
        ))}
        <span className="ml-1 text-sm text-gray-500">({rating})</span>
      </div>
    );
  };

  const handleAddToCart = () => {
    addToCart({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.image
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow group relative">
      {/* Badge */}
      {product.badges && product.badges.length > 0 && (
        <div className="absolute top-4 left-4 z-10 flex flex-col gap-2">
          {product.badges.map((badge, idx) => (
            <span key={idx} className={`px-2 py-1 text-xs font-semibold rounded-md ${
              badge === 'New' ? 'bg-blue-500 text-white' :
              badge === 'Sale' ? 'bg-red-500 text-white' :
              badge === 'Bestseller' ? 'bg-amber-500 text-white' :
              'bg-green-500 text-white'
            }`}>
              {badge}
            </span>
          ))}
        </div>
      )}
      
      {/* Product Image with hover effect */}
      <div className="relative overflow-hidden h-64">
        <img 
          src={product.image} 
          alt={product.name}
          className="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-300"
        />
        <div className="absolute inset-0 bg-black bg-opacity-20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
          <button className="bg-white text-gray-900 px-4 py-2 rounded-md font-medium text-sm transform translate-y-4 group-hover:translate-y-0 transition-transform">
            Quick View
          </button>
        </div>
      </div>
      
      {/* Product details */}
      <div className="p-4 border-t border-gray-100">
        <div className="mb-1">
          {renderRating(product.rating)}
          <span className="text-xs text-gray-500 ml-1">{product.reviews} reviews</span>
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-1 group-hover:text-blue-600 transition-colors">
          {product.name}
        </h3>
        <p className="text-sm text-gray-500 mb-4 line-clamp-2">
          {product.description}
        </p>
        <div className="flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900">
            ${product.price}
          </span>
          <button 
            onClick={handleAddToCart}
            className="bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 transition-colors flex items-center"
          >
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
