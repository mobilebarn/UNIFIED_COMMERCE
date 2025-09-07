import React, { useState } from 'react'
import { Modal } from './Modal'
import { useRegisterUser } from '../../hooks/useGraphQL'
import { useRefreshQueries } from '../../hooks/useGraphQL'

interface Props { open:boolean; onClose:()=>void }
export const AddCustomerModal:React.FC<Props> = ({ open, onClose }) => {
  const [firstName,setFirstName] = useState('')
  const [lastName,setLastName] = useState('')
  const [email,setEmail] = useState('')
  const [phone,setPhone] = useState('')
  const [saving,setSaving] = useState(false)
  const [message,setMessage] = useState<string|undefined>()
  const [error,setError] = useState<string|undefined>()
  
  const [registerUser] = useRegisterUser()
  const { refreshCustomers } = useRefreshQueries()

  async function submit(e:React.FormEvent){
    e.preventDefault()
    setSaving(true)
    setError(undefined)
    setMessage(undefined)
    
    try {
      // Generate a random username and password for the customer
      const username = email.split('@')[0] + Math.floor(Math.random() * 1000)
      const password = 'TempPass123!' // In a real app, this would be generated securely
      
      await registerUser({
        variables: {
          input: {
            firstName,
            lastName,
            email,
            phone,
            username,
            password
          }
        }
      })
      
      setMessage('Customer created successfully!')
      setFirstName('')
      setLastName('')
      setEmail('')
      setPhone('')
      
      // Refresh the customers list
      await refreshCustomers()
    } catch (err: any) {
      setError(err.message || 'Failed to create customer')
      console.error('Error creating customer:', err)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Customer" footer={<div className="flex justify-end gap-3">
      <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Close</button>
      <button form="add-customer-form" disabled={saving} className="px-4 py-2 text-sm rounded-md bg-purple-600 text-white hover:bg-purple-700 disabled:opacity-50">{saving?'Saving...':'Save Customer'}</button>
    </div>}>
      <form id="add-customer-form" onSubmit={submit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">First Name</label>
            <input 
              value={firstName} 
              onChange={e=>setFirstName(e.target.value)} 
              required 
              className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
            <input 
              value={lastName} 
              onChange={e=>setLastName(e.target.value)} 
              required 
              className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" 
            />
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input 
            type="email" 
            value={email} 
            onChange={e=>setEmail(e.target.value)} 
            required 
            className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" 
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
          <input 
            type="tel" 
            value={phone} 
            onChange={e=>setPhone(e.target.value)} 
            className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" 
          />
        </div>
        {message && <p className="text-xs text-green-600 bg-green-50 rounded-md p-2">{message}</p>}
        {error && <p className="text-xs text-red-600 bg-red-50 rounded-md p-2">{error}</p>}
      </form>
    </Modal>
  )
}