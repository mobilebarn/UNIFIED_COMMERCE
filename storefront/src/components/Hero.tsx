import Link from 'next/link';

export function Hero() {
  return (
    <div className="relative bg-gradient-to-r from-blue-600 to-blue-800 text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div className="text-center">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Welcome to Retail OS
          </h1>
          <p className="text-xl md:text-2xl mb-8 text-blue-100">
            Discover amazing products with our modern e-commerce platform
          </p>
          <div className="space-x-4">
            <Link
              href="/products"
              className="bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors inline-block"
            >
              Shop Now
            </Link>
            <Link
              href="/categories"
              className="border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white hover:text-blue-600 transition-colors inline-block"
            >
              Browse Categories
            </Link>
          </div>
        </div>
      </div>
      
      {/* Decorative background elements */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden">
        <div className="absolute top-10 left-10 w-20 h-20 bg-white opacity-10 rounded-full"></div>
        <div className="absolute top-20 right-20 w-16 h-16 bg-white opacity-10 rounded-full"></div>
        <div className="absolute bottom-10 left-1/4 w-12 h-12 bg-white opacity-10 rounded-full"></div>
        <div className="absolute bottom-20 right-1/3 w-24 h-24 bg-white opacity-10 rounded-full"></div>
      </div>
    </div>
  );
}
