import { gql } from '@apollo/client';

// Query to get payment methods for a customer
export const GET_PAYMENT_METHODS = gql`
  query GetPaymentMethods($customerId: ID) {
    paymentMethods(customerId: $customerId) {
      id
      type
      isDefault
      cardLast4
      cardBrand
      cardExpMonth
      cardExpYear
    }
  }
`;

// Mutation to create a payment
export const CREATE_PAYMENT = gql`
  mutation CreatePayment($input: CreatePaymentInput!) {
    createPayment(input: $input) {
      id
      orderId
      customerId
      merchantId
      amount
      currency
      status
      gateway
      paymentMethodId
      paymentMethodType
      createdAt
    }
  }
`;

// Mutation to capture a payment transaction
export const CAPTURE_PAYMENT = gql`
  mutation CapturePaymentTransaction($id: ID!, $amount: Float) {
    capturePaymentTransaction(id: $id, amount: $amount) {
      id
      status
      capturedAt
    }
  }
`;

// Mutation to refund a payment
export const REFUND_PAYMENT = gql`
  mutation RefundPayment($id: ID!, $input: RefundPaymentInput!) {
    refundPayment(id: $id, input: $input) {
      id
      status
      refundedAt
    }
  }
`;

// Mutation to void a payment
export const VOID_PAYMENT = gql`
  mutation VoidPayment($id: ID!) {
    voidPayment(id: $id) {
      id
      status
      voidedAt
    }
  }
`;