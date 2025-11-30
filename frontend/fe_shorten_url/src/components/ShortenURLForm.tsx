'use client';

import { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Paper,
  Typography,
  Alert,
  CircularProgress,
  Link,
} from '@mui/material';
import { ContentCopy as ContentCopyIcon } from '@mui/icons-material';
import { submitURL, getByLongURL, isValidURL, type ShortenURLResponse } from '../lib/api';

const POLL_INTERVAL = 2000; // 2 seconds
const MAX_POLL_ATTEMPTS = 30; // 1 minute max

export default function ShortenURLForm() {
  const [longURL, setLongURL] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [polling, setPolling] = useState(false);
  const [result, setResult] = useState<ShortenURLResponse | null>(null);
  const [shortURL, setShortURL] = useState('');
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setResult(null);
    setShortURL('');
    setCopied(false);

    // Validate URL
    if (!longURL.trim()) {
      setError('Please enter a URL');
      return;
    }

    if (!isValidURL(longURL)) {
      setError('Please enter a valid URL (must start with http:// or https://)');
      return;
    }

    setLoading(true);

    try {
      // Submit URL
      await submitURL(longURL);
      
      // Start polling
      setPolling(true);
      let attempts = 0;

      const poll = async () => {
        try {
          const response = await getByLongURL(longURL);
          
          if (response.data && response.data.status === 'encoded') {
            setResult(response.data);
            const redirectHost = process.env.NEXT_PUBLIC_REDIRECT_HOST || 'http://localhost:8080';
            // Remove trailing slash if present
            const cleanHost = redirectHost.replace(/\/$/, '');
            setShortURL(`${cleanHost}/${response.data.code}`);
            setPolling(false);
            setLoading(false);
            return;
          }

          attempts++;
          if (attempts >= MAX_POLL_ATTEMPTS) {
            setError('Timeout: URL encoding is taking longer than expected. Please try again later.');
            setPolling(false);
            setLoading(false);
            return;
          }

          // Continue polling
          setTimeout(poll, POLL_INTERVAL);
        } catch (err: any) {
          // If 404, continue polling (URL not ready yet)
          if (err.status === 404 || err.message?.includes('404') || err.message?.includes('not found')) {
            attempts++;
            if (attempts >= MAX_POLL_ATTEMPTS) {
              setError('Timeout: URL encoding is taking longer than expected. Please try again later.');
              setPolling(false);
              setLoading(false);
              return;
            }
            setTimeout(poll, POLL_INTERVAL);
          } else {
            setError(err.message || 'Failed to get shortened URL');
            setPolling(false);
            setLoading(false);
          }
        }
      };

      // Start polling after a short delay
      setTimeout(poll, POLL_INTERVAL);
    } catch (err: any) {
      setError(err.message || 'Failed to submit URL');
      setLoading(false);
      setPolling(false);
    }
  };

  const handleCopy = async () => {
    if (shortURL) {
      try {
        await navigator.clipboard.writeText(shortURL);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch (err) {
        setError('Failed to copy to clipboard');
      }
    }
  };


  return (
    <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
      <Typography variant="h5" gutterBottom>
        Shorten URL
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Enter a long URL to get a shortened version
      </Typography>

      <Box component="form" onSubmit={handleSubmit} sx={{ mb: 2 }}>
        <TextField
          fullWidth
          label="Long URL"
          value={longURL}
          onChange={(e) => setLongURL(e.target.value)}
          placeholder="https://example.com/very/long/url"
          error={!!error && !loading && !polling}
          disabled={loading || polling}
          sx={{ mb: 2 }}
        />

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        <Button
          type="submit"
          variant="contained"
          fullWidth
          disabled={loading || polling}
          startIcon={loading || polling ? <CircularProgress size={20} /> : null}
        >
          {loading ? 'Submitting...' : polling ? 'Processing...' : 'Shorten URL'}
        </Button>
      </Box>

      {shortURL && result && (
        <Box sx={{ mt: 3, p: 2, bgcolor: 'success.light', borderRadius: 1 }}>
          <Typography variant="subtitle2" gutterBottom>
            Shortened URL:
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Link href={shortURL} target="_blank" rel="noopener noreferrer" sx={{ flex: 1 }}>
              {shortURL}
            </Link>
            <Button
              size="small"
              variant="outlined"
              startIcon={<ContentCopyIcon />}
              onClick={handleCopy}
              sx={{ minWidth: 'auto' }}
            >
              {copied ? 'Copied!' : 'Copy'}
            </Button>
          </Box>
        </Box>
      )}
    </Paper>
  );
}

