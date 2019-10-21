import React from "react";
import PropTypes from "prop-types";
import { isURL } from "../../utils/utils";
import {
  CardContent,
  Collapse,
  List,
  ListItem,
  ListItemText,
  makeStyles,
  Tooltip,
  Typography
} from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  appContentDetails: {
    backgroundColor: "#FBFBFB"
  },
  content: {
    maxHeight: "200px",
    overflow: "auto"
  }
}));

const AppCardDetails = ({ properties, isDetailsExpanded }) => {
  const classes = useStyles();

  return (
    <Collapse
      in={isDetailsExpanded}
      timeout="auto"
      unmountOnExit
      className={classes.appContentDetails}
    >
      <CardContent className={classes.content}>
        <List>
          {properties && Object.keys(properties).map(property => (
            <ListItem key={property}>
              <ListItemText
                primary={property}
                secondary={
                  isURL(properties[property]) ? (
                    <React.Fragment>
                      <Tooltip title={properties[property]}>
                        <Typography
                          variant="subtitle2"
                          noWrap
                          color="primary"
                          component="a"
                          href={properties[property]}
                          target="_blank"
                          style={{
                            textDecoration: "none",
                            width: "100%",
                            display: "inline-block"
                          }}
                        >
                          {properties[property]}
                        </Typography>
                      </Tooltip>
                    </React.Fragment>
                  ) : (
                    properties[property]
                  )
                }
              />
            </ListItem>
          ))}

          {/* If no properties exist */}
          {Object.keys(properties).length === 0 && (
            <ListItem>
              <ListItemText secondary="No properties found." />
            </ListItem>
          )}
        </List>
      </CardContent>
    </Collapse>
  );
};

AppCardDetails.propTypes = {
  properties: PropTypes.object,
  isDetailsExpanded: PropTypes.bool.isRequired
};

AppCardDetails.defaultProps = {
  properties: {}
};

export default AppCardDetails;
