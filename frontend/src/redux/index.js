import { configureStore } from "@reduxjs/toolkit";

import appsReducer from "./app/appsModule";
import configReducer from "./app/configModule";
import filtersReducer from "./filters/filtersModule";
import uiReducer from "./slices/uiSlice";

const store = configureStore({
  reducer: {
    apps: appsReducer,
    filters: filtersReducer,
    config: configReducer,
    ui: uiReducer,
  },
});

export default store;
