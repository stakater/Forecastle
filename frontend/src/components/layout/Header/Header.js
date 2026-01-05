import React, { useState, useEffect, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Container,
  useScrollTrigger,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';

import HeaderSearch from './HeaderSearch';
import HeaderControls from './HeaderControls';
import { setQuery } from '../../../redux/filters/filtersModule';

const Header = () => {
  const theme = useTheme();
  const dispatch = useDispatch();
  const config = useSelector((state) => state.config.data);
  const query = useSelector((state) => state.filters.query);
  const [isScrolled, setIsScrolled] = useState(false);

  // Detect scroll for backdrop blur effect
  const trigger = useScrollTrigger({
    disableHysteresis: true,
    threshold: 10,
  });

  useEffect(() => {
    setIsScrolled(trigger);
  }, [trigger]);

  // Handle search input
  const handleSearch = useCallback(
    (value) => {
      dispatch(setQuery(value));
    },
    [dispatch]
  );

  // Set document title based on config
  useEffect(() => {
    if (config?.title) {
      document.title = config.title;
    }
  }, [config?.title]);

  const title = config?.title || 'Forecastle';

  return (
    <AppBar
      position="sticky"
      elevation={isScrolled ? 1 : 0}
      sx={{
        backgroundColor: isScrolled
          ? theme.palette.mode === 'dark'
            ? 'rgba(30, 41, 59, 0.9)'
            : 'rgba(255, 255, 255, 0.9)'
          : theme.palette.background.paper,
        backdropFilter: isScrolled ? 'blur(10px)' : 'none',
        borderBottom: `1px solid ${theme.palette.divider}`,
        transition: 'all 0.2s ease',
      }}
    >
      <Container maxWidth="xl">
        <Toolbar
          disableGutters
          sx={{
            minHeight: { xs: 56, md: 64 },
            gap: 2,
          }}
        >
          {/* Logo and Title */}
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 1.5,
              flexShrink: 0,
            }}
          >
            <Box
              component="img"
              src="favicon.ico"
              alt="Forecastle"
              sx={{
                width: 32,
                height: 32,
                display: { xs: 'none', sm: 'block' },
              }}
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            <Typography
              variant="h6"
              component="h1"
              sx={{
                fontWeight: 600,
                color: theme.palette.text.primary,
                display: { xs: 'none', sm: 'block' },
                whiteSpace: 'nowrap',
              }}
            >
              {title}
            </Typography>
          </Box>

          {/* Spacer */}
          <Box sx={{ flexGrow: 1 }} />

          {/* Search */}
          <HeaderSearch value={query} onChange={handleSearch} />

          {/* Controls */}
          <HeaderControls />
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Header;
