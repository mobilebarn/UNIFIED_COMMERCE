import type { PaymentMethodDetails } from './';
export declare namespace Refund {
    type Props = IOS.Props & Android.Props & {
        id: string;
        amount?: number;
        chargeId: string;
        currency?: string;
        failureReason?: string;
        metadata?: Record<string, string>;
        reason?: string;
        description?: string;
        status?: Status;
    };
    namespace IOS {
        interface Props {
            created: string;
            paymentMethodDetails?: PaymentMethodDetails;
        }
    }
    type Status = 'succeeded' | 'failed' | 'pending' | 'unknown';
    namespace Android {
        interface Props {
            balanceTransaction?: string;
            failureBalanceTransaction?: string;
            receiptNumber?: string;
            sourceTransferReversal?: string;
            transferReversal?: string;
            paymentIntentId?: string;
        }
    }
}
