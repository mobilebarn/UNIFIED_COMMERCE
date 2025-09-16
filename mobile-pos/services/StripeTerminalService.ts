import { 
  initialize, 
  discoverReaders, 
  connectBluetoothReader,
  connectInternetReader,
  createPaymentIntent,
  collectPaymentMethod,
  processPayment,
  Reader,
  Terminal,
  PaymentIntent,
  DiscoveryConfiguration,
  ConnectionConfiguration,
  PaymentIntentParameters,
  CollectConfiguration
} from '@stripe/stripe-terminal-react-native';

export interface PaymentResult {
  success: boolean;
  paymentIntent?: PaymentIntent;
  error?: string;
}

export interface ReaderInfo {
  id: string;
  label?: string;
  serialNumber: string;
  deviceType: string;
  batteryLevel?: number;
  isCharging?: boolean;
}

class StripeTerminalService {
  private isInitialized = false;
  private connectedReader: Reader | null = null;

  /**
   * Initialize Stripe Terminal SDK
   */
  async initialize(): Promise<boolean> {
    try {
      if (this.isInitialized) {
        return true;
      }

      await initialize({
        logLevel: 'verbose', // Use 'none' in production
      });

      this.isInitialized = true;
      console.log('Stripe Terminal initialized successfully');
      return true;
    } catch (error) {
      console.error('Failed to initialize Stripe Terminal:', error);
      return false;
    }
  }

  /**
   * Discover available card readers
   */
  async discoverReaders(type: 'bluetooth' | 'internet' = 'bluetooth'): Promise<Reader[]> {
    try {
      if (!this.isInitialized) {
        throw new Error('Stripe Terminal not initialized');
      }

      const config: DiscoveryConfiguration = {
        discoveryMethod: type === 'bluetooth' ? 'bluetoothScan' : 'internet',
        simulated: __DEV__, // Use simulated readers in development
      };

      console.log('Discovering readers...');
      const { readers } = await discoverReaders(config);
      console.log(`Found ${readers.length} readers`);
      
      return readers;
    } catch (error) {
      console.error('Failed to discover readers:', error);
      return [];
    }
  }

  /**
   * Connect to a card reader
   */
  async connectReader(reader: Reader, locationId?: string): Promise<boolean> {
    try {
      if (!this.isInitialized) {
        throw new Error('Stripe Terminal not initialized');
      }

      const config: ConnectionConfiguration = {
        locationId: locationId || 'tml_FakeLocationForDevelopment', // Use real location ID in production
      };

      let connectedReader: Reader;

      if (reader.deviceType === 'bluetoothLE') {
        connectedReader = await connectBluetoothReader(reader, config);
      } else {
        connectedReader = await connectInternetReader(reader, config);
      }

      this.connectedReader = connectedReader;
      console.log('Successfully connected to reader:', connectedReader.serialNumber);
      return true;
    } catch (error) {
      console.error('Failed to connect to reader:', error);
      return false;
    }
  }

  /**
   * Process a payment
   */
  async processPayment(amount: number, currency: string = 'usd'): Promise<PaymentResult> {
    try {
      if (!this.isInitialized) {
        throw new Error('Stripe Terminal not initialized');
      }

      if (!this.connectedReader) {
        throw new Error('No reader connected');
      }

      // Convert amount to cents
      const amountInCents = Math.round(amount * 100);

      // Create payment intent
      const paymentIntentParams: PaymentIntentParameters = {
        amount: amountInCents,
        currency: currency.toLowerCase(),
        paymentMethodTypes: ['card_present'],
        captureMethod: 'automatic',
      };

      console.log('Creating payment intent for amount:', amountInCents);
      const { paymentIntent } = await createPaymentIntent(paymentIntentParams);

      // Collect payment method
      const collectConfig: CollectConfiguration = {
        skipTipping: true, // Can be configured based on business needs
      };

      console.log('Collecting payment method...');
      const { paymentIntent: collectedPaymentIntent } = await collectPaymentMethod(
        paymentIntent,
        collectConfig
      );

      // Process payment
      console.log('Processing payment...');
      const { paymentIntent: processedPaymentIntent } = await processPayment(collectedPaymentIntent);

      if (processedPaymentIntent.status === 'succeeded') {
        console.log('Payment succeeded!');
        return {
          success: true,
          paymentIntent: processedPaymentIntent,
        };
      } else {
        console.log('Payment failed with status:', processedPaymentIntent.status);
        return {
          success: false,
          error: `Payment failed with status: ${processedPaymentIntent.status}`,
        };
      }
    } catch (error) {
      console.error('Payment processing failed:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred',
      };
    }
  }

  /**
   * Get connected reader info
   */
  getConnectedReader(): ReaderInfo | null {
    if (!this.connectedReader) {
      return null;
    }

    return {
      id: this.connectedReader.id || 'unknown',
      label: this.connectedReader.label,
      serialNumber: this.connectedReader.serialNumber,
      deviceType: this.connectedReader.deviceType,
      batteryLevel: this.connectedReader.batteryLevel,
      isCharging: this.connectedReader.isCharging,
    };
  }

  /**
   * Disconnect from current reader
   */
  async disconnectReader(): Promise<boolean> {
    try {
      if (this.connectedReader) {
        // The SDK handles disconnection automatically
        this.connectedReader = null;
        console.log('Disconnected from reader');
      }
      return true;
    } catch (error) {
      console.error('Failed to disconnect reader:', error);
      return false;
    }
  }

  /**
   * Check if a reader is connected
   */
  isReaderConnected(): boolean {
    return this.connectedReader !== null;
  }
}

// Export singleton instance
export const stripeTerminalService = new StripeTerminalService();
export default stripeTerminalService;