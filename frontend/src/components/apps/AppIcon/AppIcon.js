import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Box, Avatar } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import AppsIcon from '@mui/icons-material/Apps';

const AppIcon = ({ src, alt, size = 48 }) => {
  const theme = useTheme();
  const [hasError, setHasError] = useState(false);
  const [isLoaded, setIsLoaded] = useState(false);

  // Generate a consistent color based on the app name
  const getAvatarColor = (name) => {
    const colors = [
      '#3b82f6', '#8b5cf6', '#ec4899', '#f43f5e',
      '#f97316', '#eab308', '#22c55e', '#14b8a6',
      '#06b6d4', '#6366f1',
    ];
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  };

  const showFallback = !src || hasError;

  return (
    <Box
      sx={{
        width: size,
        height: size,
        borderRadius: 2,
        overflow: 'hidden',
        backgroundColor: theme.palette.mode === 'dark'
          ? 'rgba(255, 255, 255, 0.05)'
          : 'rgba(0, 0, 0, 0.03)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        flexShrink: 0,
      }}
    >
      {showFallback ? (
        <Avatar
          sx={{
            width: size,
            height: size,
            backgroundColor: getAvatarColor(alt || 'App'),
            borderRadius: 2,
            fontSize: size * 0.4,
            fontWeight: 600,
          }}
        >
          {alt ? alt.charAt(0).toUpperCase() : <AppsIcon />}
        </Avatar>
      ) : (
        <Box
          component="img"
          src={src}
          alt={alt}
          onError={() => setHasError(true)}
          onLoad={() => setIsLoaded(true)}
          sx={{
            width: '100%',
            height: '100%',
            objectFit: 'contain',
            opacity: isLoaded ? 1 : 0,
            transition: 'opacity 0.2s ease',
          }}
        />
      )}
    </Box>
  );
};

AppIcon.propTypes = {
  src: PropTypes.string,
  alt: PropTypes.string,
  size: PropTypes.number,
};

export default AppIcon;
