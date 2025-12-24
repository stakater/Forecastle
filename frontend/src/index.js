import React from "react";
import { createRoot } from "react-dom/client";
import { Provider, useSelector, useDispatch } from "react-redux";

import * as serviceWorker from "./serviceWorker";
import App from "./containers/App/App";
import store from "./redux";
import { ThemeProvider } from "./theme";
import { selectThemeMode, setThemeMode } from "./redux/slices/uiSlice";

import "./index.css";

// Root component that connects theme to Redux
const Root = () => {
  const themeMode = useSelector(selectThemeMode);
  const dispatch = useDispatch();

  const handleThemeChange = (mode) => {
    dispatch(setThemeMode(mode));
  };

  return (
    <ThemeProvider mode={themeMode} onModeChange={handleThemeChange}>
      <App />
    </ThemeProvider>
  );
};

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <Provider store={store}>
    <Root />
  </Provider>
);

serviceWorker.unregister();
