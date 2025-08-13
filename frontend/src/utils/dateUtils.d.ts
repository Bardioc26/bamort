/**
 * TypeScript Definitionen für Utility-Funktionen
 */

export interface DateFormatOptions {
  year?: 'numeric' | '2-digit'
  month?: 'numeric' | '2-digit' | 'long' | 'short' | 'narrow'
  day?: 'numeric' | '2-digit'
  hour?: 'numeric' | '2-digit'
  minute?: 'numeric' | '2-digit'
  second?: 'numeric' | '2-digit'
}

export declare function formatDate(
  dateString: string | null | undefined,
  locale?: string,
  options?: DateFormatOptions
): string

export declare function formatDateTime(
  dateString: string | null | undefined,
  locale?: string
): string

export declare function formatRelativeDate(
  dateString: string | null | undefined,
  locale?: string
): string

export declare function safeValue<T>(
  value: T | null | undefined,
  fallback?: T
): T

export declare function capitalize(str: string): string
