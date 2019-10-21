import React from "react";
import PropTypes from "prop-types";
import { Card } from "@material-ui/core";

import AppCardHeader from "../AppCardHeader/AppCardHeader";
import AppCardFooter from "../AppCardFooter/AppCardFooter";
import AppCardContent from "../AppCardContent/AppCardContent";
import AppCardDetails from "../AppCardDetails/AppCardDetails";

const AppCard = ({ card }) => {
  const [isDetailsExpanded, SetIsDetailsExpanded] = React.useState(false);

  const handleExpandClick = () => {
    SetIsDetailsExpanded(!isDetailsExpanded);
  };

  const handleOpenAppLink = url => {
    window.open(url, "_blank");
  };

  return (
    <Card>
      <AppCardHeader
        name={card.name}
        url={card.url}
        onOpenAppLink={() => handleOpenAppLink(card.url)}
      />

      <AppCardContent
        name={card.name}
        icon={card.icon}
        onOpenAppLink={() => handleOpenAppLink(card.url)}
      />

      <AppCardFooter
        discoverySource={card.discoverySource}
        networkRestricted={card.networkRestricted}
        properties={card.properties}
        isDetailsExpanded={isDetailsExpanded}
        onExpandDetails={handleExpandClick}
      />

      <AppCardDetails
        isDetailsExpanded={isDetailsExpanded}
        properties={card.properties}
      />
    </Card>
  );
};

AppCard.propTypes = {
  card: PropTypes.shape({
    discoverySource: PropTypes.string,
    group: PropTypes.string,
    icon: PropTypes.string,
    name: PropTypes.string,
    networkRestricted: PropTypes.bool,
    properties: PropTypes.object,
    url: PropTypes.string
  }).isRequired
};

export default AppCard;
