import React, { useRef, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import {
  InputBase,
  Box,
  IconButton,
  Tooltip,
  Typography,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';

const HeaderSearch = ({ value, onChange }) => {
  const theme = useTheme();
  const inputRef = useRef(null);
  const [isFocused, setIsFocused] = useState(false);

  // Keyboard shortcut: Cmd/Ctrl + K to focus search
  useEffect(() => {
    const handleKeyDown = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        inputRef.current?.focus();
      }
      // Escape to blur
      if (e.key === 'Escape' && document.activeElement === inputRef.current) {
        inputRef.current?.blur();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  const handleClear = () => {
    onChange('');
    inputRef.current?.focus();
  };

  const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
  const shortcutKey = isMac ? '⌘K' : 'Ctrl+K';

  return (
    <Box
      sx={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        backgroundColor: isFocused
          ? theme.palette.background.paper
          : theme.palette.mode === 'dark'
          ? 'rgba(255, 255, 255, 0.05)'
          : 'rgba(0, 0, 0, 0.04)',
        borderRadius: 2,
        border: `1px solid ${
          isFocused ? theme.palette.primary.main : 'transparent'
        }`,
        transition: 'all 0.2s ease',
        width: { xs: isFocused ? 240 : 180, sm: isFocused ? 320 : 240 },
        '&:hover': {
          backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.08)'
            : 'rgba(0, 0, 0, 0.06)',
        },
      }}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          pl: 1.5,
          pr: 0.5,
          color: theme.palette.text.secondary,
        }}
      >
        <SearchIcon fontSize="small" />
      </Box>

      <InputBase
        inputRef={inputRef}
        placeholder="Search apps..."
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        sx={{
          flex: 1,
          py: 1,
          pr: 1,
          fontSize: '0.875rem',
          color: theme.palette.text.primary,
          '& input::placeholder': {
            color: theme.palette.text.secondary,
            opacity: 1,
          },
        }}
        inputProps={{
          'aria-label': 'Search applications',
        }}
      />

      {/* Clear button or keyboard shortcut hint */}
      {value ? (
        <Tooltip title="Clear search">
          <IconButton
            size="small"
            onClick={handleClear}
            sx={{
              mr: 0.5,
              color: theme.palette.text.secondary,
              '&:hover': {
                color: theme.palette.text.primary,
              },
            }}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      ) : (
        <Typography
          variant="caption"
          sx={{
            mr: 1.5,
            px: 0.75,
            py: 0.25,
            borderRadius: 0.5,
            backgroundColor: theme.palette.mode === 'dark'
              ? 'rgba(255, 255, 255, 0.1)'
              : 'rgba(0, 0, 0, 0.06)',
            color: theme.palette.text.secondary,
            fontSize: '0.7rem',
            fontWeight: 500,
            display: { xs: 'none', sm: 'block' },
            whiteSpace: 'nowrap',
          }}
        >
          {shortcutKey}
        </Typography>
      )}
    </Box>
  );
};

HeaderSearch.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
};

export default HeaderSearch;
