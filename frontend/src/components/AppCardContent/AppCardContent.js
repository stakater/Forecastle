import React from "react";
import PropTypes from "prop-types";
import { Grid, CardActionArea, CardMedia, makeStyles } from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  mediaWrapper: {
    backgroundColor: "#FBFBFB",
    padding: `${theme.spacing(2)}px ${theme.spacing(4)}px`
  },
  media: {
    height: 0,
    paddingTop: "35%",
    backgroundSize: "contain"
  }
}));

const AppCardContent = ({ icon, name, url }) => {
  const classes = useStyles();

  return (
    <CardActionArea href={url} target="_blank">
      <Grid className={classes.mediaWrapper}>
        <CardMedia className={classes.media} image={icon} title={name} />
      </Grid>
    </CardActionArea>
  );
};

AppCardContent.propTypes = {
  icon: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired
};

export default AppCardContent;
