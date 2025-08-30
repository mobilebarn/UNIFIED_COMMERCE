import React, { useState } from 'react'
import { Modal } from './Modal'
import { createProductDraft } from '../../services/dashboard'

interface Props { open:boolean; onClose:()=>void }
export const AddProductModal:React.FC<Props> = ({ open, onClose }) => {
  const [name,setName] = useState('')
  const [price,setPrice] = useState('')
  const [saving,setSaving] = useState(false)
  const [message,setMessage] = useState<string|undefined>()

  async function submit(e:React.FormEvent){
    e.preventDefault()
    setSaving(true)
    const prod = await createProductDraft({ name, price: parseFloat(price)||0 })
    setMessage(`Draft product created (id: ${prod.id}) - replace with real API call.`)
    setSaving(false)
  }

  return (
    <Modal open={open} onClose={onClose} title="Add Product" footer={<div className="flex justify-end gap-3">
      <button onClick={onClose} className="px-4 py-2 text-sm rounded-md border border-gray-300 bg-white hover:bg-gray-50">Close</button>
      <button form="add-product-form" disabled={saving} className="px-4 py-2 text-sm rounded-md bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50">{saving?'Saving...':'Save Product'}</button>
    </div>}>
      <form id="add-product-form" onSubmit={submit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input value={name} onChange={e=>setName(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm" />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Price (USD)</label>
            <input type="number" step="0.01" value={price} onChange={e=>setPrice(e.target.value)} required className="w-full rounded-md border-gray-300 focus:ring-blue-500 focus:border-blue-500 text-sm" />
        </div>
        {message && <p className="text-xs text-green-600 bg-green-50 rounded-md p-2">{message}</p>}
        <p className="text-xs text-gray-500">NOTE: This currently creates a local draft only. Wire to backend API later.</p>
      </form>
    </Modal>
  )
}
