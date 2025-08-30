import React from 'react'

export default function Analytics(){
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Analytics (Placeholder)</h1>
      <p className="text-gray-600 mb-6 max-w-prose">This page will surface deeper insights (funnels, cohort retention, LTV, channel performance). Build after core order & product CRUD is fully wired.</p>
      <ul className="list-disc pl-6 space-y-2 text-sm text-gray-700">
        <li>Planned: Sales over time with comparison periods</li>
        <li>Planned: Product performance breakdown with filters</li>
        <li>Planned: Customer retention / repeat purchase rate visualization</li>
        <li>Planned: Promotion effectiveness metrics</li>
      </ul>
    </div>
  )
}
