import React, { useState, useEffect } from 'react';
import { 
  View, 
  Text, 
  StyleSheet, 
  ScrollView, 
  TouchableOpacity, 
  FlatList,
  Alert,
  Share 
} from 'react-native';
import transactionHistoryService, { Transaction, DayReport } from '../../services/TransactionHistoryService';

export default function ReportsScreen() {
  const [activeTab, setActiveTab] = useState<'today' | 'history' | 'weekly'>('today');
  const [todayReport, setTodayReport] = useState<DayReport | null>(null);
  const [weeklyReports, setWeeklyReports] = useState<DayReport[]>([]);
  const [allTransactions, setAllTransactions] = useState<Transaction[]>([]);
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    loadReportData();
  }, [refreshKey]);

  const loadReportData = () => {
    setTodayReport(transactionHistoryService.generateDayReport());
    setWeeklyReports(transactionHistoryService.getWeeklySummary(7));
    setAllTransactions(transactionHistoryService.getAllTransactions().reverse()); // Most recent first
  };

  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
  };

  const handleExportToday = async () => {
    try {
      const today = new Date();
      const csvData = transactionHistoryService.exportToCSV(today, today);
      await Share.share({
        message: csvData,
        title: `POS Report - ${today.toLocaleDateString()}`
      });
    } catch (error) {
      Alert.alert('Export Error', 'Failed to export today\'s report');
    }
  };

  const handleExportWeekly = async () => {
    try {
      const today = new Date();
      const weekAgo = new Date();
      weekAgo.setDate(weekAgo.getDate() - 7);
      const csvData = transactionHistoryService.exportToCSV(weekAgo, today);
      await Share.share({
        message: csvData,
        title: `POS Weekly Report - ${weekAgo.toLocaleDateString()} to ${today.toLocaleDateString()}`
      });
    } catch (error) {
      Alert.alert('Export Error', 'Failed to export weekly report');
    }
  };

  const renderTransaction = ({ item }: { item: Transaction }) => (
    <View style={styles.transactionCard}>
      <View style={styles.transactionHeader}>
        <Text style={styles.transactionId}>{item.id}</Text>
        <Text style={styles.transactionTime}>
          {item.timestamp.toLocaleTimeString()}
        </Text>
      </View>
      
      <View style={styles.transactionItems}>
        {item.items.map((transactionItem, index) => (
          <Text key={index} style={styles.itemText}>
            {transactionItem.quantity}x {transactionItem.name} - ${(transactionItem.quantity * transactionItem.price).toFixed(2)}
          </Text>
        ))}
      </View>
      
      <View style={styles.transactionFooter}>
        <View style={styles.paymentInfo}>
          <Text style={styles.paymentMethod}>{item.paymentMethod}</Text>
          <Text style={[styles.status, { color: getStatusColor(item.status) }]}>
            {item.status.toUpperCase()}
          </Text>
        </View>
        <Text style={styles.transactionTotal}>${item.total.toFixed(2)}</Text>
      </View>
    </View>
  );

  const renderWeeklyReport = ({ item }: { item: DayReport }) => (
    <View style={styles.weeklyCard}>
      <Text style={styles.weeklyDate}>{new Date(item.date).toLocaleDateString()}</Text>
      <View style={styles.weeklyStats}>
        <View style={styles.statItem}>
          <Text style={styles.statLabel}>Transactions</Text>
          <Text style={styles.statValue}>{item.transactionCount}</Text>
        </View>
        <View style={styles.statItem}>
          <Text style={styles.statLabel}>Total Sales</Text>
          <Text style={styles.statValue}>${item.totalGross.toFixed(2)}</Text>
        </View>
      </View>
    </View>
  );

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return '#2e7d32';
      case 'refunded': return '#f57c00';
      case 'void': return '#d32f2f';
      default: return '#666';
    }
  };

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Reports & Analytics</Text>
        <TouchableOpacity onPress={handleRefresh} style={styles.refreshButton}>
          <Text style={styles.refreshText}>↻</Text>
        </TouchableOpacity>
      </View>

      {/* Tab Navigation */}
      <View style={styles.tabContainer}>
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'today' && styles.activeTab]}
          onPress={() => setActiveTab('today')}
        >
          <Text style={[styles.tabText, activeTab === 'today' && styles.activeTabText]}>
            Today
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'weekly' && styles.activeTab]}
          onPress={() => setActiveTab('weekly')}
        >
          <Text style={[styles.tabText, activeTab === 'weekly' && styles.activeTabText]}>
            Weekly
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'history' && styles.activeTab]}
          onPress={() => setActiveTab('history')}
        >
          <Text style={[styles.tabText, activeTab === 'history' && styles.activeTabText]}>
            History
          </Text>
        </TouchableOpacity>
      </View>

      {/* Content */}
      <ScrollView style={styles.content}>
        {activeTab === 'today' && todayReport && (
          <View style={styles.todayReport}>
            <View style={styles.reportHeader}>
              <Text style={styles.reportTitle}>Today's Report</Text>
              <TouchableOpacity onPress={handleExportToday} style={styles.exportButton}>
                <Text style={styles.exportText}>Export</Text>
              </TouchableOpacity>
            </View>
            
            <View style={styles.statsGrid}>
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>{todayReport.transactionCount}</Text>
                <Text style={styles.statLabel}>Transactions</Text>
              </View>
              
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>${todayReport.totalGross.toFixed(2)}</Text>
                <Text style={styles.statLabel}>Total Sales</Text>
              </View>
              
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>${todayReport.cashSales.toFixed(2)}</Text>
                <Text style={styles.statLabel}>Cash Sales</Text>
              </View>
              
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>${todayReport.cardSales.toFixed(2)}</Text>
                <Text style={styles.statLabel}>Card Sales</Text>
              </View>
              
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>${todayReport.totalTax.toFixed(2)}</Text>
                <Text style={styles.statLabel}>Tax Collected</Text>
              </View>
              
              <View style={styles.statCard}>
                <Text style={styles.statNumber}>
                  ${(todayReport.transactionCount > 0 ? todayReport.totalGross / todayReport.transactionCount : 0).toFixed(2)}
                </Text>
                <Text style={styles.statLabel}>Avg. Sale</Text>
              </View>
            </View>

            {todayReport.transactions.length > 0 && (
              <View style={styles.transactionsList}>
                <Text style={styles.sectionTitle}>Today's Transactions</Text>
                <FlatList
                  data={todayReport.transactions}
                  renderItem={renderTransaction}
                  keyExtractor={(item) => item.id}
                  scrollEnabled={false}
                />
              </View>
            )}
          </View>
        )}

        {activeTab === 'weekly' && (
          <View style={styles.weeklyReport}>
            <View style={styles.reportHeader}>
              <Text style={styles.reportTitle}>Weekly Overview</Text>
              <TouchableOpacity onPress={handleExportWeekly} style={styles.exportButton}>
                <Text style={styles.exportText}>Export</Text>
              </TouchableOpacity>
            </View>
            
            <FlatList
              data={weeklyReports}
              renderItem={renderWeeklyReport}
              keyExtractor={(item) => item.date}
              scrollEnabled={false}
            />
          </View>
        )}

        {activeTab === 'history' && (
          <View style={styles.historyReport}>
            <Text style={styles.reportTitle}>Transaction History</Text>
            
            {allTransactions.length === 0 ? (
              <View style={styles.emptyState}>
                <Text style={styles.emptyText}>No transactions yet</Text>
              </View>
            ) : (
              <FlatList
                data={allTransactions}
                renderItem={renderTransaction}
                keyExtractor={(item) => item.id}
                scrollEnabled={false}
              />
            )}
          </View>
        )}
      </ScrollView>
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
  refreshButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(255,255,255,0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  refreshText: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
  },
  tabContainer: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  tab: {
    flex: 1,
    paddingVertical: 15,
    alignItems: 'center',
  },
  activeTab: {
    borderBottomWidth: 2,
    borderBottomColor: '#2e7d32',
  },
  tabText: {
    fontSize: 16,
    color: '#666',
  },
  activeTabText: {
    color: '#2e7d32',
    fontWeight: 'bold',
  },
  content: {
    flex: 1,
  },
  todayReport: {
    padding: 20,
  },
  reportHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 20,
  },
  reportTitle: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  exportButton: {
    backgroundColor: '#2e7d32',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 6,
  },
  exportText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: 'bold',
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
  },
  statCard: {
    width: '48%',
    backgroundColor: '#fff',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10,
    alignItems: 'center',
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  statNumber: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2e7d32',
    marginBottom: 5,
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
  },
  transactionsList: {
    marginTop: 20,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  transactionCard: {
    backgroundColor: '#fff',
    padding: 15,
    marginBottom: 10,
    borderRadius: 8,
    elevation: 1,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
  },
  transactionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 10,
  },
  transactionId: {
    fontSize: 14,
    fontWeight: 'bold',
  },
  transactionTime: {
    fontSize: 12,
    color: '#666',
  },
  transactionItems: {
    marginBottom: 10,
  },
  itemText: {
    fontSize: 12,
    color: '#666',
    marginBottom: 2,
  },
  transactionFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  paymentInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  paymentMethod: {
    fontSize: 12,
    color: '#666',
    marginRight: 10,
  },
  status: {
    fontSize: 10,
    fontWeight: 'bold',
  },
  transactionTotal: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2e7d32',
  },
  weeklyReport: {
    padding: 20,
  },
  weeklyCard: {
    backgroundColor: '#fff',
    padding: 15,
    marginBottom: 10,
    borderRadius: 8,
    elevation: 1,
  },
  weeklyDate: {
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  weeklyStats: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2e7d32',
  },
  historyReport: {
    padding: 20,
  },
  emptyState: {
    alignItems: 'center',
    padding: 40,
  },
  emptyText: {
    fontSize: 16,
    color: '#666',
  },
});