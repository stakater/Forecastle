import React, { useEffect } from "react";
import { connect } from "react-redux";

import AppList from "../AppList/AppList";
import { SearchAppBar, Footer } from "../../components";
import * as configStore from "../../redux/app/configModule";

const App = ({ loadConfig }) => {
  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  return (
    <React.Fragment>
      <SearchAppBar />
      <AppList />

      {/* Footer */}
      <Footer />
      {/* End footer */}
    </React.Fragment>
  );
};

const mapDispatchToProps = dispatch => ({
  loadConfig: () => dispatch(configStore.loadConfig())
});

export default connect(
  null,
  mapDispatchToProps
)(App);
