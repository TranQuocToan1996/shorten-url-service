const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export interface APIResponse<T = unknown> {
  code: number;
  status: string;
  message: string;
  data: T;
}

export interface ShortenURLResponse {
  id: number;
  created_at: string;
  updated_at: string;
  status: string;
  code: string;
  long_url: string;
}

export interface SubmitURLRequest {
  long_url: string;
  callback_url?: string;
}

export async function submitURL(longURL: string, callbackURL?: string): Promise<APIResponse> {
  const response = await fetch(`${API_BASE_URL}/encode`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      long_url: longURL,
      callback_url: callbackURL,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to submit URL');
  }

  return response.json();
}

export async function getByLongURL(longURL: string): Promise<APIResponse<ShortenURLResponse>> {
  const response = await fetch(`${API_BASE_URL}/urls/long?long_url=${encodeURIComponent(longURL)}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    if (response.status === 404) {
      const error = new Error('URL not found') as Error & { status?: number };
      error.status = 404;
      throw error;
    }
    const error = await response.json() as { message?: string };
    throw new Error(error.message || 'Failed to get URL');
  }

  return response.json();
}

export async function decodeURL(shortenURL: string): Promise<APIResponse<ShortenURLResponse>> {
  const response = await fetch(`${API_BASE_URL}/decode?shorten_url=${encodeURIComponent(shortenURL)}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to decode URL');
  }

  return response.json();
}

export function isValidURL(url: string): boolean {
  try {
    const urlObj = new URL(url);
    return urlObj.protocol === 'http:' || urlObj.protocol === 'https:';
  } catch {
    return false;
  }
}

