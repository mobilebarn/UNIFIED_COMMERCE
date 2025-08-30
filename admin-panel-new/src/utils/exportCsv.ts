export function exportToCsv(filename:string, rows:Record<string, any>[]) {
  if(!rows.length) return
  const headers = Object.keys(rows[0])
  const esc = (v:any)=> '"'+String(v).replace(/"/g,'""')+'"'
  const csv = [headers.join(','), ...rows.map(r=> headers.map(h=> esc(r[h] ?? '')).join(','))].join('\n')
  const blob = new Blob([csv], { type:'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
