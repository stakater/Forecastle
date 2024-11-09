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

    // todo: move to utils
    data.sort((a, b) => a.name.toLowerCase().localeCompare(b.name.toLowerCase(), 'en', {numeric: true}));
    
    let groups = [...new Set(data.map(i => i.group))];
 
    data = { groups, apps: data };

    dispatch(loadAppsSuccess(data));
  } catch (e) {
    if (e.response && e.response.data) {
      dispatch(fail(e.response.data));
    } else {
      dispatch(fail(e.message));
    }
  }
};

// Export required thunks
export { loadApps };

// Export reducer as default
export default reducer;
