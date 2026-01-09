import React from 'react';
import PropTypes from 'prop-types';
import { Box, Typography, Button } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import RefreshIcon from '@mui/icons-material/Refresh';

const ErrorState = ({
  title = 'Something went wrong',
  description = 'We encountered an error while loading the data',
  error,
  onRetry,
}) => {
  const theme = useTheme();

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        py: 8,
        px: 3,
        textAlign: 'center',
      }}
    >
      <Box
        sx={{
          width: 80,
          height: 80,
          borderRadius: '50%',
          backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(239, 68, 68, 0.15)'
            : 'rgba(220, 38, 38, 0.1)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          mb: 3,
        }}
      >
        <ErrorOutlineIcon
          sx={{
            fontSize: 40,
            color: theme.palette.error.main,
          }}
        />
      </Box>

      <Typography
        variant="h6"
        sx={{
          fontWeight: 600,
          color: theme.palette.text.primary,
          mb: 1,
        }}
      >
        {title}
      </Typography>

      <Typography
        variant="body2"
        sx={{
          color: theme.palette.text.secondary,
          maxWidth: 400,
          mb: error ? 2 : 3,
        }}
      >
        {description}
      </Typography>

      {error && (
        <Typography
          variant="caption"
          sx={{
            color: theme.palette.error.main,
            backgroundColor: theme.palette.mode === 'dark'
              ? 'rgba(239, 68, 68, 0.1)'
              : 'rgba(220, 38, 38, 0.05)',
            px: 2,
            py: 1,
            borderRadius: 1,
            maxWidth: 500,
            mb: 3,
            fontFamily: 'monospace',
            wordBreak: 'break-word',
          }}
        >
          {(() => {
            // Handle different error formats
            if (typeof error === 'string') {
              // Check if it's HTML and extract text
              if (error.includes('<') && error.includes('>')) {
                const match = error.match(/<pre>(.*?)<\/pre>/);
                return match ? match[1] : 'Server returned an HTML error';
              }
              return error;
            }
            return error?.message || 'Unknown error';
          })()}
        </Typography>
      )}

      {onRetry && (
        <Button
          variant="contained"
          onClick={onRetry}
          startIcon={<RefreshIcon />}
          sx={{
            textTransform: 'none',
            fontWeight: 500,
          }}
        >
          Try again
        </Button>
      )}
    </Box>
  );
};

ErrorState.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
  error: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
  onRetry: PropTypes.func,
};

export default ErrorState;
