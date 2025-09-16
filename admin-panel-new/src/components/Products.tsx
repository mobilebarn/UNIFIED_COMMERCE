import React, { useState, useMemo } from 'react';
import { useProducts, useCreateProduct, useUpdateProduct, useDeleteProduct } from '../hooks/useGraphQL';
import { Product } from '../lib/graphql';

export default function Products() {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [sortBy, setSortBy] = useState('newest');
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);
  
  // Note: The useProducts hook fetches all products. For large catalogs, 
  // it would be more performant to implement server-side filtering, sorting, and pagination
  // by passing variables to the useProducts hook.
  const { data, loading, error, refetch } = useProducts();
  const products = data?.products;

  const [createProduct] = useCreateProduct();
  const [updateProduct] = useUpdateProduct();
  const [deleteProduct] = useDeleteProduct();

  // Calculate category counts from actual data
  const categoryStats = useMemo(() => {
    return products?.reduce((acc, product) => {
      const categoryName = product.categories?.[0]?.name || 'Uncategorized';
      acc[categoryName] = (acc[categoryName] || 0) + 1;
      return acc;
    }, {} as Record<string, number>) || {};
  }, [products]);

  const uniqueCategories = Object.keys(categoryStats);
  const totalProducts = products?.length || 0;

  // Filter and sort products
  const filteredProducts = useMemo(() => {
    return products?.filter(product => {
      const matchesSearch = product.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         product.description?.toLowerCase().includes(searchTerm.toLowerCase());
      const categoryName = product.categories?.[0]?.name || 'Uncategorized';
      const matchesCategory = !selectedCategory || categoryName === selectedCategory;
      return matchesSearch && matchesCategory;
    }) || [];
  }, [products, searchTerm, selectedCategory]);

  const sortedProducts = useMemo(() => {
    return [...filteredProducts].sort((a, b) => {
      switch (sortBy) {
        case 'oldest':
          return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
        case 'price-low':
          return a.price - b.price;
        case 'price-high':
          return b.price - a.price;
        case 'newest':
        default:
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
      }
    });
  }, [filteredProducts, sortBy]);

  const handleCreateProduct = async () => {
    try {
      await createProduct({
        variables: {
          input: {
            title: 'New Draft Product',
            description: 'A description for the new product.',
            price: 99.99,
            status: 'DRAFT',
          }
        }
      });
      // The useCreateProduct hook should automatically refetch queries, 
      // but an explicit refetch can be used if needed.
    } catch (error) {
      console.error('Error creating product:', error);
    }
  };

  const handleEditProduct = (product: Product) => {
    setEditingProduct(product);
    setShowEditModal(true);
  };

  const handleDeleteProduct = async (productId: string) => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      try {
        await deleteProduct({
          variables: {
            id: productId
          }
        });
      } catch (error) {
        console.error('Error deleting product:', error);
        alert('Failed to delete product: ' + (error as Error).message);
      }
    }
  };

  const getProductStock = (product: Product) => {
    return product.variants?.reduce((sum, variant) => {
      if (!variant.inventory) return sum;
      
      // Handle both array and object types for inventory
      if (Array.isArray(variant.inventory)) {
        return sum + variant.inventory.reduce((total, item) => total + ((item as any).quantity || 0), 0);
      } else {
        return sum + ((variant.inventory as any).quantity || 0);
      }
    }, 0) || 0;
  }

  const getStockStatus = (stockLevel: number) => {
    if (stockLevel === 0) return { status: 'Out of Stock', color: 'bg-red-100 text-red-800', dot: 'bg-red-500' };
    if (stockLevel < 20) return { status: 'Low Stock', color: 'bg-yellow-100 text-yellow-800', dot: 'bg-yellow-500' };
    return { status: 'Active', color: 'bg-green-100 text-green-800', dot: 'bg-green-500' };
  };

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading products...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">Error loading products: {error.message}</p>
          <button 
            onClick={() => refetch()}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full">
      <div className="p-6">
        {/* Header */}
        <div className="mb-8 flex flex-col lg:flex-row lg:justify-between lg:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Products</h1>
            <p className="text-gray-600 mt-1">Manage your product catalog and inventory</p>
          </div>
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="relative">
              <input
                type="text"
                placeholder="Search products..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 pr-4 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 w-full sm:w-64"
              />
              <svg className="absolute left-3 top-3 w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
            </div>
            <button 
              onClick={handleCreateProduct}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm flex items-center gap-2 transition-colors whitespace-nowrap"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4"></path>
              </svg>
              Add Product
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 lg:gap-8">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-4 lg:p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Categories</h3>
              <ul className="space-y-2">
                <li>
                  <button 
                    onClick={() => setSelectedCategory('')}
                    className={`flex items-center justify-between py-2 px-3 rounded-md w-full text-left ${
                      !selectedCategory ? 'bg-blue-50 text-blue-700' : 'hover:bg-gray-50 text-gray-700'
                    }`}
                  >
                    <span className="text-sm font-medium">All Products</span>
                    <span className={`text-xs py-1 px-2 rounded-full ${
                      !selectedCategory ? 'bg-blue-100 text-blue-600' : 'bg-gray-100 text-gray-600'
                    }`}>
                      {totalProducts}
                    </span>
                  </button>
                </li>
                {uniqueCategories.map(category => (
                  <li key={category}>
                    <button 
                      onClick={() => setSelectedCategory(category)}
                      className={`flex items-center justify-between py-2 px-3 rounded-md w-full text-left ${
                        selectedCategory === category ? 'bg-blue-50 text-blue-700' : 'hover:bg-gray-50 text-gray-700'
                      }`}
                    >
                      <span className="text-sm">{category}</span>
                      <span className={`text-xs py-1 px-2 rounded-full ${
                        selectedCategory === category ? 'bg-blue-100 text-blue-600' : 'bg-gray-100 text-gray-600'
                      }`}>
                        {categoryStats[category]}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>
          
          {/* Main Content */}
          <div className="lg:col-span-3">
            <div className="bg-white rounded-xl shadow-sm border border-gray-100">
              <div className="p-4 lg:p-6 border-b border-gray-100 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                <div className="flex items-center gap-4">
                  <h3 className="text-lg font-semibold text-gray-900">Product List</h3>
                  <span className="text-sm bg-gray-100 text-gray-600 py-1 px-3 rounded-full">
                    {sortedProducts.length} products
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <select 
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value)}
                    className="text-sm border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="newest">Sort by: Newest</option>
                    <option value="oldest">Sort by: Oldest</option>
                    <option value="price-low">Sort by: Price (Low to High)</option>
                    <option value="price-high">Sort by: Price (High to Low)</option>
                  </select>
                </div>
              </div>
              
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Inventory</th>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                      <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {sortedProducts.map((product) => {
                      const stockLevel = getProductStock(product);
                      const stockStatus = getStockStatus(stockLevel);
                      return (
                        <tr key={product.id} className="hover:bg-gray-50">
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <div className="flex-shrink-0 h-12 w-12">
                                <img 
                                  className="h-12 w-12 rounded-lg object-cover" 
                                  src={product.imageUrl || `https://via.placeholder.com/300x300/3B82F6/FFFFFF?text=${product.title.charAt(0)}`} 
                                  alt={product.title}
                                  onError={(e) => {
                                    const target = e.target as HTMLImageElement;
                                    target.src = `https://via.placeholder.com/300x300/3B82F6/FFFFFF?text=${product.title.charAt(0)}`;
                                  }}
                                />
                              </div>
                              <div className="ml-4">
                                <div className="text-sm font-medium text-gray-900">{product.title}</div>
                                <div className="text-sm text-gray-500">{product.categories?.[0]?.name || 'N/A'}</div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm font-medium text-gray-900">${product.price.toFixed(2)}</div>
                            <div className="text-sm text-gray-500">USD</div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <div className={`mr-2 h-2 w-2 rounded-full ${stockStatus.dot}`}></div>
                              <div className="text-sm text-gray-900">{stockLevel} in stock</div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <span className={`inline-flex px-2 py-1 text-xs font-semibold leading-5 rounded-full ${stockStatus.color}`}>
                              {product.status}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                            <button 
                              onClick={() => handleEditProduct(product)}
                              className="text-blue-600 hover:text-blue-900 mr-3"
                            >
                              Edit
                            </button>
                            <button 
                              onClick={() => handleDeleteProduct(product.id)}
                              className="text-red-600 hover:text-red-900"
                            >
                              Delete
                            </button>
                          </td>
                        </tr>
                      );
                    })}
                    {sortedProducts.length === 0 && (
                      <tr>
                        <td colSpan={5} className="px-6 py-12 text-center text-gray-500">
                          {searchTerm || selectedCategory ? 'No products match your filters.' : 'No products found.'}
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
              
              <div className="px-4 lg:px-6 py-4 flex flex-col sm:flex-row sm:items-center sm:justify-between border-t border-gray-200 gap-4">
                <div className="text-sm text-gray-500">
                  Showing {sortedProducts.length > 0 ? '1' : '0'} to {sortedProducts.length} of {totalProducts} results
                  {(searchTerm || selectedCategory) && ` (filtered)`}
                </div>
                <div className="flex items-center space-x-2">
                  <button 
                    disabled={true}
                    className="px-3 py-2 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm transition-colors disabled:opacity-50"
                  >
                    Previous
                  </button>
                  <button className="px-3 py-2 rounded-md bg-blue-600 border border-blue-600 text-white text-sm">1</button>
                  <button 
                    disabled={true}
                    className="px-3 py-2 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm transition-colors disabled:opacity-50"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      {/* Edit Product Modal */}
      {showEditModal && editingProduct && (
        <EditProductModal 
          product={editingProduct} 
          onClose={() => setShowEditModal(false)} 
          onUpdate={async (updatedProduct) => {
            try {
              await updateProduct({
                variables: {
                  id: updatedProduct.id,
                  input: {
                    title: updatedProduct.title,
                    description: updatedProduct.description,
                    price: updatedProduct.price,
                    status: updatedProduct.status,
                  }
                }
              });
              setShowEditModal(false);
            } catch (error) {
              console.error('Error updating product:', error);
              alert('Failed to update product: ' + (error as Error).message);
            }
          }} 
        />
      )}
    </div>
  );
}

// Edit Product Modal Component
const EditProductModal: React.FC<{
  product: Product;
  onClose: () => void;
  onUpdate: (product: Product) => void;
}> = ({ product, onClose, onUpdate }) => {
  const [title, setTitle] = useState(product.title);
  const [description, setDescription] = useState(product.description || '');
  const [price, setPrice] = useState(product.price.toString());
  const [status, setStatus] = useState(product.status);
  const [saving, setSaving] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    
    try {
      await onUpdate({
        ...product,
        title,
        description,
        price: parseFloat(price),
        status
      });
    } catch (error) {
      console.error('Error updating product:', error);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/30" onClick={onClose} />
      <div className="relative bg-white rounded-xl shadow-lg w-full max-w-lg border border-gray-100 animate-fadeIn">
        <div className="px-5 py-4 border-b border-gray-100 flex justify-between items-center">
          <h2 className="text-lg font-semibold text-gray-900">Edit Product</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600" aria-label="Close">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>
        <div className="px-5 py-4 space-y-4 max-h-[60vh] overflow-y-auto">
          <form id="edit-product-form" onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <input 
                value={title} 
                onChange={e => setTitle(e.target.value)} 
                required 
                className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm" 
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea 
                value={description} 
                onChange={e => setDescription(e.target.value)} 
                className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm" 
                rows={3}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Price (USD)</label>
              <input 
                type="number" 
                step="0.01" 
                value={price} 
                onChange={e => setPrice(e.target.value)} 
                required 
                className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm" 
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select 
                value={status} 
                onChange={e => setStatus(e.target.value as any)} 
                className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm"
              >
                <option value="active">Active</option>
                <option value="draft">Draft</option>
                <option value="archived">Archived</option>
              </select>
            </div>
          </form>
        </div>
        <div className="px-5 py-4 border-t border-gray-100 bg-gray-50 rounded-b-xl">
          <div className="flex justify-end gap-3">
            <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Cancel</button>
            <button 
              form="edit-product-form" 
              disabled={saving} 
              className="px-4 py-2 text-sm rounded-md bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50"
            >
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};