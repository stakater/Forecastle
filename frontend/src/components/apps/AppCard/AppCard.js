import React, { useState, useMemo } from 'react';
import PropTypes from 'prop-types';
import { useSelector, useDispatch } from 'react-redux';
import {
  Card,
  CardContent,
  Typography,
  Box,
  IconButton,
  Tooltip,
  Collapse,
  Link,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import VpnLockIcon from '@mui/icons-material/VpnLock';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';

import AppIcon from '../AppIcon';
import AppBadge from '../AppBadge';
import { isURL } from '../../../utils/utils';
import { toggleCardExpanded, selectExpandedCards } from '../../../redux/slices/uiSlice';

const AppCard = ({ app }) => {
  const theme = useTheme();
  const dispatch = useDispatch();
  const expandedCards = useSelector(selectExpandedCards);
  const [isHovered, setIsHovered] = useState(false);

  const {
    name,
    icon,
    url,
    discoverySource,
    networkRestricted,
    properties,
  } = app;

  // Generate a unique ID for this card
  const appId = useMemo(() => `${name}-${url}`, [name, url]);
  const isExpanded = !!expandedCards[appId];

  const hasProperties = properties && Object.keys(properties).length > 0;

  const handleExpandClick = (e) => {
    e.preventDefault();
    e.stopPropagation();
    dispatch(toggleCardExpanded(appId));
  };

  return (
    <Card
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      elevation={0}
      sx={{
        display: 'flex',
        flexDirection: 'column',
        position: 'relative',
        borderRadius: 1.5,
        border: `1px solid ${isHovered ? theme.palette.primary.main : theme.palette.divider}`,
        backgroundColor: theme.palette.background.paper,
        transition: 'all 0.2s ease',
        transform: isHovered ? 'translateY(-2px)' : 'none',
        boxShadow: isHovered
          ? theme.palette.mode === 'dark'
            ? '0 8px 24px rgba(0, 0, 0, 0.4)'
            : '0 8px 24px rgba(0, 0, 0, 0.12)'
          : 'none',
        overflow: 'hidden',
        minWidth: 200,
        height: '100%',
      }}
    >
      <Box
        component="a"
        href={url}
        target="_blank"
        rel="noopener noreferrer"
        sx={{
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'stretch',
          justifyContent: 'flex-start',
          textDecoration: 'none',
          color: 'inherit',
          cursor: 'pointer',
        }}
      >
        <CardContent
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            pt: 3,
            pb: 2,
            px: 2,
            flexGrow: 1,
          }}
        >
          {/* App Icon */}
          <AppIcon src={icon} alt={name} size={56} />

          {/* App Name */}
          <Typography
            variant="subtitle1"
            component="h3"
            sx={{
              mt: 2,
              fontWeight: 600,
              textAlign: 'center',
              color: theme.palette.text.primary,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              lineHeight: 1.3,
            }}
          >
            {name}
          </Typography>

          {/* URL (truncated) */}
          <Tooltip title={url} arrow placement="bottom">
            <Typography
              variant="caption"
              sx={{
                mt: 0.5,
                color: theme.palette.text.secondary,
                textAlign: 'center',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                maxWidth: '100%',
                fontSize: '0.7rem',
              }}
            >
              {url?.replace(/^https?:\/\//, '').split('/')[0]}
            </Typography>
          </Tooltip>
        </CardContent>
      </Box>

      {/* Footer with badges and actions */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          px: 2,
          py: 1.5,
          borderTop: `1px solid ${theme.palette.divider}`,
          backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.02)'
            : 'rgba(0, 0, 0, 0.01)',
        }}
      >
        {/* Badges */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <AppBadge source={discoverySource} />
          {networkRestricted && (
            <Tooltip title="Network Restricted" arrow>
              <VpnLockIcon
                sx={{
                  fontSize: 18,
                  color: theme.palette.warning.main,
                  ml: 0.5,
                }}
              />
            </Tooltip>
          )}
        </Box>

        {/* Actions */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          {hasProperties && (
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
                '&:hover': {
                  color: theme.palette.primary.main,
                },
              }}
            >
              <OpenInNewIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Expandable Properties */}
      <Collapse in={isExpanded} timeout="auto" unmountOnExit>
        <Box
          sx={{
            borderTop: `1px solid ${theme.palette.divider}`,
            backgroundColor: theme.palette.mode === 'dark'
              ? 'rgba(255, 255, 255, 0.02)'
              : 'rgba(0, 0, 0, 0.01)',
            maxHeight: 200,
            overflow: 'auto',
          }}
        >
          <Box
            sx={{
              px: 2,
              py: 1.5,
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(80px, 1fr))',
              gap: 1.5,
            }}
          >
            {properties && Object.keys(properties).map((key) => (
              <Box key={key}>
                <Typography
                  variant="caption"
                  component="div"
                  sx={{
                    fontWeight: 600,
                    color: theme.palette.text.secondary,
                    textTransform: 'uppercase',
                    fontSize: '0.6rem',
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
                      fontSize: '0.75rem',
                      color: theme.palette.primary.main,
                      textDecoration: 'none',
                      display: 'block',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
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
                      fontSize: '0.75rem',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {properties[key]}
                  </Typography>
                )}
              </Box>
            ))}
            {(!properties || Object.keys(properties).length === 0) && (
              <Typography
                variant="caption"
                sx={{ color: theme.palette.text.secondary }}
              >
                No properties available
              </Typography>
            )}
          </Box>
        </Box>
      </Collapse>
    </Card>
  );
};

AppCard.propTypes = {
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

export default AppCard;
