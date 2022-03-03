export async function downloadFile(name: string, data: string) {
  const fileName = name
  // const json = JSON.stringify(data)
  const blob = new Blob([data], { type: 'application/json' })
  const href = await URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = href
  link.download = fileName + '.txt'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
