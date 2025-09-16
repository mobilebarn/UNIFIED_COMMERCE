import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, TextInput, Alert } from 'react-native';
import { useMutation, useQuery } from '@apollo/client';
import { GET_CART_BY_SESSION, CREATE_PAYMENT, CAPTURE_PAYMENT } from '../../graphql';
import { localStorage } from '../../utils/storage';

export default function CheckoutScreen() {
  const [paymentMethod, setPaymentMethod] = useState('card');
  const [cardNumber, setCardNumber] = useState('');
  const [expiryDate, setExpiryDate] = useState('');
  const [cvv, setCvv] = useState('');
  const [cashAmount, setCashAmount] = useState('');
  const [processing, setProcessing] = useState(false);

  // Get cart data
  const sessionId = localStorage.getItem('posSessionId') || 'session_' + Date.now();
  const { data: cartData, loading: cartLoading, error: cartError } = useQuery(
    GET_CART_BY_SESSION,
    {
      variables: { sessionId },
      fetchPolicy: 'network-only'
    }
  );

  // Payment mutations
  const [createPayment] = useMutation(CREATE_PAYMENT);
  const [capturePayment] = useMutation(CAPTURE_PAYMENT);

  if (cartLoading) {
    return (
      <View style={styles.container}>
        <Text>Loading cart data...</Text>
      </View>
    );
  }

  if (cartError) {
    return (
      <View style={styles.container}>
        <Text>Error loading cart: {cartError.message}</Text>
      </View>
    );
  }

  const cart = cartData?.cartBySession;
  const cartItems = cart?.lineItems || [];

  // Calculate subtotal
  const subtotal = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
  
  // Calculate tax (assuming 8.5% tax rate)
  const tax = subtotal * 0.085;
  
  // Calculate total
  const total = subtotal + tax;

  // Handle payment processing
  const handlePayment = async () => {
    if (!cart) {
      Alert.alert('Error', 'No cart found');
      return;
    }

    setProcessing(true);

    try {
      // Create payment record
      const { data: paymentData } = await createPayment({
        variables: {
          input: {
            merchantId: 'merchant_1', // In a real app, this would come from auth context
            amount: total,
            currency: 'USD',
            paymentMethodId: 'pm_card', // This would be the actual payment method ID
            description: `POS payment for cart ${cart.id}`
          }
        }
      });

      const paymentId = paymentData.createPayment.id;

      // Capture the payment (in a real app, this would connect to a payment processor)
      await capturePayment({
        variables: {
          id: paymentId,
          amount: total
        }
      });

      // Show success message
      Alert.alert(
        'Payment Successful',
        `Payment processed successfully! Total: $${total.toFixed(2)}`,
        [
          {
            text: 'OK',
            onPress: () => {
              // Clear cart and navigate to receipt or new transaction
              // This would be implemented based on your navigation setup
            }
          }
        ]
      );
    } catch (error) {
      Alert.alert('Payment Error', 'Failed to process payment. Please try again.');
      console.error(error);
    } finally {
      setProcessing(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.header}>Checkout</Text>
      
      {/* Order Summary */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Order Summary</Text>
        {cartItems.map(item => (
          <View key={item.id} style={styles.cartItem}>
            <View style={styles.itemInfo}>
              <Text style={styles.itemName}>{item.name}</Text>
              <Text style={styles.itemQuantity}>Qty: {item.quantity}</Text>
            </View>
            <Text style={styles.itemPrice}>${(item.price * item.quantity).toFixed(2)}</Text>
          </View>
        ))}
        
        <View style={styles.orderTotal}>
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Subtotal</Text>
            <Text style={styles.totalValue}>${subtotal.toFixed(2)}</Text>
          </View>
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Tax</Text>
            <Text style={styles.totalValue}>${tax.toFixed(2)}</Text>
          </View>
          <View style={[styles.totalRow, styles.grandTotalRow]}>
            <Text style={styles.grandTotalLabel}>Total</Text>
            <Text style={styles.grandTotalValue}>${total.toFixed(2)}</Text>
          </View>
        </View>
      </View>
      
      {/* Payment Method */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Payment Method</Text>
        
        <View style={styles.paymentOptions}>
          <TouchableOpacity 
            style={[styles.paymentOption, paymentMethod === 'card' && styles.selectedPaymentOption]}
            onPress={() => setPaymentMethod('card')}
          >
            <Text style={styles.paymentOptionText}>Card</Text>
          </TouchableOpacity>
          
          <TouchableOpacity 
            style={[styles.paymentOption, paymentMethod === 'cash' && styles.selectedPaymentOption]}
            onPress={() => setPaymentMethod('cash')}
          >
            <Text style={styles.paymentOptionText}>Cash</Text>
          </TouchableOpacity>
          
          <TouchableOpacity 
            style={[styles.paymentOption, paymentMethod === 'other' && styles.selectedPaymentOption]}
            onPress={() => setPaymentMethod('other')}
          >
            <Text style={styles.paymentOptionText}>Other</Text>
          </TouchableOpacity>
        </View>
        
        {paymentMethod === 'card' && (
          <View style={styles.cardForm}>
            <TextInput
              style={styles.input}
              placeholder="Card Number"
              value={cardNumber}
              onChangeText={setCardNumber}
              keyboardType="numeric"
            />
            <View style={styles.cardDetails}>
              <TextInput
                style={[styles.input, styles.halfInput]}
                placeholder="MM/YY"
                value={expiryDate}
                onChangeText={setExpiryDate}
                keyboardType="numeric"
              />
              <TextInput
                style={[styles.input, styles.halfInput]}
                placeholder="CVV"
                value={cvv}
                onChangeText={setCvv}
                keyboardType="numeric"
                secureTextEntry
              />
            </View>
          </View>
        )}
        
        {paymentMethod === 'cash' && (
          <View style={styles.cashForm}>
            <TextInput
              style={styles.input}
              placeholder="Cash Amount"
              value={cashAmount}
              onChangeText={setCashAmount}
              keyboardType="numeric"
            />
            {cashAmount !== '' && (
              <View style={styles.changeRow}>
                <Text style={styles.changeLabel}>Change</Text>
                <Text style={styles.changeValue}>
                  ${(parseFloat(cashAmount) - total).toFixed(2)}
                </Text>
              </View>
            )}
          </View>
        )}
      </View>
      
      {/* Process Payment Button */}
      <TouchableOpacity 
        style={[styles.processButton, processing && styles.disabledButton]} 
        onPress={handlePayment}
        disabled={processing}
      >
        <Text style={styles.processButtonText}>
          {processing ? 'Processing...' : `Process Payment - $${total.toFixed(2)}`}
        </Text>
      </TouchableOpacity>
    </ScrollView>
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
  section: {
    backgroundColor: '#fff',
    marginVertical: 10,
    padding: 15,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  cartItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  itemInfo: {
    flex: 1,
  },
  itemName: {
    fontSize: 16,
  },
  itemQuantity: {
    fontSize: 14,
    color: '#757575',
  },
  itemPrice: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  orderTotal: {
    marginTop: 15,
  },
  totalRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 5,
  },
  totalLabel: {
    fontSize: 16,
  },
  totalValue: {
    fontSize: 16,
  },
  grandTotalRow: {
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
    marginTop: 10,
    paddingTop: 10,
  },
  grandTotalLabel: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  grandTotalValue: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#2e7d32',
  },
  paymentOptions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 20,
  },
  paymentOption: {
    flex: 1,
    padding: 15,
    marginHorizontal: 5,
    borderWidth: 1,
    borderColor: '#e0e0e0',
    borderRadius: 8,
    alignItems: 'center',
  },
  selectedPaymentOption: {
    borderColor: '#2e7d32',
    backgroundColor: '#e8f5e9',
  },
  paymentOptionText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  cardForm: {
    marginTop: 10,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    marginBottom: 15,
  },
  cardDetails: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  halfInput: {
    flex: 1,
    marginHorizontal: 5,
  },
  cashForm: {
    marginTop: 10,
  },
  changeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 10,
    padding: 10,
    backgroundColor: '#e8f5e9',
    borderRadius: 8,
  },
  changeLabel: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  changeValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2e7d32',
  },
  processButton: {
    backgroundColor: '#2e7d32',
    padding: 20,
    margin: 20,
    borderRadius: 8,
    alignItems: 'center',
  },
  disabledButton: {
    backgroundColor: '#a5d6a7',
  },
  processButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
  },
});