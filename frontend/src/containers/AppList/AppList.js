import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Box, Grid, Container, makeStyles } from "@material-ui/core";

import * as appsStore from "../../redux/app/appsModule";
import selectApps from "../../redux/app/appsSelector";
import { sortAlphabetically } from "../../utils/utils";

import { AppCard, PageLoader } from "../../components";

import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import Typography from "@material-ui/core/Typography";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

const useStyles = makeStyles(theme => ({
  root: {
    minHeight: "calc(100vh - 118px)"
  },
  expansionPanel: {
    width: "100%",
    boxShadow: "none"
  },
  expansionPanelHeader: {
    backgroundColor: "#f1f1f1"
  },
  cardGrid: {
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(8)
  },
  panelDetails: {
    paddingLeft: 0,
    paddingRight: 0,
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2)
  }
}));

export const AppList = ({ apps, groups, isLoading, isLoaded, error, loadApps }) => {
  const classes = useStyles();

  useEffect(() => {
    loadApps();
  }, [loadApps]);

  const groups = sortAlphabetically(Object.keys(apps))

  return (
    <main className={classes.root}>
      {/* Show loader when fetching apps collections */}
      <PageLoader show={isLoading} />

      {/* Display list of apps  */}
      <Container className={classes.cardGrid} fixed>
        {groups.map(group => (
          <ExpansionPanel
            defaultExpanded
            className={classes.expansionPanel}
            key={group}
          >
            <ExpansionPanelSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
              id="panel1a-header"
              className={classes.expansionPanelHeader}
            >
              <Typography className={classes.heading}>
                {group.toUpperCase()} ({apps[group].length})
              </Typography>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails className={classes.panelDetails}>
              <Grid container spacing={4}>
                {apps[group].map((app, idx) => (
                  <Grid key={idx} item xs={12} sm={6} md={3}>
                    <AppCard card={app} />
                  </Grid>
                ))}
              </Grid>
            </ExpansionPanelDetails>
          </ExpansionPanel>
        ))}

        {/* Display message if result list is empty */}
        {Object.keys(apps).length === 0 && isLoaded && (
          <Box>No results found matching your query</Box>
        )}

        {/* Display error message in case of failed api call or erroneous response */}
        {error && !isLoaded && (
          <Box>Results couldn't be loaded due to an error</Box>
        )}
      </Container>
    </main>
  );
};

AppList.props = {
  apps: PropTypes.object,
  isLoading: PropTypes.bool.isRequired,
  isLoaded: PropTypes.bool.isRequired,
  error: PropTypes.oneOf([PropTypes.string, PropTypes.object])
};

AppList.defaultProps = {
  apps: {},
  isLoading: false,
  isLoaded: false,
  error: null
};

const mapStateToProps = state => ({
  apps: selectApps(state.apps.data, state.filters),
  isLoading: state.apps.isLoading,
  isLoaded: state.apps.isLoaded,
  error: state.apps.error
});

const mapDispatchToProps = dispatch => ({
  loadApps: () => dispatch(appsStore.loadApps())
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AppList);
