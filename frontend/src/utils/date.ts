// src/utils/date.ts

/**
 * Convert date string to RFC3339 format for backend
 * Input: "2026-09-01"
 * Output: "2026-09-01T00:00:00Z"
 * 
 * Input: null or undefined
 * Output: undefined
 */
export const formatDateForBackend = (dateStr: string | null | undefined): string | undefined => {
    if (!dateStr) return undefined
    return `${dateStr}T00:00:00Z`
}

/**
 * Format date for display in date input fields
 * Input: "2026-09-01T00:00:00Z"
 * Output: "2026-09-01"
 */
export const formatDateForDisplay = (dateStr: string | null | undefined): string => {
    if (!dateStr) return ''
    try {
        return new Date(dateStr).toISOString().slice(0, 10)
    } catch {
        return ''
    }
}