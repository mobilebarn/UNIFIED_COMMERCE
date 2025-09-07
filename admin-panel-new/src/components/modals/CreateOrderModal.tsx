import React, { useState } from 'react'
import { Modal } from './Modal'
import { useCreateOrder } from '../../hooks/useGraphQL'
import { useRefreshQueries } from '../../hooks/useGraphQL'

interface Props { open:boolean; onClose:()=>void }
export const CreateOrderModal:React.FC<Props> = ({ open, onClose }) => {
  const [customerEmail,setCustomerEmail] = useState('')
  const [customerFirstName,setCustomerFirstName] = useState('')
  const [customerLastName,setCustomerLastName] = useState('')
  const [total,setTotal] = useState('')
  const [currency,setCurrency] = useState('USD')
  const [saving,setSaving] = useState(false)
  const [message,setMessage] = useState<string|undefined>()
  const [error,setError] = useState<string|undefined>()
  
  const [createOrder] = useCreateOrder()
  const { refreshOrders } = useRefreshQueries()

  async function submit(e:React.FormEvent){
    e.preventDefault()
    setSaving(true)
    setError(undefined)
    setMessage(undefined)
    
    try {
      // In a real implementation, we would look up the customer by email
      // For now, we'll create a basic order with the provided information
      await createOrder({
        variables: {
          input: {
            merchantId: 'merchant-1', // This would be dynamically determined in a real app
            customer: {
              email: customerEmail,
              firstName: customerFirstName,
              lastName: customerLastName,
              phone: ''
            },
            billingAddress: {
              firstName: customerFirstName,
              lastName: customerLastName,
              street1: '',
              city: '',
              state: '',
              country: '',
              postalCode: ''
            },
            shippingAddress: {
              firstName: customerFirstName,
              lastName: customerLastName,
              street1: '',
              city: '',
              state: '',
              country: '',
              postalCode: ''
            },
            currency: currency,
            source: 'ONLINE'
          }
        }
      })
      
      setMessage('Order created successfully!')
      setCustomerEmail('')
      setCustomerFirstName('')
      setCustomerLastName('')
      setTotal('')
      
      // Refresh the orders list
      await refreshOrders()
    } catch (err: any) {
      setError(err.message || 'Failed to create order')
      console.error('Error creating order:', err)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Create Order" footer={<div className="flex justify-end gap-3">
      <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Close</button>
      <button form="create-order-form" disabled={saving} className="px-4 py-2 text-sm rounded-md bg-green-600 text-white hover:bg-green-700 disabled:opacity-50">{saving?'Creating...':'Create Order'}</button>
    </div>}>
      <form id="create-order-form" onSubmit={submit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Customer First Name</label>
            <input 
              value={customerFirstName} 
              onChange={e=>setCustomerFirstName(e.target.value)} 
              required 
              className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Customer Last Name</label>
            <input 
              value={customerLastName} 
              onChange={e=>setCustomerLastName(e.target.value)} 
              required 
              className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" 
            />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Customer Email</label>
          <input 
            type="email" 
            value={customerEmail} 
            onChange={e=>setCustomerEmail(e.target.value)} 
            required 
            className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" 
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Order Total ({currency})</label>
            <input 
              type="number" 
              step="0.01" 
              value={total} 
              onChange={e=>setTotal(e.target.value)} 
              required 
              className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
            <select 
              value={currency} 
              onChange={e=>setCurrency(e.target.value)} 
              className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm"
            >
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="GBP">GBP</option>
            </select>
          </div>
        </div>
        {message && <p className="text-xs text-green-600 bg-green-50 rounded-md p-2">{message}</p>}
        {error && <p className="text-xs text-red-600 bg-red-50 rounded-md p-2">{error}</p>}
      </form>
    </Modal>
  )
}