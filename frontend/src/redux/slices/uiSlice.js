import { createSlice } from '@reduxjs/toolkit';

// Detect system preference for dark mode
const getSystemThemePreference = () => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }
  return 'light';
};

// Get initial theme from localStorage or system preference
const getInitialTheme = () => {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('forecastle-theme');
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch {
        return getSystemThemePreference();
      }
    }
  }
  return getSystemThemePreference();
};

// Get initial view mode from localStorage
const getInitialViewMode = () => {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('forecastle-view-mode');
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch {
        return 'grid';
      }
    }
  }
  return 'grid';
};

const initialState = {
  themeMode: getInitialTheme(),
  viewMode: getInitialViewMode(), // 'grid' or 'list'
  sidebarOpen: false,
  expandedCards: {}, // { [appId]: true } - tracks expanded cards
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    setThemeMode: (state, action) => {
      state.themeMode = action.payload;
      // Persist to localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem('forecastle-theme', JSON.stringify(action.payload));
      }
    },
    toggleThemeMode: (state) => {
      const newMode = state.themeMode === 'light' ? 'dark' : 'light';
      state.themeMode = newMode;
      // Persist to localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem('forecastle-theme', JSON.stringify(newMode));
      }
    },
    setViewMode: (state, action) => {
      state.viewMode = action.payload;
      // Persist to localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem('forecastle-view-mode', JSON.stringify(action.payload));
      }
    },
    toggleViewMode: (state) => {
      const newMode = state.viewMode === 'grid' ? 'list' : 'grid';
      state.viewMode = newMode;
      // Persist to localStorage
      if (typeof window !== 'undefined') {
        localStorage.setItem('forecastle-view-mode', JSON.stringify(newMode));
      }
    },
    setSidebarOpen: (state, action) => {
      state.sidebarOpen = action.payload;
    },
    toggleSidebar: (state) => {
      state.sidebarOpen = !state.sidebarOpen;
    },
    toggleCardExpanded: (state, action) => {
      const appId = action.payload;
      if (state.expandedCards[appId]) {
        delete state.expandedCards[appId];
      } else {
        state.expandedCards[appId] = true;
      }
    },
    setCardExpanded: (state, action) => {
      const { appId, expanded } = action.payload;
      if (expanded) {
        state.expandedCards[appId] = true;
      } else {
        delete state.expandedCards[appId];
      }
    },
  },
});

export const {
  setThemeMode,
  toggleThemeMode,
  setViewMode,
  toggleViewMode,
  setSidebarOpen,
  toggleSidebar,
  toggleCardExpanded,
  setCardExpanded,
} = uiSlice.actions;

// Selectors
export const selectThemeMode = (state) => state.ui.themeMode;
export const selectViewMode = (state) => state.ui.viewMode;
export const selectSidebarOpen = (state) => state.ui.sidebarOpen;
export const selectExpandedCards = (state) => state.ui.expandedCards;
export const selectIsCardExpanded = (appId) => (state) => !!state.ui.expandedCards[appId];

export default uiSlice.reducer;
