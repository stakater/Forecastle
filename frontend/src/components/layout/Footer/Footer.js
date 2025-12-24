import React from 'react';
import { Typography, Link, Box } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';

const Footer = () => {
  const theme = useTheme();

  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: theme.palette.mode === 'dark'
          ? 'rgba(255, 255, 255, 0.02)'
          : 'rgba(0, 0, 0, 0.02)',
        borderTop: `1px solid ${theme.palette.divider}`,
        py: 2,
        px: 3,
      }}
    >
      <Typography
        variant="body2"
        align="center"
        sx={{
          color: theme.palette.text.secondary,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: 0.5,
        }}
      >
        Made with
        <FavoriteBorderIcon
          sx={{
            fontSize: 16,
            color: theme.palette.error.main,
            mx: 0.25,
          }}
        />
        by
        <Link
          href="https://stakater.com/"
          target="_blank"
          rel="noopener noreferrer"
          sx={{
            color: theme.palette.text.secondary,
            textDecoration: 'none',
            fontWeight: 500,
            ml: 0.5,
            '&:hover': {
              color: theme.palette.primary.main,
              textDecoration: 'underline',
            },
          }}
        >
          Stakater
        </Link>
      </Typography>
    </Box>
  );
};

export default Footer;
