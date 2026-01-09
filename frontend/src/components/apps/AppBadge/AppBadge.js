import React from 'react';
import PropTypes from 'prop-types';
import { Chip, Tooltip } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import CloudIcon from '@mui/icons-material/Cloud';
import SettingsIcon from '@mui/icons-material/Settings';
import ExtensionIcon from '@mui/icons-material/Extension';

const AppBadge = ({ source, size = 'small' }) => {
  const theme = useTheme();

  // Configuration for different discovery sources
  const sourceConfig = {
    Ingress: {
      label: 'Ingress',
      icon: <CloudIcon sx={{ fontSize: 14 }} />,
      color: theme.palette.mode === 'dark' ? '#3b82f6' : '#2563eb',
      bgColor: theme.palette.mode === 'dark'
        ? 'rgba(59, 130, 246, 0.15)'
        : 'rgba(37, 99, 235, 0.1)',
      tooltip: 'Discovered from Kubernetes Ingress',
    },
    Config: {
      label: 'Config',
      icon: <SettingsIcon sx={{ fontSize: 14 }} />,
      color: theme.palette.mode === 'dark' ? '#8b5cf6' : '#7c3aed',
      bgColor: theme.palette.mode === 'dark'
        ? 'rgba(139, 92, 246, 0.15)'
        : 'rgba(124, 58, 237, 0.1)',
      tooltip: 'Added via configuration',
    },
    ForecastleApp: {
      label: 'CRD',
      icon: <ExtensionIcon sx={{ fontSize: 14 }} />,
      color: theme.palette.mode === 'dark' ? '#10b981' : '#059669',
      bgColor: theme.palette.mode === 'dark'
        ? 'rgba(16, 185, 129, 0.15)'
        : 'rgba(5, 150, 105, 0.1)',
      tooltip: 'Discovered from ForecastleApp CRD',
    },
  };

  const config = sourceConfig[source] || {
    label: source || 'Unknown',
    icon: null,
    color: theme.palette.text.secondary,
    bgColor: theme.palette.action.hover,
    tooltip: `Discovery source: ${source}`,
  };

  return (
    <Tooltip title={config.tooltip} arrow>
      <Chip
        size={size}
        label={config.label}
        icon={config.icon}
        sx={{
          fontWeight: 500,
          fontSize: '0.7rem',
          height: size === 'small' ? 22 : 28,
          color: config.color,
          backgroundColor: config.bgColor,
          border: 'none',
          '& .MuiChip-icon': {
            color: config.color,
            marginLeft: '6px',
          },
          '& .MuiChip-label': {
            px: 1,
          },
        }}
      />
    </Tooltip>
  );
};

AppBadge.propTypes = {
  source: PropTypes.string,
  size: PropTypes.oneOf(['small', 'medium']),
};

export default AppBadge;
