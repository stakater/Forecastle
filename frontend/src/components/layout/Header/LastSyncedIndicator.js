import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { Box, Typography, Tooltip } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import SyncIcon from '@mui/icons-material/Sync';

const formatElapsedTime = (seconds) => {
  if (seconds < 5) return 'Just now';
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  return `${hours}h ago`;
};

const LastSyncedIndicator = () => {
  const theme = useTheme();
  const lastUpdated = useSelector((state) => state.apps.lastUpdated);
  const [elapsed, setElapsed] = useState(0);

  useEffect(() => {
    if (!lastUpdated) return;

    // Calculate initial elapsed time
    setElapsed(Math.floor((Date.now() - lastUpdated) / 1000));

    // Update every second
    const interval = setInterval(() => {
      setElapsed(Math.floor((Date.now() - lastUpdated) / 1000));
    }, 1000);

    return () => clearInterval(interval);
  }, [lastUpdated]);

  // Don't render until first sync
  if (!lastUpdated) {
    return (
      <Box
        sx={{
          display: { xs: 'none', md: 'flex' },
          alignItems: 'center',
          gap: 0.5,
        }}
      >
        <SyncIcon
          sx={{
            fontSize: 14,
            color: theme.palette.text.disabled,
            animation: 'spin 1s linear infinite',
            '@keyframes spin': {
              '0%': { transform: 'rotate(0deg)' },
              '100%': { transform: 'rotate(360deg)' },
            },
          }}
        />
        <Typography
          variant="caption"
          sx={{
            color: theme.palette.text.disabled,
            fontSize: '0.75rem',
          }}
        >
          Syncing...
        </Typography>
      </Box>
    );
  }

  const formattedTime = formatElapsedTime(elapsed);
  const absoluteTime = new Date(lastUpdated).toLocaleTimeString();

  return (
    <Tooltip title={`Last synced at ${absoluteTime}`}>
      <Box
        sx={{
          display: { xs: 'none', md: 'flex' },
          alignItems: 'center',
          justifyContent: 'center',
          gap: 0.5,
          px: 1.5,
          py: 0.5,
          width: 95, // Fixed width to prevent layout shift
          backgroundColor:
            theme.palette.mode === 'dark'
              ? 'rgba(255, 255, 255, 0.05)'
              : 'rgba(0, 0, 0, 0.04)',
          borderRadius: 1.5,
        }}
      >
        <SyncIcon
          sx={{
            fontSize: 14,
            color: theme.palette.text.secondary,
          }}
        />
        <Typography
          variant="caption"
          sx={{
            color: theme.palette.text.secondary,
            fontSize: '0.75rem',
            whiteSpace: 'nowrap',
            fontVariantNumeric: 'tabular-nums', // Monospace numbers
          }}
        >
          {formattedTime}
        </Typography>
      </Box>
    </Tooltip>
  );
};

export default LastSyncedIndicator;
