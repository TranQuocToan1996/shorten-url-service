'use client';

import { CacheProvider } from '@emotion/react';
import { ThemeProvider, CssBaseline } from '@mui/material';
import createEmotionCache from './emotion-cache';
import theme from './theme';
import { useState } from 'react';


export default function Providers({ children }: { children: React.ReactNode }) {
    const [emotionCache] = useState(() => createEmotionCache());
    return (
        <CacheProvider value={emotionCache}>
            <ThemeProvider theme={theme}>
                <CssBaseline />
                {children}
            </ThemeProvider>
        </CacheProvider>
    );
}
