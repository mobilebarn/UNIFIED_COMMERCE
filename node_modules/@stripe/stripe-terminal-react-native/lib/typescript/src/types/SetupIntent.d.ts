import type { PaymentMethodType } from './';
export declare namespace SetupIntent {
    type Type = Android.Type & IOS.Type & {
        id: string;
        created: string;
        customer?: string;
        metadata?: Record<string, string>;
        status: Status;
        latestAttempt: SetupAttempt;
        usage: Usage;
        sdkUuid: string;
        paymentMethodTypes?: string[];
    };
    type Status = 'canceled' | 'processing' | 'requiresAction' | 'requiresConfirmation' | 'requiresPaymentMethod' | 'succeeded' | 'unknown';
    type Usage = 'offSession' | 'onSession';
    type SetupAttempt = IOS.SetupAttempt & Android.SetupAttempt & {
        id: string;
        created?: number;
        status: string;
        customer?: string;
        setupIntentId?: string;
        paymentMethodDetails: SetupAttemptPaymentMethodDetails;
        onBehalfOfId?: string;
        applicationId?: string;
        paymentMethodId?: string;
    };
    interface SetupAttemptPaymentMethodDetails {
        cardPresent: SetupAttemptCardPresentDetails;
        interacPresent: SetupAttemptCardPresentDetails;
        type: PaymentMethodType;
    }
    interface SetupAttemptCardPresentDetails {
        emvAuthData: string;
        generatedCard: string;
    }
    namespace IOS {
        type SetupAttempt = {};
        type Type = {};
    }
    namespace Android {
        type SetupAttempt = {
            isLiveMode: boolean;
            usage?: Usage;
        };
        type Type = {
            applicationId?: string;
            clientSecret?: string;
            description?: string;
            mandateId?: string;
            onBehalfOfId?: string;
            paymentMethodId?: string;
            singleUseMandateId?: string;
        };
    }
}
