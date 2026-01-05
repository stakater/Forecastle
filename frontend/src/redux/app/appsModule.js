import { createSlice } from "@reduxjs/toolkit";
import { getApps } from "../../services/api";
import { groupBy } from "../../utils/utils";

const initialState = {
  data: [],
  isLoading: true,
  isLoaded: false,
  error: null,
  lastUpdated: null
};

// Deep compare two objects for equality
const isEqual = (a, b) => JSON.stringify(a) === JSON.stringify(b);

const appsSlice = createSlice({
  name: "apps",
  initialState,
  reducers: {
    loading: state => ({
      ...state,
      isLoading: true,
      error: null
    }),
    loadAppsSuccess: (state, action) => ({
      ...state,
      data: action.payload,
      isLoading: false,
      isLoaded: true,
      lastUpdated: Date.now()
    }),
    // Silent update - no loading state, always updates lastUpdated on successful sync
    refreshAppsSuccess: (state, action) => {
      const dataChanged = !isEqual(state.data, action.payload);
      return {
        ...state,
        data: dataChanged ? action.payload : state.data,
        lastUpdated: Date.now()
      };
    },
    fail: (state, action) => ({
      ...state,
      error: action.payload,
      isLoading: false,
      isLoaded: false
    }),
    // Silent fail for background refresh - don't overwrite existing data
    refreshFail: (state) => state
  }
});

// Extract the actions creators object and reducer
const { actions, reducer } = appsSlice;

// Extract action creators by their names
const { loading, loadAppsSuccess, refreshAppsSuccess, fail, refreshFail } = actions;

// Helper to extract error message
const getErrorMessage = (e) => {
  if (e.response) {
    const status = e.response.status;
    const statusText = e.response.statusText || 'Error';

    if (status === 404) {
      return 'API endpoint not found. Make sure the backend is running.';
    } else if (status >= 500) {
      return `Server error (${status}): ${statusText}`;
    } else {
      return `Request failed (${status}): ${statusText}`;
    }
  } else if (e.request) {
    return 'No response from server. Check if the backend is running.';
  } else if (e.message) {
    return e.message;
  }
  return 'An unexpected error occurred';
};

// Initial load - shows loading state
const loadApps = () => async dispatch => {
  try {
    dispatch(loading());
    let { data } = await getApps();
    data = groupBy("group")(data);
    dispatch(loadAppsSuccess(data));
  } catch (e) {
    dispatch(fail(getErrorMessage(e)));
  }
};

// Background refresh - silent, only updates if data changed
const refreshApps = () => async dispatch => {
  try {
    let { data } = await getApps();
    data = groupBy("group")(data);
    dispatch(refreshAppsSuccess(data));
  } catch (e) {
    // Silent fail - don't disrupt the UI on background refresh errors
    dispatch(refreshFail());
  }
};

// Export required thunks
export { loadApps, refreshApps };

// Export reducer as default
export default reducer;
