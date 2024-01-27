import { type ClassValue, clsx } from 'clsx'
import { format } from 'date-fns-tz'
import { twMerge } from 'tailwind-merge'
import { camelize, getCurrentInstance, toHandlerKey } from 'vue'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


export const formatDate = (date: Date): string => {
  return format(date, 'yyyy-MM-dd hh:mm:ss', {
    timeZone: 'Australia/Melbourne'
  })
}

export const parseDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat("en-US", {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }).format(date);
}