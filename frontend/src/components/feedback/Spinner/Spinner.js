import React from 'react';
import PropTypes from 'prop-types';
import { Box, keyframes } from '@mui/material';
import { useTheme } from '@mui/material/styles';

const pulseAnimation = keyframes`
  0% {
    transform: scale(0);
    opacity: 0;
  }
  5% {
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 0;
  }
`;

const Spinner = ({ size = 60 }) => {
  const theme = useTheme();
  const color = theme.palette.primary.main;

  return (
    <Box
      sx={{
        position: 'relative',
        width: size,
        height: size,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      {[0, 1, 2].map((index) => (
        <Box
          key={index}
          sx={{
            position: 'absolute',
            width: size,
            height: size,
            borderRadius: '50%',
            backgroundColor: color,
            opacity: 0,
            animation: `${pulseAnimation} 1s linear infinite`,
            animationDelay: `${-0.4 + index * 0.2}s`,
          }}
        />
      ))}
    </Box>
  );
};

Spinner.propTypes = {
  size: PropTypes.number,
};

export default Spinner;
