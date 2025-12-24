// Color palette for Forecastle UI
// Modern minimal design with light and dark theme support

export const lightPalette = {
  mode: 'light',
  primary: {
    main: '#2563eb',
    light: '#3b82f6',
    dark: '#1d4ed8',
    contrastText: '#ffffff',
  },
  secondary: {
    main: '#7c3aed',
    light: '#8b5cf6',
    dark: '#6d28d9',
    contrastText: '#ffffff',
  },
  background: {
    default: '#f8fafc',
    paper: '#ffffff',
  },
  surface: {
    main: '#ffffff',
    hover: '#f1f5f9',
    active: '#e2e8f0',
  },
  text: {
    primary: '#0f172a',
    secondary: '#475569',
    disabled: '#94a3b8',
  },
  divider: '#e2e8f0',
  error: {
    main: '#dc2626',
    light: '#ef4444',
    dark: '#b91c1c',
  },
  warning: {
    main: '#f59e0b',
    light: '#fbbf24',
    dark: '#d97706',
  },
  success: {
    main: '#10b981',
    light: '#34d399',
    dark: '#059669',
  },
  info: {
    main: '#0ea5e9',
    light: '#38bdf8',
    dark: '#0284c7',
  },
  // Custom colors for app-specific elements
  app: {
    cardShadow: 'rgba(0, 0, 0, 0.08)',
    cardShadowHover: 'rgba(0, 0, 0, 0.12)',
    headerBlur: 'rgba(255, 255, 255, 0.8)',
    iconBackground: '#f1f5f9',
  },
};

export const darkPalette = {
  mode: 'dark',
  primary: {
    main: '#3b82f6',
    light: '#60a5fa',
    dark: '#2563eb',
    contrastText: '#ffffff',
  },
  secondary: {
    main: '#8b5cf6',
    light: '#a78bfa',
    dark: '#7c3aed',
    contrastText: '#ffffff',
  },
  background: {
    default: '#0f172a',
    paper: '#1e293b',
  },
  surface: {
    main: '#1e293b',
    hover: '#334155',
    active: '#475569',
  },
  text: {
    primary: '#f8fafc',
    secondary: '#cbd5e1',
    disabled: '#64748b',
  },
  divider: '#334155',
  error: {
    main: '#ef4444',
    light: '#f87171',
    dark: '#dc2626',
  },
  warning: {
    main: '#fbbf24',
    light: '#fcd34d',
    dark: '#f59e0b',
  },
  success: {
    main: '#34d399',
    light: '#6ee7b7',
    dark: '#10b981',
  },
  info: {
    main: '#38bdf8',
    light: '#7dd3fc',
    dark: '#0ea5e9',
  },
  // Custom colors for app-specific elements
  app: {
    cardShadow: 'rgba(0, 0, 0, 0.3)',
    cardShadowHover: 'rgba(0, 0, 0, 0.4)',
    headerBlur: 'rgba(15, 23, 42, 0.8)',
    iconBackground: '#334155',
  },
};
