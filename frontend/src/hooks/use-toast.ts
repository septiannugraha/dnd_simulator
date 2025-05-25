import { useState } from 'react'

interface Toast {
  title: string
  description?: string
  variant?: 'default' | 'destructive'
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([])

  const toast = (toast: Toast) => {
    setToasts(prev => [...prev, toast])
    
    // Simple console log for now - you can implement a proper toast UI later
    console.log(`Toast: ${toast.title}`, toast.description)
    
    // Auto remove after 5 seconds
    setTimeout(() => {
      setToasts(prev => prev.slice(1))
    }, 5000)
  }

  return { toast, toasts }
}