import { type ClassValue, clsx } from 'clsx'
import { format } from 'date-fns-tz'
import { twMerge } from 'tailwind-merge'
import { camelize, getCurrentInstance, toHandlerKey } from 'vue'

/**
 * Combines multiple class strings or objects into a single class string.
 * @param {...ClassValue} inputs - Class strings or objects to merge.
 * @returns {string} - The merged class string.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Formats a date into a string using the 'yyyy-MM-dd hh:mm:ss' format and the 'Australia/Melbourne' timezone.
 * @param {Date} date - The date to format.
 * @returns {string} - The formatted date string.
 */
export const formatDate = (date: Date): string => {
  return format(date, 'yyyy-MM-dd hh:mm:ss', {
    timeZone: 'Australia/Melbourne'
  })
}

/**
 * Parses a date string into a human-readable format.
 * @param {string} dateStr - The date string to parse.
 * @returns {string} - The parsed date string.
 */
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
