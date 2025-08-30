import React, { useState } from 'react'
import { Modal } from './Modal'
import { createOrderDraft } from '../../services/dashboard'

interface Props { open:boolean; onClose:()=>void }
export const CreateOrderModal:React.FC<Props> = ({ open, onClose }) => {
  const [customer,setCustomer] = useState('')
  const [total,setTotal] = useState('')
  const [saving,setSaving] = useState(false)
  const [message,setMessage] = useState<string|undefined>()

  async function submit(e:React.FormEvent){
    e.preventDefault()
    setSaving(true)
    const order = await createOrderDraft({ customer, total: parseFloat(total)||0 })
    setMessage(`Draft order created (id: ${order.id}). Replace with real Create Order API call.`)
    setSaving(false)
  }

  return (
    <Modal open={open} onClose={onClose} title="Create Order" footer={<div className="flex justify-end gap-3">
      <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Close</button>
      <button form="create-order-form" disabled={saving} className="px-4 py-2 text-sm rounded-md bg-green-600 text-white hover:bg-green-700 disabled:opacity-50">{saving?'Creating...':'Create Order'}</button>
    </div>}>
      <form id="create-order-form" onSubmit={submit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Customer Name</label>
          <input value={customer} onChange={e=>setCustomer(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Order Total (USD)</label>
          <input type="number" step="0.01" value={total} onChange={e=>setTotal(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-green-500 focus:border-green-500 text-sm" />
        </div>
        {message && <p className="text-xs text-green-600 bg-green-50 rounded-md p-2">{message}</p>}
        <p className="text-xs text-gray-500">NOTE: This is a front-end simulation. Integrate with Orders service later.</p>
      </form>
    </Modal>
  )
}
