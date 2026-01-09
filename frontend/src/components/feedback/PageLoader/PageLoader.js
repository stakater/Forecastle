import React from 'react';
import PropTypes from 'prop-types';
import { Box, CircularProgress, Typography } from '@mui/material';
import { useTheme } from '@mui/material/styles';

const PageLoader = ({ show, message }) => {
  const theme = useTheme();

  if (!show) return null;

  return (
    <Box
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        zIndex: 9999,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: theme.palette.mode === 'dark'
          ? 'rgba(15, 23, 42, 0.95)'
          : 'rgba(248, 250, 252, 0.95)',
        backdropFilter: 'blur(4px)',
        opacity: show ? 1 : 0,
        visibility: show ? 'visible' : 'hidden',
        transition: 'all 0.25s ease-in-out',
      }}
    >
      <CircularProgress
        size={48}
        thickness={3}
        sx={{
          color: theme.palette.primary.main,
        }}
      />
      {message && (
        <Typography
          variant="body2"
          sx={{
            mt: 2,
            color: theme.palette.text.secondary,
          }}
        >
          {message}
        </Typography>
      )}
    </Box>
  );
};

PageLoader.propTypes = {
  show: PropTypes.bool,
  message: PropTypes.string,
};

PageLoader.defaultProps = {
  show: false,
  message: '',
};

export default PageLoader;
