import { configureStore } from "redux-starter-kit";

import appsReducer from "./app/appsModule";
import configReducer from "./app/configModule";
import filtersReducer from "./filters/filtersModule";

const store = configureStore({
  reducer: {
    apps: appsReducer,
    filters: filtersReducer,
    config: configReducer
  }
});

export default store;
