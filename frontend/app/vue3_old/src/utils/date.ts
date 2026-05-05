export function formatDate(dateString: string) {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString()
}
