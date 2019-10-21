import { createSlice } from 'redux-starter-kit';
import { getConfig } from '../../services/api'

const initialState = {
    data: [],
    isLoading: true,
    isLoaded: false,
    error: null
  };

const configSlice = createSlice({
    slice: "config",
    initialState,
    reducers: {
        loading: state => ({
            ...state,
            isLoading: true,
            error: null
        }),
        loadConfigSuccess: (state, action) => ({
            ...state,
            data: action.payload,
            isLoading: false,
            isLoaded: true
        }),
        fail: (state, action) => ({
            ...state,
            error: action.payload,
            isLoading: false,
            isLoaded: true
        })
    }
});

// Extract the actions creators object and reducer
const { actions, reducer } = configSlice;

// Extract action creators by their names
const { loading, loadConfigSuccess, fail } = actions;

const loadConfig = () => async dispatch => {
    try {
        dispatch(loading());
        let { data } = await getConfig();

        dispatch(loadConfigSuccess(data));
    } catch(e) {
        if(e.response && e.response.data) {
            dispatch(fail(e.response.data));
        } else {
            dispatch(fail(e.message));
        }
    }
};

// Export required thunks
export { loadConfig };

// Export reducer as default
export default reducer;
