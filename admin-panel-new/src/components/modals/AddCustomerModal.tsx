import React, { useState } from 'react'
import { Modal } from './Modal'
import { createCustomerDraft } from '../../services/dashboard'

interface Props { open:boolean; onClose:()=>void }
export const AddCustomerModal:React.FC<Props> = ({ open, onClose }) => {
  const [name,setName] = useState('')
  const [email,setEmail] = useState('')
  const [saving,setSaving] = useState(false)
  const [message,setMessage] = useState<string|undefined>()

  async function submit(e:React.FormEvent){
    e.preventDefault()
    setSaving(true)
    const cust = await createCustomerDraft({ name, email })
    setMessage(`Draft customer created (id: ${cust.id}). Replace with real Customer API call.`)
    setSaving(false)
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Customer" footer={<div className="flex justify-end gap-3">
      <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Close</button>
      <button form="add-customer-form" disabled={saving} className="px-4 py-2 text-sm rounded-md bg-purple-600 text-white hover:bg-purple-700 disabled:opacity-50">{saving?'Saving...':'Save Customer'}</button>
    </div>}>
      <form id="add-customer-form" onSubmit={submit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
          <input value={name} onChange={e=>setName(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input type="email" value={email} onChange={e=>setEmail(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-purple-500 focus:border-purple-500 text-sm" />
        </div>
        {message && <p className="text-xs text-green-600 bg-green-50 rounded-md p-2">{message}</p>}
        <p className="text-xs text-gray-500">NOTE: This is local only until backend integration.</p>
      </form>
    </Modal>
  )
}
