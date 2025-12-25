import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
  Box,
  Typography,
  IconButton,
  Tooltip,
  Collapse,
  Link,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import VpnLockIcon from '@mui/icons-material/VpnLock';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

import AppIcon from '../AppIcon';
import AppBadge from '../AppBadge';
import { isURL } from '../../../utils/utils';

const AppListItem = ({ app }) => {
  const theme = useTheme();
  const [isHovered, setIsHovered] = useState(false);
  const [isExpanded, setIsExpanded] = useState(false);

  const {
    name,
    icon,
    url,
    discoverySource,
    networkRestricted,
    properties,
  } = app;

  const hasProperties = properties && Object.keys(properties).length > 0;

  const handleExpandClick = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsExpanded(!isExpanded);
  };

  const handleRowClick = () => {
    window.open(url, '_blank', 'noopener,noreferrer');
  };

  return (
    <Box
      sx={{
        borderBottom: `1px solid ${theme.palette.divider}`,
        '&:last-child': {
          borderBottom: 'none',
        },
      }}
    >
      {/* Main Row */}
      <Box
        onClick={handleRowClick}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 2,
          px: 2,
          py: 1.5,
          cursor: 'pointer',
          backgroundColor: isHovered
            ? theme.palette.action.hover
            : 'transparent',
          transition: 'background-color 0.15s ease',
          '&:active': {
            backgroundColor: theme.palette.action.selected,
          },
        }}
      >
        {/* Icon */}
        <AppIcon src={icon} alt={name} size={40} />

        {/* Name and URL */}
        <Box sx={{ flexGrow: 1, minWidth: 0 }}>
          <Typography
            variant="subtitle2"
            sx={{
              fontWeight: 600,
              color: theme.palette.text.primary,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            }}
          >
            {name}
          </Typography>
          <Tooltip title={url} arrow placement="bottom-start">
            <Typography
              variant="caption"
              sx={{
                color: theme.palette.text.secondary,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                display: 'block',
                fontSize: '0.7rem',
              }}
            >
              {url?.replace(/^https?:\/\//, '').split('/')[0]}
            </Typography>
          </Tooltip>
        </Box>

        {/* Badges */}
        <Box
          sx={{
            display: { xs: 'none', sm: 'flex' },
            alignItems: 'center',
            gap: 1,
            flexShrink: 0,
          }}
        >
          <AppBadge source={discoverySource} />
          {networkRestricted && (
            <Tooltip title="Network Restricted" arrow>
              <VpnLockIcon
                sx={{
                  fontSize: 18,
                  color: theme.palette.warning.main,
                }}
              />
            </Tooltip>
          )}
        </Box>

        {/* Actions */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 0.5,
            flexShrink: 0,
          }}
        >
          {hasProperties ? (
            <Tooltip title={isExpanded ? 'Hide details' : 'Show details'} arrow>
              <IconButton
                size="small"
                onClick={handleExpandClick}
                sx={{
                  color: theme.palette.text.secondary,
                  transform: isExpanded ? 'rotate(180deg)' : 'none',
                  transition: 'transform 0.2s ease',
                  '&:hover': {
                    color: theme.palette.text.primary,
                  },
                }}
              >
                <ExpandMoreIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          ) : (
            // Placeholder to maintain alignment when no properties
            <Box sx={{ width: 28, height: 28 }} />
          )}
          <Tooltip title="Open in new tab" arrow>
            <IconButton
              size="small"
              component="a"
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              onClick={(e) => e.stopPropagation()}
              sx={{
                color: theme.palette.text.secondary,
                opacity: isHovered ? 1 : 0.5,
                transition: 'opacity 0.15s ease',
                '&:hover': {
                  color: theme.palette.primary.main,
                  opacity: 1,
                },
              }}
            >
              <OpenInNewIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <ChevronRightIcon
            sx={{
              fontSize: 20,
              color: theme.palette.text.secondary,
              opacity: isHovered ? 1 : 0.3,
              transition: 'opacity 0.15s ease',
            }}
          />
        </Box>
      </Box>

      {/* Expandable Properties */}
      <Collapse in={isExpanded} timeout="auto" unmountOnExit>
        <Box
          sx={{
            pl: 9,
            pr: 3,
            pt: 1,
            pb: 2,
            display: 'flex',
            flexWrap: 'wrap',
            gap: 3,
            backgroundColor: theme.palette.mode === 'dark'
              ? 'rgba(255, 255, 255, 0.02)'
              : 'rgba(0, 0, 0, 0.01)',
          }}
        >
          {properties && Object.keys(properties).map((key) => (
            <Box key={key} sx={{ minWidth: 100 }}>
              <Typography
                variant="caption"
                component="div"
                sx={{
                  fontWeight: 600,
                  color: theme.palette.text.secondary,
                  textTransform: 'uppercase',
                  fontSize: '0.65rem',
                  letterSpacing: '0.05em',
                  mb: 0.25,
                }}
              >
                {key}
              </Typography>
              {isURL(properties[key]) ? (
                <Link
                  href={properties[key]}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(e) => e.stopPropagation()}
                  sx={{
                    fontSize: '0.8rem',
                    color: theme.palette.primary.main,
                    textDecoration: 'none',
                    '&:hover': {
                      textDecoration: 'underline',
                    },
                  }}
                >
                  {properties[key]}
                </Link>
              ) : (
                <Typography
                  variant="body2"
                  sx={{
                    color: theme.palette.text.primary,
                    fontSize: '0.8rem',
                  }}
                >
                  {properties[key]}
                </Typography>
              )}
            </Box>
          ))}
        </Box>
      </Collapse>
    </Box>
  );
};

AppListItem.propTypes = {
  app: PropTypes.shape({
    name: PropTypes.string.isRequired,
    icon: PropTypes.string,
    url: PropTypes.string.isRequired,
    group: PropTypes.string,
    discoverySource: PropTypes.string,
    networkRestricted: PropTypes.bool,
    properties: PropTypes.object,
  }).isRequired,
};

export default AppListItem;
