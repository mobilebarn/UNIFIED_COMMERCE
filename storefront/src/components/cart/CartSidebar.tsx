'use client';

import { Fragment } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { XMarkIcon, MinusIcon, PlusIcon } from '@heroicons/react/24/outline';
import { useCartStore } from '@/stores/cart';
import Link from 'next/link';

export function CartSidebar() {
  const { isOpen, closeCart, items, updateQuantity, removeItem, totalPrice } = useCartStore();

  return (
    <Transition.Root show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={closeCart}>
        <Transition.Child
          as={Fragment}
          enter="ease-in-out duration-500"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in-out duration-500"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10">
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-500 sm:duration-700"
                enterFrom="translate-x-full"
                enterTo="translate-x-0"
                leave="transform transition ease-in-out duration-500 sm:duration-700"
                leaveFrom="translate-x-0"
                leaveTo="translate-x-full"
              >
                <Dialog.Panel className="pointer-events-auto w-screen max-w-md">
                  <div className="flex h-full flex-col bg-white shadow-xl">
                    {/* Header */}
                    <div className="flex items-start justify-between p-4 border-b">
                      <Dialog.Title className="text-lg font-medium text-gray-900">
                        Shopping Cart
                      </Dialog.Title>
                      <div className="ml-3 flex h-7 items-center">
                        <button
                          type="button"
                          className="relative -m-2 p-2 text-gray-400 hover:text-gray-500"
                          onClick={closeCart}
                        >
                          <XMarkIcon className="h-6 w-6" />
                        </button>
                      </div>
                    </div>

                    {/* Cart Items */}
                    <div className="flex-1 overflow-y-auto p-4">
                      {items.length === 0 ? (
                        <div className="text-center py-8">
                          <p className="text-gray-500">Your cart is empty</p>
                          <Link
                            href="/products"
                            className="mt-4 inline-block text-blue-600 hover:text-blue-500"
                            onClick={closeCart}
                          >
                            Continue Shopping
                          </Link>
                        </div>
                      ) : (
                        <div className="space-y-4">
                          {items.map((item) => (
                            <div key={item.id} className="flex items-center space-x-4 bg-gray-50 p-4 rounded-lg">
                              {/* Product Image */}
                              <div className="h-16 w-16 flex-shrink-0 bg-gray-200 rounded-md">
                                {item.image && (
                                  <img
                                    src={item.image}
                                    alt={item.name}
                                    className="h-full w-full object-cover rounded-md"
                                  />
                                )}
                              </div>

                              {/* Product Details */}
                              <div className="flex-1 min-w-0">
                                <h4 className="text-sm font-medium text-gray-900 truncate">
                                  {item.name}
                                </h4>
                                <p className="text-sm text-gray-500">
                                  ${item.price.toFixed(2)}
                                </p>
                                
                                {/* Quantity Controls */}
                                <div className="flex items-center mt-2 space-x-2">
                                  <button
                                    onClick={() => updateQuantity(item.id, item.quantity - 1)}
                                    className="p-1 text-gray-400 hover:text-gray-600"
                                  >
                                    <MinusIcon className="h-4 w-4" />
                                  </button>
                                  <span className="text-sm font-medium">{item.quantity}</span>
                                  <button
                                    onClick={() => updateQuantity(item.id, item.quantity + 1)}
                                    className="p-1 text-gray-400 hover:text-gray-600"
                                  >
                                    <PlusIcon className="h-4 w-4" />
                                  </button>
                                </div>
                              </div>

                              {/* Remove Button */}
                              <button
                                onClick={() => removeItem(item.id)}
                                className="text-red-400 hover:text-red-600"
                              >
                                <XMarkIcon className="h-5 w-5" />
                              </button>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Footer */}
                    {items.length > 0 && (
                      <div className="border-t p-4">
                        <div className="flex items-center justify-between mb-4">
                          <span className="text-lg font-medium text-gray-900">Total:</span>
                          <span className="text-lg font-bold text-gray-900">
                            ${totalPrice().toFixed(2)}
                          </span>
                        </div>
                        <Link
                          href="/checkout"
                          className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors text-center block"
                          onClick={closeCart}
                        >
                          Proceed to Checkout
                        </Link>
                      </div>
                    )}
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  );
}
