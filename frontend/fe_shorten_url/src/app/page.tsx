import { Container, Typography, Box } from '@mui/material';
import ShortenURLForm from '../components/ShortenURLForm';
import DecodeURLForm from '../components/DecodeURLForm';

export default function Home() {
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Box sx={{ mb: 4, textAlign: 'center' }}>
        <Typography variant="h3" component="h1" gutterBottom>
          URL Shortener Service
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Shorten your long URLs or decode shortened URLs
        </Typography>
      </Box>

      <ShortenURLForm />
      <DecodeURLForm />
    </Container>
  );
}
