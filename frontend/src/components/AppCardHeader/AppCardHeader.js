import React from "react";
import PropTypes from "prop-types";
import {
  CardHeader,
  Tooltip,
  Typography,
  Grid,
  makeStyles
} from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  header: {
    position: "relative"
  },
  headerAction: {
    position: "absolute",
    right: "0.5rem"
  }
}));
const AppCardHeader = ({ name, url }) => {
  const classes = useStyles();

  return (
    <CardHeader
      className={classes.header}
      title={name}
      titleTypographyProps={{ variant: "subtitle1" }}
      subheader={
        <Tooltip title={url}>
          <Grid item xs zeroMinWidth style={{ maxWidth: "100%" }}>
            <Typography
              variant="subtitle2"
              noWrap
              color="primary"
              component="a"
              href={url}
              target="_blank"
              style={{
                textDecoration: "none",
                width: "100%",
                display: "inline-block"
              }}
            >
              {url}
            </Typography>
          </Grid>
        </Tooltip>
      }
    />
  );
};

AppCardHeader.propTypes = {
  name: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
};

export default AppCardHeader;
