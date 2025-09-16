import type { AmountDetails, Charge, OfflineDetails, PaymentMethod } from './';
export declare namespace PaymentIntent {
    interface Type {
        id: string;
        amount: number;
        amountDetails: AmountDetails;
        amountTip: number;
        captureMethod: string;
        charges: Charge[];
        clientSecret: string;
        created: string;
        currency: string;
        statementDescriptor: string;
        statementDescriptorSuffix: string;
        status: Status;
        sdkUuid: string;
        paymentMethodId: string;
        paymentMethod: PaymentMethod.Type;
        offlineDetails: OfflineDetails;
        metadata: Record<string, string>;
    }
    type Status = 'canceled' | 'processing' | 'requiresAction' | 'requiresCapture' | 'requiresConfirmation' | 'requiresPaymentMethod' | 'succeeded' | 'unknown';
}
