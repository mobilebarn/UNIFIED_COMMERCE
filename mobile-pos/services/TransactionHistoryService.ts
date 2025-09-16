export interface Transaction {
  id: string;
  timestamp: Date;
  items: TransactionItem[];
  subtotal: number;
  tax: number;
  total: number;
  paymentMethod: string;
  transactionId?: string;
  cashReceived?: number;
  change?: number;
  status: 'completed' | 'refunded' | 'void';
}

export interface TransactionItem {
  id: string;
  name: string;
  quantity: number;
  price: number;
}

export interface DayReport {
  date: string;
  transactions: Transaction[];
  totalSales: number;
  totalTax: number;
  totalGross: number;
  cashSales: number;
  cardSales: number;
  transactionCount: number;
}

class TransactionHistoryService {
  private transactions: Transaction[] = [];

  // Add a new transaction
  addTransaction(transaction: Transaction): void {
    this.transactions.push(transaction);
    this.saveToStorage();
  }

  // Get all transactions
  getAllTransactions(): Transaction[] {
    return [...this.transactions];
  }

  // Get transactions for a specific date
  getTransactionsByDate(date: string): Transaction[] {
    return this.transactions.filter(transaction => 
      transaction.timestamp.toDateString() === new Date(date).toDateString()
    );
  }

  // Get transactions for today
  getTodaysTransactions(): Transaction[] {
    const today = new Date().toDateString();
    return this.transactions.filter(transaction => 
      transaction.timestamp.toDateString() === today
    );
  }

  // Generate end-of-day report
  generateDayReport(date?: string): DayReport {
    const reportDate = date || new Date().toDateString();
    const dayTransactions = this.getTransactionsByDate(reportDate);
    
    const totalSales = dayTransactions.reduce((sum, t) => sum + t.subtotal, 0);
    const totalTax = dayTransactions.reduce((sum, t) => sum + t.tax, 0);
    const totalGross = dayTransactions.reduce((sum, t) => sum + t.total, 0);
    
    const cashSales = dayTransactions
      .filter(t => t.paymentMethod === 'Cash')
      .reduce((sum, t) => sum + t.total, 0);
    
    const cardSales = dayTransactions
      .filter(t => t.paymentMethod === 'Card')
      .reduce((sum, t) => sum + t.total, 0);

    return {
      date: reportDate,
      transactions: dayTransactions,
      totalSales,
      totalTax,
      totalGross,
      cashSales,
      cardSales,
      transactionCount: dayTransactions.length
    };
  }

  // Get summary for last N days
  getWeeklySummary(days: number = 7): DayReport[] {
    const reports: DayReport[] = [];
    const today = new Date();
    
    for (let i = 0; i < days; i++) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);
      reports.push(this.generateDayReport(date.toDateString()));
    }
    
    return reports.reverse(); // Most recent first
  }

  // Update transaction status (for refunds/voids)
  updateTransactionStatus(transactionId: string, status: 'completed' | 'refunded' | 'void'): boolean {
    const transaction = this.transactions.find(t => t.id === transactionId);
    if (transaction) {
      transaction.status = status;
      this.saveToStorage();
      return true;
    }
    return false;
  }

  // Search transactions
  searchTransactions(query: string): Transaction[] {
    const lowerQuery = query.toLowerCase();
    return this.transactions.filter(transaction => 
      transaction.id.toLowerCase().includes(lowerQuery) ||
      transaction.transactionId?.toLowerCase().includes(lowerQuery) ||
      transaction.items.some(item => 
        item.name.toLowerCase().includes(lowerQuery)
      )
    );
  }

  // Get best selling items
  getBestSellingItems(days: number = 30): { name: string; quantity: number; revenue: number }[] {
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - days);
    
    const recentTransactions = this.transactions.filter(t => t.timestamp >= cutoffDate);
    const itemStats: { [key: string]: { quantity: number; revenue: number } } = {};
    
    recentTransactions.forEach(transaction => {
      transaction.items.forEach(item => {
        if (!itemStats[item.name]) {
          itemStats[item.name] = { quantity: 0, revenue: 0 };
        }
        itemStats[item.name].quantity += item.quantity;
        itemStats[item.name].revenue += item.quantity * item.price;
      });
    });
    
    return Object.entries(itemStats)
      .map(([name, stats]) => ({ name, ...stats }))
      .sort((a, b) => b.quantity - a.quantity);
  }

  // Clear all transactions (for testing)
  clearAllTransactions(): void {
    this.transactions = [];
    this.saveToStorage();
  }

  // Save transactions to local storage
  private saveToStorage(): void {
    try {
      if (typeof Storage !== 'undefined') {
        localStorage.setItem('pos_transactions', JSON.stringify(this.transactions));
      }
    } catch (error) {
      console.error('Failed to save transactions to storage:', error);
    }
  }

  // Load transactions from local storage
  loadFromStorage(): void {
    try {
      if (typeof Storage !== 'undefined') {
        const stored = localStorage.getItem('pos_transactions');
        if (stored) {
          const parsed = JSON.parse(stored);
          // Convert timestamp strings back to Date objects
          this.transactions = parsed.map((t: any) => ({
            ...t,
            timestamp: new Date(t.timestamp)
          }));
        }
      }
    } catch (error) {
      console.error('Failed to load transactions from storage:', error);
      this.transactions = [];
    }
  }

  // Export transactions to CSV format
  exportToCSV(startDate?: Date, endDate?: Date): string {
    let transactions = this.transactions;
    
    if (startDate || endDate) {
      transactions = this.transactions.filter(t => {
        if (startDate && t.timestamp < startDate) return false;
        if (endDate && t.timestamp > endDate) return false;
        return true;
      });
    }
    
    const headers = [
      'Transaction ID',
      'Date',
      'Time',
      'Items',
      'Subtotal',
      'Tax',
      'Total',
      'Payment Method',
      'Payment ID',
      'Status'
    ].join(',');
    
    const rows = transactions.map(t => [
      t.id,
      t.timestamp.toLocaleDateString(),
      t.timestamp.toLocaleTimeString(),
      t.items.map(item => `${item.name} (${item.quantity})`).join('; '),
      t.subtotal.toFixed(2),
      t.tax.toFixed(2),
      t.total.toFixed(2),
      t.paymentMethod,
      t.transactionId || '',
      t.status
    ].join(','));
    
    return [headers, ...rows].join('\n');
  }
}

// Create singleton instance
const transactionHistoryService = new TransactionHistoryService();

// Load existing transactions on initialization
transactionHistoryService.loadFromStorage();

export default transactionHistoryService;