import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Share, Alert } from 'react-native';

export interface ReceiptItem {
  id: string;
  name: string;
  quantity: number;
  price: number;
}

export interface ReceiptData {
  id: string;
  timestamp: Date;
  items: ReceiptItem[];
  subtotal: number;
  tax: number;
  total: number;
  paymentMethod: string;
  transactionId?: string;
  cashReceived?: number;
  change?: number;
}

interface ReceiptGeneratorProps {
  receiptData: ReceiptData;
  onPrint?: () => void;
  onShare?: () => void;
  onClose?: () => void;
}

export default function ReceiptGenerator({ 
  receiptData, 
  onPrint, 
  onShare, 
  onClose 
}: ReceiptGeneratorProps) {
  
  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const generateReceiptText = () => {
    let receipt = `
========================================
              RETAIL OS POS
            Point of Sale System
========================================

Receipt #: ${receiptData.id}
Date: ${formatDate(receiptData.timestamp)}

----------------------------------------
ITEMS
----------------------------------------
`;

    receiptData.items.forEach(item => {
      const lineTotal = item.price * item.quantity;
      receipt += `${item.name}\n`;
      receipt += `  ${item.quantity} x $${item.price.toFixed(2)} = $${lineTotal.toFixed(2)}\n\n`;
    });

    receipt += `----------------------------------------
TOTALS
----------------------------------------
Subtotal:            $${receiptData.subtotal.toFixed(2)}
Tax:                 $${receiptData.tax.toFixed(2)}
TOTAL:               $${receiptData.total.toFixed(2)}

Payment Method: ${receiptData.paymentMethod}`;

    if (receiptData.transactionId) {
      receipt += `\nTransaction ID: ${receiptData.transactionId}`;
    }

    if (receiptData.cashReceived) {
      receipt += `\nCash Received: $${receiptData.cashReceived.toFixed(2)}`;
      receipt += `\nChange: $${receiptData.change?.toFixed(2) || '0.00'}`;
    }

    receipt += `

========================================
        Thank you for your business!
        
        Powered by Unified Commerce OS
========================================
`;

    return receipt;
  };

  const handleShare = async () => {
    try {
      const receiptText = generateReceiptText();
      await Share.share({
        message: receiptText,
        title: `Receipt #${receiptData.id}`
      });
      onShare?.();
    } catch (error) {
      Alert.alert('Share Error', 'Failed to share receipt');
    }
  };

  const handlePrint = () => {
    // In a real implementation, this would integrate with a printer SDK
    Alert.alert(
      'Print Receipt', 
      'This would normally send the receipt to a connected printer. For now, you can share the receipt instead.',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Share Instead', onPress: handleShare }
      ]
    );
    onPrint?.();
  };

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Receipt Generated</Text>
        {onClose && (
          <TouchableOpacity onPress={onClose} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>×</Text>
          </TouchableOpacity>
        )}
      </View>

      <View style={styles.receipt}>
        <Text style={styles.receiptHeader}>UNIFIED COMMERCE OS POS</Text>
        <Text style={styles.receiptSubheader}>Point of Sale System</Text>
        
        <View style={styles.divider} />
        
        <View style={styles.receiptInfo}>
          <Text style={styles.receiptId}>Receipt #: {receiptData.id}</Text>
          <Text style={styles.receiptDate}>{formatDate(receiptData.timestamp)}</Text>
        </View>
        
        <View style={styles.divider} />
        
        <Text style={styles.sectionTitle}>ITEMS</Text>
        {receiptData.items.map((item) => (
          <View key={item.id} style={styles.receiptItem}>
            <Text style={styles.itemName}>{item.name}</Text>
            <View style={styles.itemDetails}>
              <Text style={styles.itemQuantity}>
                {item.quantity} × ${item.price.toFixed(2)}
              </Text>
              <Text style={styles.itemTotal}>
                ${(item.price * item.quantity).toFixed(2)}
              </Text>
            </View>
          </View>
        ))}
        
        <View style={styles.divider} />
        
        <View style={styles.totals}>
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Subtotal</Text>
            <Text style={styles.totalValue}>${receiptData.subtotal.toFixed(2)}</Text>
          </View>
          <View style={styles.totalRow}>
            <Text style={styles.totalLabel}>Tax</Text>
            <Text style={styles.totalValue}>${receiptData.tax.toFixed(2)}</Text>
          </View>
          <View style={[styles.totalRow, styles.grandTotal]}>
            <Text style={styles.grandTotalLabel}>TOTAL</Text>
            <Text style={styles.grandTotalValue}>${receiptData.total.toFixed(2)}</Text>
          </View>
        </View>
        
        <View style={styles.divider} />
        
        <View style={styles.paymentInfo}>
          <Text style={styles.paymentMethod}>Payment: {receiptData.paymentMethod}</Text>
          {receiptData.transactionId && (
            <Text style={styles.transactionId}>ID: {receiptData.transactionId}</Text>
          )}
          {receiptData.cashReceived && (
            <>
              <Text style={styles.cashReceived}>
                Cash: ${receiptData.cashReceived.toFixed(2)}
              </Text>
              <Text style={styles.change}>
                Change: ${receiptData.change?.toFixed(2) || '0.00'}
              </Text>
            </>
          )}
        </View>
        
        <View style={styles.divider} />
        
        <Text style={styles.footer}>Thank you for your business!</Text>
        <Text style={styles.poweredBy}>Powered by Unified Commerce OS</Text>
      </View>

      <View style={styles.actions}>
        <TouchableOpacity style={styles.printButton} onPress={handlePrint}>
          <Text style={styles.buttonText}>Print Receipt</Text>
        </TouchableOpacity>
        
        <TouchableOpacity style={styles.shareButton} onPress={handleShare}>
          <Text style={styles.buttonText}>Share Receipt</Text>
        </TouchableOpacity>
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
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#2e7d32',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#fff',
  },
  closeButton: {
    width: 30,
    height: 30,
    borderRadius: 15,
    backgroundColor: 'rgba(255,255,255,0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  closeButtonText: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
  },
  receipt: {
    flex: 1,
    backgroundColor: '#fff',
    margin: 20,
    padding: 20,
    borderRadius: 8,
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  receiptHeader: {
    fontSize: 18,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 5,
  },
  receiptSubheader: {
    fontSize: 14,
    textAlign: 'center',
    color: '#666',
    marginBottom: 15,
  },
  divider: {
    height: 1,
    backgroundColor: '#ddd',
    marginVertical: 15,
  },
  receiptInfo: {
    marginBottom: 10,
  },
  receiptId: {
    fontSize: 14,
    fontWeight: 'bold',
  },
  receiptDate: {
    fontSize: 14,
    color: '#666',
    marginTop: 2,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  receiptItem: {
    marginBottom: 10,
  },
  itemName: {
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 2,
  },
  itemDetails: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  itemQuantity: {
    fontSize: 12,
    color: '#666',
  },
  itemTotal: {
    fontSize: 12,
    fontWeight: 'bold',
  },
  totals: {
    marginTop: 10,
  },
  totalRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 5,
  },
  totalLabel: {
    fontSize: 14,
  },
  totalValue: {
    fontSize: 14,
  },
  grandTotal: {
    borderTopWidth: 1,
    borderTopColor: '#ddd',
    paddingTop: 5,
    marginTop: 5,
  },
  grandTotalLabel: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  grandTotalValue: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  paymentInfo: {
    marginTop: 10,
  },
  paymentMethod: {
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 2,
  },
  transactionId: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
  },
  cashReceived: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
  },
  change: {
    fontSize: 12,
    color: '#666',
  },
  footer: {
    fontSize: 14,
    textAlign: 'center',
    fontWeight: 'bold',
    marginTop: 10,
  },
  poweredBy: {
    fontSize: 12,
    textAlign: 'center',
    color: '#666',
    marginTop: 5,
  },
  actions: {
    flexDirection: 'row',
    padding: 20,
    gap: 10,
  },
  printButton: {
    flex: 1,
    backgroundColor: '#2e7d32',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
  },
  shareButton: {
    flex: 1,
    backgroundColor: '#1976d2',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});