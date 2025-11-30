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
import { decodeURL, type ShortenURLResponse } from '../lib/api';

export default function DecodeURLForm() {
  const [shortenURL, setShortenURL] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ShortenURLResponse | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setResult(null);

    if (!shortenURL.trim()) {
      setError('Please enter a shortened URL');
      return;
    }

    // Normalize the URL - add protocol if missing
    let normalizedURL = shortenURL.trim();
    if (!normalizedURL.startsWith('http://') && !normalizedURL.startsWith('https://')) {
      // If it's just a code, prepend default host with /api/v1 prefix
      if (!normalizedURL.includes('/')) {
        const defaultHost = process.env.NEXT_PUBLIC_REDIRECT_HOST || 'http://localhost:8080/api/v1';
        const cleanHost = defaultHost.replace(/\/$/, '');
        normalizedURL = `${cleanHost}/${normalizedURL}`;
      } else {
        normalizedURL = `http://${normalizedURL}`;
      }
    }

    setLoading(true);

    try {
      const response = await decodeURL(normalizedURL);

      if (response.data) {
        setResult(response.data);
      } else {
        setError('No data returned from server');
      }
    } catch (err) {
      const error = err as Error;
      setError(error.message || 'Failed to decode URL');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper elevation={3} sx={{ p: 3 }}>
      <Typography variant="h5" gutterBottom>
        Decode Shortened URL
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Enter a shortened URL to get the original long URL
      </Typography>

      <Box component="form" onSubmit={handleSubmit} sx={{ mb: 2 }}>
        <TextField
          fullWidth
          label="Shortened URL"
          value={shortenURL}
          onChange={(e) => setShortenURL(e.target.value)}
          placeholder="http://localhost:8080/api/v1/abc123"
          error={!!error && !loading}
          disabled={loading}
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
          disabled={loading}
          startIcon={loading ? <CircularProgress size={20} /> : null}
        >
          {loading ? 'Decoding...' : 'Decode URL'}
        </Button>
      </Box>

      {result && (
        <Box sx={{ mt: 3, p: 2, borderRadius: 1, bgcolor: "#E8EFFD" }}>
          <Typography variant="subtitle2" gutterBottom>
            Original Long URL:
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 1 }}>
            <Link href={result.long_url} target="_blank" rel="noopener noreferrer" sx={{ overflowWrap: 'anywhere' }}>
              {result.long_url}
            </Link>
            <Button
              size="small"
              variant="outlined"
              onClick={async () => {
                if (navigator.clipboard) await navigator.clipboard.writeText(result.long_url);
              }}
              sx={{ minWidth: 'auto' }}
            >
              Copy
            </Button>
          </Box>
          <Box sx={{ mt: 2 }}>
            <Typography variant="caption" display="block" color="text.secondary">
              Status: {result.status}
            </Typography>
            <Typography variant="caption" display="block" color="text.secondary">
              Code: {result.code}
            </Typography>
          </Box>
        </Box>
      )}
    </Paper>
  );
}

