// src/utils/image.ts
import api from '@/utils/axios'
import { push } from 'notivue'

interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data: T
}

interface ImageData {
  image_url: string // Changed from "url" to "image_url"
}

export const uploadImage = async (
  entity: string,
  id: number,
  file: File,
  onProgress?: (progress: number) => void,
): Promise<string | null> => {
  try {
    const formData = new FormData()
    formData.append('image', file)

    const res = await api.post<ApiResponse<ImageData>>(`/images/${entity}/${id}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total && onProgress) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      },
    })

    // Check response - use "image_url" instead of "url"
    if (!res.data?.success || !res.data?.data?.image_url) {
      push.error(res.data?.message || 'Upload failed')
      return null
    }

    // Build full URL
    const baseURL = (api.defaults.baseURL || 'http://localhost:8080').replace(/\/+$/, '')
    const imagePath = res.data.data.image_url.startsWith('/')
      ? res.data.data.image_url
      : `/${res.data.data.image_url}`

    return `${baseURL}${imagePath}`
  } catch (error) {
    const err = error as { response?: { data?: { message?: string } } }
    push.error(err.response?.data?.message || 'Failed to upload image')
    return null
  }
}

export const deleteImage = async (
  entity: string,
  id: number,
  imageUrl: string,
): Promise<boolean> => {
  try {
    let relativePath = imageUrl
    try {
      const url = new URL(imageUrl)
      relativePath = url.pathname
    } catch {
      // Not a full URL, use as-is
    }

    const res = await api.delete<ApiResponse>(`/images/${entity}/${id}`, {
      data: { image_url: relativePath },
    })

    if (!res.data?.success) {
      push.error(res.data?.message || 'Delete failed')
      return false
    }

    return true
  } catch (error) {
    const err = error as { response?: { data?: { message?: string } } }
    push.error(err.response?.data?.message || 'Failed to delete image')
    return false
  }
}

export const getImageUrl = (path: string | null | undefined): string | undefined => {
  if (!path) return undefined

  // If it's already a full URL, return it
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }

  // Get the base URL and remove /api from it
  const apiBaseURL = api.defaults.baseURL || 'http://localhost:8080/api'
  const baseURL = apiBaseURL.replace(/\/api\/?$/, '') // Remove /api from the end

  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${baseURL}${cleanPath}`
}
