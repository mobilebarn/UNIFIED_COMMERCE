import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, TextInput, Alert, Modal } from 'react-native';
import { useMutation, useQuery } from '@apollo/client';
import { GET_CART_BY_SESSION, CREATE_PAYMENT, CAPTURE_PAYMENT } from '../../graphql';
import { localStorage } from '../../utils/storage';
import CardReaderScreen from '../../components/CardReaderScreen';
import ReceiptGenerator, { ReceiptData } from '../../components/ReceiptGenerator';
import stripeTerminalService, { ReaderInfo } from '../../services/StripeTerminalService';
import transactionHistoryService, { Transaction } from '../../services/TransactionHistoryService';

export default function CheckoutScreen() {
  const [paymentMethod, setPaymentMethod] = useState('card');
  const [cardNumber, setCardNumber] = useState('');
  const [expiryDate, setExpiryDate] = useState('');
  const [cvv, setCvv] = useState('');
  const [cashAmount, setCashAmount] = useState('');
  const [processing, setProcessing] = useState(false);
  const [showCardReader, setShowCardReader] = useState(false);
  const [connectedReader, setConnectedReader] = useState<ReaderInfo | null>(null);
  const [showReceipt, setShowReceipt] = useState(false);
  const [receiptData, setReceiptData] = useState<ReceiptData | null>(null);

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
  const subtotal = cartItems.reduce((sum: number, item: any) => sum + item.price * item.quantity, 0);
  
  // Calculate tax (assuming 8.5% tax rate)
  const tax = subtotal * 0.085;
  
  // Calculate total
  const total = subtotal + tax;

  // Generate receipt data
  const generateReceiptData = (paymentMethod: string, transactionId?: string, cashReceived?: number): ReceiptData => {
    return {
      id: 'RCP_' + Date.now(),
      timestamp: new Date(),
      items: cartItems.map((item: any) => ({
        id: item.id,
        name: item.name,
        quantity: item.quantity,
        price: item.price
      })),
      subtotal,
      tax,
      total,
      paymentMethod,
      transactionId,
      cashReceived,
      change: cashReceived ? cashReceived - total : undefined
    };
  };

  // Save transaction and generate receipt
  const completeTransaction = (paymentMethod: string, transactionId?: string, cashReceived?: number) => {
    const transaction: Transaction = {
      id: 'TXN_' + Date.now(),
      timestamp: new Date(),
      items: cartItems.map((item: any) => ({
        id: item.id,
        name: item.name,
        quantity: item.quantity,
        price: item.price
      })),
      subtotal,
      tax,
      total,
      paymentMethod,
      transactionId,
      cashReceived,
      change: cashReceived ? cashReceived - total : undefined,
      status: 'completed'
    };
    
    // Save transaction to history
    transactionHistoryService.addTransaction(transaction);
    
    // Generate receipt
    const receipt = generateReceiptData(paymentMethod, transactionId, cashReceived);
    setReceiptData(receipt);
    setShowReceipt(true);
  };

  // Handle card payment with Stripe Terminal
  const handleCardPayment = async () => {
    if (!connectedReader) {
      setShowCardReader(true);
      return;
    }

    setProcessing(true);
    try {
      const result = await stripeTerminalService.processPayment(total);
      
      if (result.success) {
        completeTransaction('Card', result.paymentIntent?.id);
      } else {
        Alert.alert('Payment Failed', result.error || 'Card payment failed. Please try again.');
      }
    } catch (error) {
      Alert.alert('Payment Error', 'Failed to process card payment. Please try again.');
    } finally {
      setProcessing(false);
    }
  };

  // Handle cash payment
  const handleCashPayment = async () => {
    const cashAmountValue = parseFloat(cashAmount);
    if (isNaN(cashAmountValue) || cashAmountValue < total) {
      Alert.alert('Invalid Amount', 'Please enter a valid cash amount that covers the total.');
      return;
    }

    const change = cashAmountValue - total;
    completeTransaction('Cash', undefined, cashAmountValue);
  };

  // Handle payment processing based on selected method
  const handlePayment = async () => {
    if (!cart) {
      Alert.alert('Error', 'No cart found');
      return;
    }

    if (paymentMethod === 'card') {
      await handleCardPayment();
    } else if (paymentMethod === 'cash') {
      await handleCashPayment();
    } else {
      Alert.alert('Payment Method', 'Please select a payment method');
    }
  };

  const onReaderConnected = (readerInfo: ReaderInfo) => {
    setConnectedReader(readerInfo);
    setShowCardReader(false);
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.header}>Checkout</Text>
      
      {/* Order Summary */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Order Summary</Text>
        {cartItems.map((item: any) => (
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
            {connectedReader ? (
              <View style={styles.readerStatus}>
                <Text style={styles.readerStatusText}>
                  ✓ Connected: {connectedReader.label || connectedReader.serialNumber}
                </Text>
                <TouchableOpacity 
                  style={styles.changeReaderButton}
                  onPress={() => setShowCardReader(true)}
                >
                  <Text style={styles.changeReaderText}>Change Reader</Text>
                </TouchableOpacity>
              </View>
            ) : (
              <TouchableOpacity 
                style={styles.setupReaderButton}
                onPress={() => setShowCardReader(true)}
              >
                <Text style={styles.setupReaderText}>Set Up Card Reader</Text>
              </TouchableOpacity>
            )}
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
      
      {/* Card Reader Modal */}
      <Modal
        visible={showCardReader}
        animationType="slide"
        presentationStyle="pageSheet"
        onRequestClose={() => setShowCardReader(false)}
      >
        <CardReaderScreen
          onReaderConnected={onReaderConnected}
          onClose={() => setShowCardReader(false)}
        />
      </Modal>
      
      {/* Receipt Modal */}
      <Modal
        visible={showReceipt}
        animationType="slide"
        presentationStyle="pageSheet"
        onRequestClose={() => setShowReceipt(false)}
      >
        {receiptData && (
          <ReceiptGenerator
            receiptData={receiptData}
            onClose={() => {
              setShowReceipt(false);
              setReceiptData(null);
              // Clear cart and reset session
              // This would be implemented based on your navigation setup
            }}
            onPrint={() => {
              console.log('Print receipt');
            }}
            onShare={() => {
              console.log('Share receipt');
            }}
          />
        )}
      </Modal>
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
  readerStatus: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 15,
    backgroundColor: '#e8f5e9',
    borderRadius: 8,
    marginBottom: 10,
  },
  readerStatusText: {
    fontSize: 16,
    color: '#2e7d32',
    fontWeight: 'bold',
  },
  changeReaderButton: {
    backgroundColor: '#fff',
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: '#2e7d32',
  },
  changeReaderText: {
    color: '#2e7d32',
    fontSize: 14,
    fontWeight: 'bold',
  },
  setupReaderButton: {
    backgroundColor: '#2e7d32',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
    marginBottom: 10,
  },
  setupReaderText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});