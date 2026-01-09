import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  Box,
  IconButton,
  Tooltip,
  Divider,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import GridViewIcon from '@mui/icons-material/GridView';
import ViewListIcon from '@mui/icons-material/ViewList';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';

import {
  selectViewMode,
  selectThemeMode,
  setViewMode,
  toggleThemeMode,
} from '../../../redux/slices/uiSlice';

const HeaderControls = () => {
  const theme = useTheme();
  const dispatch = useDispatch();
  const viewMode = useSelector(selectViewMode);
  const themeMode = useSelector(selectThemeMode);

  const handleViewModeChange = (mode) => {
    dispatch(setViewMode(mode));
  };

  const handleThemeToggle = () => {
    dispatch(toggleThemeMode());
  };

  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        gap: 0.5,
      }}
    >
      {/* View Mode Toggle */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.05)'
            : 'rgba(0, 0, 0, 0.04)',
          borderRadius: 1.5,
          p: 0.5,
        }}
      >
        <Tooltip title="Grid view">
          <IconButton
            size="small"
            onClick={() => handleViewModeChange('grid')}
            sx={{
              color: viewMode === 'grid'
                ? theme.palette.primary.main
                : theme.palette.text.secondary,
              backgroundColor: viewMode === 'grid'
                ? theme.palette.mode === 'dark'
                  ? 'rgba(59, 130, 246, 0.2)'
                  : 'rgba(37, 99, 235, 0.1)'
                : 'transparent',
              borderRadius: 1,
              '&:hover': {
                backgroundColor: viewMode === 'grid'
                  ? theme.palette.mode === 'dark'
                    ? 'rgba(59, 130, 246, 0.3)'
                    : 'rgba(37, 99, 235, 0.15)'
                  : theme.palette.action.hover,
              },
            }}
          >
            <GridViewIcon fontSize="small" />
          </IconButton>
        </Tooltip>

        <Tooltip title="List view">
          <IconButton
            size="small"
            onClick={() => handleViewModeChange('list')}
            sx={{
              color: viewMode === 'list'
                ? theme.palette.primary.main
                : theme.palette.text.secondary,
              backgroundColor: viewMode === 'list'
                ? theme.palette.mode === 'dark'
                  ? 'rgba(59, 130, 246, 0.2)'
                  : 'rgba(37, 99, 235, 0.1)'
                : 'transparent',
              borderRadius: 1,
              '&:hover': {
                backgroundColor: viewMode === 'list'
                  ? theme.palette.mode === 'dark'
                    ? 'rgba(59, 130, 246, 0.3)'
                    : 'rgba(37, 99, 235, 0.15)'
                  : theme.palette.action.hover,
              },
            }}
          >
            <ViewListIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>

      <Divider
        orientation="vertical"
        flexItem
        sx={{
          mx: 1,
          height: 24,
          alignSelf: 'center',
        }}
      />

      {/* Theme Toggle */}
      <Tooltip title={themeMode === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}>
        <IconButton
          onClick={handleThemeToggle}
          size="small"
          sx={{
            color: theme.palette.text.secondary,
            '&:hover': {
              color: theme.palette.text.primary,
              backgroundColor: theme.palette.action.hover,
            },
          }}
        >
          {themeMode === 'light' ? (
            <DarkModeIcon fontSize="small" />
          ) : (
            <LightModeIcon fontSize="small" />
          )}
        </IconButton>
      </Tooltip>
    </Box>
  );
};

export default HeaderControls;
