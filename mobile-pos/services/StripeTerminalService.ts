// Mock Stripe Terminal Service for development
// In production, this would be replaced with actual Stripe Terminal SDK integration

export interface PaymentResult {
  success: boolean;
  paymentIntent?: any;
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

export interface Reader {
  id?: string;
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

      // Mock initialization
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      this.isInitialized = true;
      console.log('Stripe Terminal initialized successfully (Mock)');
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

      // Mock discovery - return simulated readers
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      const mockReaders: Reader[] = [
        {
          id: 'mock_reader_1',
          label: 'Simulated Card Reader',
          serialNumber: 'MOCK-12345',
          deviceType: 'bluetoothLE',
          batteryLevel: 0.85,
          isCharging: false,
        },
        {
          id: 'mock_reader_2',
          label: 'Test Terminal',
          serialNumber: 'MOCK-67890',
          deviceType: 'bluetoothLE',
          batteryLevel: 0.92,
          isCharging: true,
        },
      ];

      console.log(`Found ${mockReaders.length} readers (Mock)`);
      return mockReaders;
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

      // Mock connection
      await new Promise(resolve => setTimeout(resolve, 3000));
      
      this.connectedReader = reader;
      console.log('Successfully connected to reader (Mock):', reader.serialNumber);
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

      console.log('Creating payment intent for amount (Mock):', amountInCents);
      await new Promise(resolve => setTimeout(resolve, 2000));

      console.log('Collecting payment method (Mock)...');
      await new Promise(resolve => setTimeout(resolve, 3000));

      console.log('Processing payment (Mock)...');
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Simulate successful payment
      const mockPaymentIntent = {
        id: 'pi_mock_' + Date.now(),
        amount: amountInCents,
        currency: currency.toLowerCase(),
        status: 'succeeded',
        created: Date.now() / 1000,
      };

      console.log('Payment succeeded (Mock)!');
      return {
        success: true,
        paymentIntent: mockPaymentIntent,
      };
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
        await new Promise(resolve => setTimeout(resolve, 1000));
        this.connectedReader = null;
        console.log('Disconnected from reader (Mock)');
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