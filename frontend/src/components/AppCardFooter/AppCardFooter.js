import React from "react";
import clsx from "clsx";
import PropTypes from "prop-types";
import {
  Chip,
  CardActions,
  IconButton,
  Tooltip,
  makeStyles
} from "@material-ui/core";
import VpnLockIcon from "@material-ui/icons/VpnLock";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

const useStyles = makeStyles(theme => ({
  expand: {
    transform: "rotate(0deg)",
    marginLeft: "auto",
    transition: theme.transitions.create("transform", {
      duration: theme.transitions.duration.shortest
    })
  },
  expandOpen: {
    transform: "rotate(180deg)"
  }
}));

const AppCardFooter = ({
  discoverySource,
  networkRestricted,
  properties,
  isDetailsExpanded,
  onExpandDetails
}) => {
  const classes = useStyles();

  return (
    <CardActions disableSpacing>
      <Tooltip title="Discovery source">
        <Chip size="small" label={discoverySource} />
      </Tooltip>
      {networkRestricted && (
        <IconButton>
          <VpnLockIcon />
        </IconButton>
      )}
      {Object.keys(properties).length > 0 && (
        <Tooltip title={isDetailsExpanded ? "Hide details" : "Show details"}>
          <IconButton
            className={clsx(classes.expand, {
              [classes.expandOpen]: isDetailsExpanded
            })}
            onClick={onExpandDetails}
            aria-label="show more"
          >
            <ExpandMoreIcon />
          </IconButton>
        </Tooltip>
      )}
    </CardActions>
  );
};

AppCardFooter.propTypes = {
  networkRestricted: PropTypes.bool,
  discoverySource: PropTypes.string.isRequired,
  isDetailsExpanded: PropTypes.bool.isRequired,
  properties: PropTypes.object,

  onExpandDetails: PropTypes.func.isRequired
};

AppCardFooter.defaultProps = {
  networkRestricted: false,
  properties: {}
};

export default AppCardFooter;
