import React from 'react';
import PropTypes from 'prop-types';
import { Box, Paper } from '@mui/material';
import { useTheme } from '@mui/material/styles';

import AppListItem from '../apps/AppListItem';
import AppGroup from '../apps/AppGroup';

const AppListView = ({ apps }) => {
  const theme = useTheme();
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
          <Paper
            variant="outlined"
            sx={{
              borderRadius: 2,
              overflow: 'hidden',
              borderColor: theme.palette.divider,
            }}
          >
            {apps[groupName].map((app, index) => (
              <AppListItem key={`${app.name}-${index}`} app={app} />
            ))}
          </Paper>
        </AppGroup>
      ))}
    </Box>
  );
};

AppListView.propTypes = {
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

export default AppListView;
