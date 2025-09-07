'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { GET_PAYMENT_METHODS, ADD_PAYMENT_METHOD, REMOVE_PAYMENT_METHOD, SET_DEFAULT_PAYMENT_METHOD } from '@/graphql/queries';
import type { PaymentMethod } from '@/graphql/queries';

export default function PaymentMethods() {
  const [showAddForm, setShowAddForm] = useState(false);
  const { data, loading, error, refetch } = useQuery(GET_PAYMENT_METHODS);
  const [addPaymentMethod] = useMutation(ADD_PAYMENT_METHOD, {
    onCompleted: () => {
      refetch();
      setShowAddForm(false);
    }
  });
  const [removePaymentMethod] = useMutation(REMOVE_PAYMENT_METHOD, {
    onCompleted: () => {
      refetch();
    }
  });
  const [setDefaultPaymentMethod] = useMutation(SET_DEFAULT_PAYMENT_METHOD, {
    onCompleted: () => {
      refetch();
    }
  });

  const paymentMethods = data?.paymentMethods || [];

  const handleAddPaymentMethod = async (formData: FormData) => {
    try {
      const type = formData.get('type') as 'credit_card' | 'debit_card' | 'paypal' | 'bank_account';
      
      let input: any = {
        type,
        name: formData.get('name') as string,
        isDefault: formData.get('isDefault') === 'on'
      };

      if (type === 'credit_card' || type === 'debit_card') {
        input = {
          ...input,
          number: formData.get('number') as string,
          expiryMonth: parseInt(formData.get('expiryMonth') as string),
          expiryYear: parseInt(formData.get('expiryYear') as string),
          cvv: formData.get('cvv') as string
        };
      } else if (type === 'paypal') {
        input = {
          ...input,
          email: formData.get('email') as string
        };
      } else if (type === 'bank_account') {
        input = {
          ...input,
          accountNumber: formData.get('accountNumber') as string,
          routingNumber: formData.get('routingNumber') as string
        };
      }

      await addPaymentMethod({ variables: { input } });
    } catch (err) {
      console.error('Error adding payment method:', err);
    }
  };

  const handleRemovePaymentMethod = async (id: string) => {
    if (window.confirm('Are you sure you want to remove this payment method?')) {
      try {
        await removePaymentMethod({ variables: { id } });
      } catch (err) {
        console.error('Error removing payment method:', err);
      }
    }
  };

  const handleSetDefault = async (id: string) => {
    try {
      await setDefaultPaymentMethod({ variables: { id } });
    } catch (err) {
      console.error('Error setting default payment method:', err);
    }
  };

  const getPaymentMethodIcon = (type: string) => {
    switch (type) {
      case 'credit_card':
      case 'debit_card':
        return (
          <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
          </svg>
        );
      case 'paypal':
        return (
          <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        );
      case 'bank_account':
        return (
          <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
          </svg>
        );
      default:
        return (
          <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        );
    }
  };

  const formatExpiryDate = (month?: number, year?: number) => {
    if (!month || !year) return '';
    return `${month.toString().padStart(2, '0')}/${year.toString().slice(-2)}`;
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-700">Error loading payment methods: {error.message}</p>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Payment Methods</h1>
        <button
          onClick={() => setShowAddForm(true)}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          Add Payment Method
        </button>
      </div>

      {showAddForm && (
        <div className="bg-white shadow rounded-lg mb-6">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Add New Payment Method</h3>
            <form action={handleAddPaymentMethod}>
              <div className="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
                <div className="sm:col-span-6">
                  <label htmlFor="type" className="block text-sm font-medium text-gray-700">
                    Payment Method Type
                  </label>
                  <div className="mt-1">
                    <select
                      id="type"
                      name="type"
                      required
                      className="py-2 px-3 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 rounded-md"
                    >
                      <option value="">Select a payment method</option>
                      <option value="credit_card">Credit Card</option>
                      <option value="debit_card">Debit Card</option>
                      <option value="paypal">PayPal</option>
                      <option value="bank_account">Bank Account</option>
                    </select>
                  </div>
                </div>

                <div className="sm:col-span-6">
                  <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                    Name
                  </label>
                  <div className="mt-1">
                    <input
                      type="text"
                      name="name"
                      id="name"
                      required
                      className="py-2 px-3 block w-full shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300 rounded-md"
                    />
                  </div>
                </div>

                <div className="sm:col-span-6">
                  <div className="flex items-center">
                    <input
                      id="isDefault"
                      name="isDefault"
                      type="checkbox"
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label htmlFor="isDefault" className="ml-2 block text-sm text-gray-900">
                      Set as default payment method
                    </label>
                  </div>
                </div>
              </div>

              <div className="mt-6 flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setShowAddForm(false)}
                  className="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  Save Payment Method
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {paymentMethods.length === 0 ? (
        <div className="text-center py-12">
          <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No payment methods</h3>
          <p className="mt-1 text-sm text-gray-500">Get started by adding a new payment method.</p>
          <div className="mt-6">
            <button
              onClick={() => setShowAddForm(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Add Payment Method
            </button>
          </div>
        </div>
      ) : (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {paymentMethods.map((paymentMethod: PaymentMethod) => (
              <li key={paymentMethod.id}>
                <div className="px-4 py-5 sm:px-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div className="flex-shrink-0">
                        {getPaymentMethodIcon(paymentMethod.type)}
                      </div>
                      <div className="ml-4">
                        <p className="text-sm font-medium text-gray-900">
                          {paymentMethod.name}
                          {paymentMethod.isDefault && (
                            <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                              Default
                            </span>
                          )}
                        </p>
                        <p className="text-sm text-gray-500">
                          {paymentMethod.type === 'credit_card' || paymentMethod.type === 'debit_card' ? (
                            <>
                              **** **** **** {paymentMethod.last4}
                              {paymentMethod.expiryMonth && paymentMethod.expiryYear && (
                                <span className="ml-2">
                                  Expires {formatExpiryDate(paymentMethod.expiryMonth, paymentMethod.expiryYear)}
                                </span>
                              )}
                            </>
                          ) : paymentMethod.type === 'paypal' ? (
                            'PayPal Account'
                          ) : (
                            `Bank Account ending in ${paymentMethod.last4}`
                          )}
                        </p>
                      </div>
                    </div>
                    <div className="flex space-x-2">
                      {!paymentMethod.isDefault && (
                        <button
                          onClick={() => handleSetDefault(paymentMethod.id)}
                          className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-blue-700 bg-blue-100 hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        >
                          Make Default
                        </button>
                      )}
                      <button
                        onClick={() => handleRemovePaymentMethod(paymentMethod.id)}
                        className="inline-flex items-center px-3 py-1 border border-gray-300 text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}