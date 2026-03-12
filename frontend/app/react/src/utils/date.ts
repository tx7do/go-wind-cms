export function formatDate(dateString: string | undefined) {
    if (!dateString) return ''
    return new Date(dateString).toLocaleString()
}
