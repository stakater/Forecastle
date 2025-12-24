import React, { createContext, useContext, useMemo } from 'react';
import { createTheme, ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { lightPalette, darkPalette } from './palette';
import { typography } from './typography';
import { transitions } from './transitions';

// Theme mode context for switching between light and dark
export const ThemeModeContext = createContext({
  mode: 'light',
  toggleMode: () => {},
  setMode: () => {},
});

// Hook to access theme mode
export const useThemeMode = () => {
  const context = useContext(ThemeModeContext);
  if (!context) {
    throw new Error('useThemeMode must be used within a ThemeProvider');
  }
  return context;
};

// Create MUI theme based on mode
const createAppTheme = (mode) => {
  const palette = mode === 'dark' ? darkPalette : lightPalette;

  return createTheme({
    palette,
    typography,
    transitions: {
      duration: transitions.duration,
      easing: transitions.easing,
    },
    shape: {
      borderRadius: 8,
    },
    shadows: [
      'none',
      `0 1px 2px ${palette.app.cardShadow}`,
      `0 1px 3px ${palette.app.cardShadow}`,
      `0 2px 4px ${palette.app.cardShadow}`,
      `0 2px 8px ${palette.app.cardShadow}`,
      `0 4px 12px ${palette.app.cardShadow}`,
      `0 6px 16px ${palette.app.cardShadow}`,
      `0 8px 24px ${palette.app.cardShadow}`,
      `0 12px 32px ${palette.app.cardShadow}`,
      ...Array(16).fill(`0 16px 48px ${palette.app.cardShadow}`),
    ],
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            scrollbarColor: mode === 'dark' ? '#475569 #1e293b' : '#cbd5e1 #f1f5f9',
            '&::-webkit-scrollbar': {
              width: '8px',
              height: '8px',
            },
            '&::-webkit-scrollbar-track': {
              background: mode === 'dark' ? '#1e293b' : '#f1f5f9',
            },
            '&::-webkit-scrollbar-thumb': {
              background: mode === 'dark' ? '#475569' : '#cbd5e1',
              borderRadius: '4px',
            },
            '&::-webkit-scrollbar-thumb:hover': {
              background: mode === 'dark' ? '#64748b' : '#94a3b8',
            },
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            textTransform: 'none',
            fontWeight: 500,
            borderRadius: '8px',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: '12px',
            boxShadow: `0 1px 3px ${palette.app.cardShadow}`,
            transition: `all ${transitions.duration.shorter}ms ${transitions.easing.easeOut}`,
            '&:hover': {
              boxShadow: `0 4px 12px ${palette.app.cardShadowHover}`,
            },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            fontWeight: 500,
            fontSize: '0.75rem',
          },
          sizeSmall: {
            height: '24px',
          },
        },
      },
      MuiAccordion: {
        styleOverrides: {
          root: {
            borderRadius: '8px',
            '&:before': {
              display: 'none',
            },
            '&.Mui-expanded': {
              margin: 0,
            },
          },
        },
      },
      MuiAccordionSummary: {
        styleOverrides: {
          root: {
            borderRadius: '8px',
            minHeight: '56px',
            '&.Mui-expanded': {
              minHeight: '56px',
            },
          },
          content: {
            '&.Mui-expanded': {
              margin: '12px 0',
            },
          },
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            transition: `all ${transitions.duration.shortest}ms ${transitions.easing.easeOut}`,
          },
        },
      },
      MuiInputBase: {
        styleOverrides: {
          root: {
            borderRadius: '8px',
          },
        },
      },
      MuiTooltip: {
        styleOverrides: {
          tooltip: {
            fontSize: '0.75rem',
            borderRadius: '6px',
          },
        },
      },
    },
  });
};

// Theme provider component
export const ThemeProvider = ({ children, mode, onModeChange }) => {
  const theme = useMemo(() => createAppTheme(mode), [mode]);

  const contextValue = useMemo(
    () => ({
      mode,
      toggleMode: () => onModeChange(mode === 'light' ? 'dark' : 'light'),
      setMode: onModeChange,
    }),
    [mode, onModeChange]
  );

  return (
    <ThemeModeContext.Provider value={contextValue}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeModeContext.Provider>
  );
};

export default ThemeProvider;
