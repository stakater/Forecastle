import React from 'react';
import PropTypes from 'prop-types';
import { Grid, Box } from '@mui/material';

import AppCard from '../apps/AppCard';
import AppGroup from '../apps/AppGroup';

const AppGridView = ({ apps }) => {
  const groupNames = Object.keys(apps).sort();

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
      {groupNames.map((groupName) => (
        <AppGroup
          key={groupName}
          name={groupName}
          count={apps[groupName].length}
          defaultExpanded
        >
          <Grid container spacing={2} alignItems="flex-start">
            {apps[groupName].map((app, index) => (
              <Grid
                item
                key={`${app.name}-${index}`}
                xs={12}
                sm={6}
                md={4}
                lg={3}
                xl={2.4}
              >
                <AppCard app={app} />
              </Grid>
            ))}
          </Grid>
        </AppGroup>
      ))}
    </Box>
  );
};

AppGridView.propTypes = {
  apps: PropTypes.objectOf(
    PropTypes.arrayOf(
      PropTypes.shape({
        name: PropTypes.string.isRequired,
        icon: PropTypes.string,
        url: PropTypes.string.isRequired,
        group: PropTypes.string,
        discoverySource: PropTypes.string,
        networkRestricted: PropTypes.bool,
        properties: PropTypes.object,
      })
    )
  ).isRequired,
};

export default AppGridView;
