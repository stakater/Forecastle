import { createSlice } from "@reduxjs/toolkit";
import { getApps } from "../../services/api";
import { groupBy } from "../../utils/utils";

const initialState = {
  data: [],
  isLoading: true,
  isLoaded: false,
  error: null
};

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
      isLoaded: true
    }),
    fail: (state, action) => ({
      ...state,
      error: action.payload,
      isLoading: false,
      isLoaded: false
    })
  }
});

// Extract the actions creators object and reducer
const { actions, reducer } = appsSlice;

// Extract action creators by their names
const { loading, loadAppsSuccess, fail } = actions;

const loadApps = () => async dispatch => {
  try {
    dispatch(loading());
    let { data } = await getApps();

    data = groupBy("group")(data);

    dispatch(loadAppsSuccess(data));
  } catch (e) {
    // Extract a clean error message
    let errorMessage = 'An unexpected error occurred';

    if (e.response) {
      // Server responded with an error status
      const status = e.response.status;
      const statusText = e.response.statusText || 'Error';

      if (status === 404) {
        errorMessage = 'API endpoint not found. Make sure the backend is running.';
      } else if (status >= 500) {
        errorMessage = `Server error (${status}): ${statusText}`;
      } else {
        errorMessage = `Request failed (${status}): ${statusText}`;
      }
    } else if (e.request) {
      // Request was made but no response received
      errorMessage = 'No response from server. Check if the backend is running.';
    } else if (e.message) {
      errorMessage = e.message;
    }

    dispatch(fail(errorMessage));
  }
};

// Export required thunks
export { loadApps };

// Export reducer as default
export default reducer;
