import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { AppBar, Toolbar, Typography, InputBase } from "@material-ui/core";
import { connect } from "react-redux";
import { fade, makeStyles } from "@material-ui/core/styles";
import SearchIcon from "@material-ui/icons/Search";

import * as filtersStore from "../../redux/filters/filtersModule";

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1
  },
  menuButton: {
    marginRight: theme.spacing(2)
  },
  title: {
    flexGrow: 1,
    display: "none",
    [theme.breakpoints.up("sm")]: {
      display: "block"
    }
  },
  search: {
    position: "relative",
    borderRadius: theme.shape.borderRadius,
    backgroundColor: fade(theme.palette.common.white, 0.15),
    "&:hover": {
      backgroundColor: fade(theme.palette.common.white, 0.25)
    },
    marginLeft: 0,
    width: "100%",
    [theme.breakpoints.up("sm")]: {
      marginLeft: theme.spacing(1),
      width: "auto"
    }
  },
  searchIcon: {
    width: theme.spacing(7),
    height: "100%",
    position: "absolute",
    pointerEvents: "none",
    display: "flex",
    alignItems: "center",
    justifyContent: "center"
  },
  inputRoot: {
    color: "inherit"
  },
  inputInput: {
    padding: theme.spacing(1, 1, 1, 7),
    transition: theme.transitions.create("width"),
    width: "100%",
    [theme.breakpoints.up("sm")]: {
      width: 120,
      "&:focus": {
        width: 200
      }
    }
  }
}));

const SearchAppBar = ({ config, query, setQuery }) => {
  const classes = useStyles();

  const title = config.title || "Forecastle - Stakater"
  useEffect(() => {
    document.title = title;
  }, [title]);

  const handleSearchInput = e => {
    setQuery(e.target.value);
  };

  return (
    <div className={classes.root}>
      <AppBar
        position="static"
        style={{
          backgroundColor: config.headerBackground || "#3f51b5",
          color: config.headerForeground || "white"
        }}
      >
        <Toolbar>
          <Typography className={classes.title} variant="h6" noWrap>
            {title}
          </Typography>
          <div className={classes.search}>
            <div className={classes.searchIcon}>
              <SearchIcon />
            </div>
            <InputBase
              placeholder="Search…"
              classes={{
                root: classes.inputRoot,
                input: classes.inputInput
              }}
              inputProps={{ "aria-label": "search" }}
              value={query}
              onChange={handleSearchInput}
            />
          </div>
        </Toolbar>
      </AppBar>
    </div>
  );
};

SearchAppBar.props = {
  query: PropTypes.string.isRequired,
  config: PropTypes.shape({
    title: PropTypes.string,
    headerBackground: PropTypes.string,
    headerForeground: PropTypes.string
  })
};

const mapStateToProps = state => ({
  query: state.filters.query,
  config: state.config.data
});

const mapDispatchToProps = dispatch => ({
  setQuery: query => dispatch(filtersStore.setQuery(query))
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(SearchAppBar);
