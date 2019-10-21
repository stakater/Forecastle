import { createSlice } from "redux-starter-kit";

const initialState = {
  query: ""
};

const filtersSlice = createSlice({
  slice: "filters",
  initialState,
  reducers: {
    setQuery: (state, action) => ({
      ...state,
      query: action.payload
    })
  }
});

const { actions, reducer } = filtersSlice;

const setQuery = query => dispatch => {
  dispatch(actions.setQuery(query));
};

export { setQuery };

export default reducer;
