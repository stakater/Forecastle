import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "react-redux";

import * as serviceWorker from "./serviceWorker";
import App from "./containers/App/App";
import store from "./redux";

import "./index.css";

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById("root")
);

serviceWorker.unregister();
