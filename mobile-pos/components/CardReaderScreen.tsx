import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  Alert,
  ActivityIndicator,
} from 'react-native';
import stripeTerminalService, { ReaderInfo, Reader } from '../services/StripeTerminalService';

interface CardReaderScreenProps {
  onReaderConnected: (readerInfo: ReaderInfo) => void;
  onClose: () => void;
}

export default function CardReaderScreen({ onReaderConnected, onClose }: CardReaderScreenProps) {
  const [isInitializing, setIsInitializing] = useState(false);
  const [isDiscovering, setIsDiscovering] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);
  const [readers, setReaders] = useState<Reader[]>([]);
  const [connectedReader, setConnectedReader] = useState<ReaderInfo | null>(null);

  useEffect(() => {
    initializeTerminal();
    
    // Check if already connected
    const connected = stripeTerminalService.getConnectedReader();
    if (connected) {
      setConnectedReader(connected);
    }
  }, []);

  const initializeTerminal = async () => {
    setIsInitializing(true);
    try {
      const initialized = await stripeTerminalService.initialize();
      if (!initialized) {
        Alert.alert('Error', 'Failed to initialize Stripe Terminal');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to initialize Stripe Terminal');
    } finally {
      setIsInitializing(false);
    }
  };

  const discoverReaders = async () => {
    setIsDiscovering(true);
    try {
      const discoveredReaders = await stripeTerminalService.discoverReaders('bluetooth');
      setReaders(discoveredReaders);
      
      if (discoveredReaders.length === 0) {
        Alert.alert(
          'No Readers Found',
          'Make sure your card reader is turned on and in pairing mode.'
        );
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to discover card readers');
    } finally {
      setIsDiscovering(false);
    }
  };

  const connectToReader = async (reader: Reader) => {
    setIsConnecting(true);
    try {
      const connected = await stripeTerminalService.connectReader(reader);
      if (connected) {
        const readerInfo = stripeTerminalService.getConnectedReader();
        if (readerInfo) {
          setConnectedReader(readerInfo);
          onReaderConnected(readerInfo);
          Alert.alert('Success', `Connected to ${readerInfo.label || readerInfo.serialNumber}`);
        }
      } else {
        Alert.alert('Error', 'Failed to connect to card reader');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to connect to card reader');
    } finally {
      setIsConnecting(false);
    }
  };

  const disconnectReader = async () => {
    try {
      await stripeTerminalService.disconnectReader();
      setConnectedReader(null);
      Alert.alert('Success', 'Disconnected from card reader');
    } catch (error) {
      Alert.alert('Error', 'Failed to disconnect from card reader');
    }
  };

  const renderReader = ({ item }: { item: Reader }) => (
    <TouchableOpacity
      style={styles.readerItem}
      onPress={() => connectToReader(item)}
      disabled={isConnecting}
    >
      <View style={styles.readerInfo}>
        <Text style={styles.readerName}>{item.label || 'Unknown Reader'}</Text>
        <Text style={styles.readerSerial}>Serial: {item.serialNumber}</Text>
        <Text style={styles.readerType}>Type: {item.deviceType}</Text>
        {item.batteryLevel && (
          <Text style={styles.readerBattery}>
            Battery: {Math.round(item.batteryLevel * 100)}%
            {item.isCharging ? ' (Charging)' : ''}
          </Text>
        )}
      </View>
      <View style={styles.connectButton}>
        <Text style={styles.connectButtonText}>Connect</Text>
      </View>
    </TouchableOpacity>
  );

  if (isInitializing) {
    return (
      <View style={styles.container}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#2e7d32" />
          <Text style={styles.loadingText}>Initializing Stripe Terminal...</Text>
        </View>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Card Reader Setup</Text>
        <TouchableOpacity style={styles.closeButton} onPress={onClose}>
          <Text style={styles.closeButtonText}>✕</Text>
        </TouchableOpacity>
      </View>

      {connectedReader ? (
        <View style={styles.connectedSection}>
          <Text style={styles.sectionTitle}>Connected Reader</Text>
          <View style={styles.connectedReader}>
            <Text style={styles.connectedReaderName}>
              {connectedReader.label || 'Card Reader'}
            </Text>
            <Text style={styles.connectedReaderSerial}>
              Serial: {connectedReader.serialNumber}
            </Text>
            {connectedReader.batteryLevel && (
              <Text style={styles.connectedReaderBattery}>
                Battery: {Math.round(connectedReader.batteryLevel * 100)}%
                {connectedReader.isCharging ? ' (Charging)' : ''}
              </Text>
            )}
            <TouchableOpacity 
              style={styles.disconnectButton} 
              onPress={disconnectReader}
            >
              <Text style={styles.disconnectButtonText}>Disconnect</Text>
            </TouchableOpacity>
          </View>
        </View>
      ) : (
        <View style={styles.discoverySection}>
          <Text style={styles.sectionTitle}>Available Readers</Text>
          
          <TouchableOpacity
            style={[styles.discoverButton, isDiscovering && styles.disabledButton]}
            onPress={discoverReaders}
            disabled={isDiscovering}
          >
            {isDiscovering ? (
              <ActivityIndicator size="small" color="#fff" />
            ) : (
              <Text style={styles.discoverButtonText}>Scan for Readers</Text>
            )}
          </TouchableOpacity>

          {readers.length > 0 && (
            <FlatList
              data={readers}
              keyExtractor={(item) => item.serialNumber}
              renderItem={renderReader}
              style={styles.readersList}
            />
          )}

          {isConnecting && (
            <View style={styles.connectingOverlay}>
              <ActivityIndicator size="large" color="#2e7d32" />
              <Text style={styles.connectingText}>Connecting to reader...</Text>
            </View>
          )}
        </View>
      )}

      <View style={styles.instructions}>
        <Text style={styles.instructionsTitle}>Instructions:</Text>
        <Text style={styles.instructionsText}>
          1. Make sure your card reader is turned on{'\n'}
          2. Put the reader in pairing mode{'\n'}
          3. Tap "Scan for Readers" to discover devices{'\n'}
          4. Select your reader from the list to connect
        </Text>
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
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  closeButton: {
    padding: 10,
  },
  closeButtonText: {
    fontSize: 24,
    color: '#757575',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 20,
    fontSize: 16,
    color: '#757575',
  },
  connectedSection: {
    padding: 20,
  },
  discoverySection: {
    flex: 1,
    padding: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  connectedReader: {
    backgroundColor: '#e8f5e9',
    padding: 15,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#2e7d32',
  },
  connectedReaderName: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2e7d32',
  },
  connectedReaderSerial: {
    fontSize: 14,
    color: '#1b5e20',
    marginTop: 5,
  },
  connectedReaderBattery: {
    fontSize: 14,
    color: '#1b5e20',
    marginTop: 5,
  },
  disconnectButton: {
    backgroundColor: '#d32f2f',
    padding: 10,
    borderRadius: 5,
    marginTop: 10,
    alignItems: 'center',
  },
  disconnectButtonText: {
    color: '#fff',
    fontWeight: 'bold',
  },
  discoverButton: {
    backgroundColor: '#2e7d32',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
    marginBottom: 20,
  },
  disabledButton: {
    backgroundColor: '#a5d6a7',
  },
  discoverButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  readersList: {
    flex: 1,
  },
  readerItem: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    padding: 15,
    marginBottom: 10,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
  },
  readerInfo: {
    flex: 1,
  },
  readerName: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  readerSerial: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  readerType: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  readerBattery: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  connectButton: {
    backgroundColor: '#2e7d32',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 5,
    justifyContent: 'center',
  },
  connectButtonText: {
    color: '#fff',
    fontWeight: 'bold',
  },
  connectingOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  connectingText: {
    color: '#fff',
    fontSize: 16,
    marginTop: 10,
  },
  instructions: {
    backgroundColor: '#fff',
    padding: 20,
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
  },
  instructionsTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  instructionsText: {
    fontSize: 14,
    color: '#757575',
    lineHeight: 20,
  },
});