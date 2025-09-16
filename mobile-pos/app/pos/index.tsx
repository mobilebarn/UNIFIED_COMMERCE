import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, FlatList, TextInput, Alert } from 'react-native';
import { useQuery, useMutation, useApolloClient } from '@apollo/client';
import { GET_PRODUCTS, SEARCH_PRODUCTS, GET_CART_BY_SESSION, CREATE_CART, ADD_CART_LINE_ITEM, UPDATE_CART_LINE_ITEM, REMOVE_CART_LINE_ITEM } from '../../graphql';
import { localStorage } from '../../utils/storage';

// Product type based on GraphQL schema
type Product = {
  id: string;
  title: string;
  priceRange: {
    minVariantPrice: number;
  };
  variants: {
    id: string;
    title: string;
    price: number;
    sku: string;
    barcode: string;
    inventoryQuantity: number;
  }[];
};

// Cart item type
type CartItem = {
  id: string;
  productId: string;
  productVariantId: string;
  name: string;
  sku: string;
  price: number;
  quantity: number;
  productImage?: string;
};

export default function POSScreen() {
  const [cart, setCart] = useState<CartItem[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [sessionId, setSessionId] = useState('');
  const client = useApolloClient();

  // Generate a session ID if one doesn't exist
  useEffect(() => {
    const sessionId = localStorage.getItem('posSessionId') || 'session_' + Date.now();
    localStorage.setItem('posSessionId', sessionId);
    setSessionId(sessionId);
    
    // Load existing cart for this session
    loadCartForSession(sessionId);
  }, []);

  // Load cart for the current session
  const loadCartForSession = async (sessionId: string) => {
    try {
      const { data } = await client.query({
        query: GET_CART_BY_SESSION,
        variables: { sessionId },
        fetchPolicy: 'network-only'
      });
      
      if (data.cartBySession) {
        setCart(data.cartBySession.lineItems);
      }
    } catch (error) {
      console.log('No existing cart found for session');
    }
  };

  // Get products query
  const { data: productsData, loading: productsLoading, error: productsError, refetch } = useQuery(
    SEARCH_PRODUCTS,
    {
      variables: { query: searchQuery || '' },
      skip: searchQuery.length > 0 && searchQuery.length < 2,
    }
  );

  // Create cart mutation
  const [createCart] = useMutation(CREATE_CART);
  
  // Add item to cart mutation
  const [addCartLineItem] = useMutation(ADD_CART_LINE_ITEM);
  
  // Update item quantity mutation
  const [updateCartLineItem] = useMutation(UPDATE_CART_LINE_ITEM);
  
  // Remove item from cart mutation
  const [removeCartLineItem] = useMutation(REMOVE_CART_LINE_ITEM);

  // Add item to cart
  const addToCart = async (product: Product) => {
    try {
      // Use the first variant if available, otherwise use product info
      const variant = product.variants && product.variants.length > 0 
        ? product.variants[0] 
        : null;
      
      const itemToAdd = {
        productId: product.id,
        productVariantId: variant?.id || null,
        name: product.title,
        sku: variant?.sku || '',
        price: variant?.price || product.priceRange.minVariantPrice,
        quantity: 1,
      };

      // First, check if we have a cart for this session
      let cartId = localStorage.getItem(`cartId_${sessionId}`);
      
      if (!cartId) {
        // Create a new cart if one doesn't exist
        const { data: cartData } = await createCart({
          variables: {
            input: {
              sessionId,
              merchantId: 'merchant_1', // In a real app, this would come from auth context
              currency: 'USD',
            }
          }
        });
        
        cartId = cartData.createCart.id;
        localStorage.setItem(`cartId_${sessionId}`, cartId);
      }

      // Add item to cart
      const { data } = await addCartLineItem({
        variables: {
          input: {
            cartId,
            ...itemToAdd
          }
        }
      });

      // Update local cart state
      setCart(prevCart => {
        const existingItem = prevCart.find(item => item.productId === product.id);
        if (existingItem) {
          return prevCart.map(item =>
            item.productId === product.id
              ? { ...item, quantity: item.quantity + 1 }
              : item
          );
        } else {
          return [...prevCart, { ...itemToAdd, id: data.addCartLineItem.id }];
        }
      });
    } catch (error) {
      Alert.alert('Error', 'Failed to add item to cart');
      console.error(error);
    }
  };

  // Remove item from cart
  const removeFromCart = async (id: string) => {
    try {
      await removeCartLineItem({
        variables: { id }
      });
      
      setCart(prevCart => prevCart.filter(item => item.id !== id));
    } catch (error) {
      Alert.alert('Error', 'Failed to remove item from cart');
      console.error(error);
    }
  };

  // Update item quantity
  const updateQuantity = async (id: string, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(id);
      return;
    }
    
    try {
      await updateCartLineItem({
        variables: {
          id,
          input: { quantity }
        }
      });
      
      setCart(prevCart =>
        prevCart.map(item => (item.id === id ? { ...item, quantity } : item))
      );
    } catch (error) {
      Alert.alert('Error', 'Failed to update item quantity');
      console.error(error);
    }
  };

  // Calculate total
  const calculateTotal = () => {
    return cart.reduce((total, item) => total + item.price * item.quantity, 0);
  };

  // Handle search
  const handleSearch = (query: string) => {
    setSearchQuery(query);
    if (query.length === 0 || query.length > 1) {
      refetch({ query });
    }
  };

  if (productsLoading) {
    return (
      <View style={styles.container}>
        <Text>Loading products...</Text>
      </View>
    );
  }

  if (productsError) {
    return (
      <View style={styles.container}>
        <Text>Error loading products: {productsError.message}</Text>
      </View>
    );
  }

  const products = productsData?.products || [];

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Retail OS Point of Sale</Text>
      
      {/* Search Bar */}
      <View style={styles.searchContainer}>
        <TextInput
          style={styles.searchInput}
          placeholder="Search products by name or barcode"
          value={searchQuery}
          onChangeText={handleSearch}
        />
      </View>
      
      {/* Product Grid */}
      <FlatList
        data={products}
        keyExtractor={item => item.id}
        numColumns={2}
        renderItem={({ item }) => (
          <TouchableOpacity
            style={styles.productCard}
            onPress={() => addToCart(item)}
          >
            <Text style={styles.productName}>{item.title}</Text>
            <Text style={styles.productPrice}>${item.priceRange.minVariantPrice.toFixed(2)}</Text>
            {item.variants && item.variants.length > 0 && (
              <Text style={styles.productBarcode}>{item.variants[0].barcode}</Text>
            )}
          </TouchableOpacity>
        )}
        contentContainerStyle={styles.productGrid}
      />
      
      {/* Cart Summary */}
      <View style={styles.cartContainer}>
        <Text style={styles.cartHeader}>Cart ({cart.length} items)</Text>
        <FlatList
          data={cart}
          keyExtractor={item => item.id}
          renderItem={({ item }) => (
            <View style={styles.cartItem}>
              <View style={styles.cartItemInfo}>
                <Text style={styles.cartItemName}>{item.name}</Text>
                <Text style={styles.cartItemPrice}>${item.price.toFixed(2)}</Text>
              </View>
              <View style={styles.quantityControls}>
                <TouchableOpacity
                  style={styles.quantityButton}
                  onPress={() => updateQuantity(item.id, item.quantity - 1)}
                >
                  <Text style={styles.quantityButtonText}>-</Text>
                </TouchableOpacity>
                <Text style={styles.quantityText}>{item.quantity}</Text>
                <TouchableOpacity
                  style={styles.quantityButton}
                  onPress={() => updateQuantity(item.id, item.quantity + 1)}
                >
                  <Text style={styles.quantityButtonText}>+</Text>
                </TouchableOpacity>
              </View>
            </View>
          )}
        />
        
        <View style={styles.totalContainer}>
          <Text style={styles.totalText}>Total: ${calculateTotal().toFixed(2)}</Text>
          <TouchableOpacity style={styles.checkoutButton}>
            <Text style={styles.checkoutButtonText}>Checkout</Text>
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    padding: 20,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  searchContainer: {
    padding: 15,
    backgroundColor: '#fff',
  },
  searchInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
  },
  productGrid: {
    padding: 10,
  },
  productCard: {
    flex: 1,
    backgroundColor: '#fff',
    margin: 5,
    padding: 15,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
    alignItems: 'center',
  },
  productName: {
    fontSize: 16,
    fontWeight: 'bold',
    textAlign: 'center',
  },
  productPrice: {
    fontSize: 18,
    color: '#2e7d32',
    marginVertical: 5,
  },
  productBarcode: {
    fontSize: 12,
    color: '#757575',
  },
  cartContainer: {
    flex: 1,
    backgroundColor: '#fff',
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
    padding: 15,
  },
  cartHeader: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  cartItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  cartItemInfo: {
    flex: 1,
  },
  cartItemName: {
    fontSize: 16,
  },
  cartItemPrice: {
    fontSize: 14,
    color: '#757575',
  },
  quantityControls: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  quantityButton: {
    backgroundColor: '#e0e0e0',
    width: 30,
    height: 30,
    borderRadius: 15,
    justifyContent: 'center',
    alignItems: 'center',
  },
  quantityButtonText: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  quantityText: {
    marginHorizontal: 10,
    fontSize: 16,
  },
  totalContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 20,
    paddingTop: 15,
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
  },
  totalText: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  checkoutButton: {
    backgroundColor: '#2e7d32',
    paddingVertical: 12,
    paddingHorizontal: 20,
    borderRadius: 8,
  },
  checkoutButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});