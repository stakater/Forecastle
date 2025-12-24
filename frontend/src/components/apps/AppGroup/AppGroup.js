import React from 'react';
import PropTypes from 'prop-types';
import {
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Typography,
  Box,
  Chip,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FolderIcon from '@mui/icons-material/Folder';

const AppGroup = ({ name, count, children, defaultExpanded = true }) => {
  const theme = useTheme();

  // Format group name for display
  const displayName = name
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');

  return (
    <Accordion
      defaultExpanded={defaultExpanded}
      disableGutters
      sx={{
        backgroundColor: 'transparent',
        boxShadow: 'none',
        borderRadius: '12px !important',
        overflow: 'hidden',
        '&:before': {
          display: 'none',
        },
        '&.Mui-expanded': {
          margin: 0,
        },
      }}
    >
      <AccordionSummary
        expandIcon={
          <ExpandMoreIcon
            sx={{
              color: theme.palette.text.secondary,
            }}
          />
        }
        sx={{
          backgroundColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.05)'
            : 'rgba(0, 0, 0, 0.03)',
          borderRadius: '12px',
          minHeight: 52,
          px: 2,
          '&:hover': {
            backgroundColor: theme.palette.mode === 'dark'
              ? 'rgba(255, 255, 255, 0.08)'
              : 'rgba(0, 0, 0, 0.05)',
          },
          '&.Mui-expanded': {
            minHeight: 52,
            borderBottomLeftRadius: 0,
            borderBottomRightRadius: 0,
          },
          '& .MuiAccordionSummary-content': {
            margin: '12px 0',
            '&.Mui-expanded': {
              margin: '12px 0',
            },
          },
        }}
      >
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 1.5,
            width: '100%',
          }}
        >
          <FolderIcon
            sx={{
              fontSize: 20,
              color: theme.palette.primary.main,
              opacity: 0.8,
            }}
          />
          <Typography
            variant="subtitle2"
            sx={{
              fontWeight: 600,
              color: theme.palette.text.primary,
              textTransform: 'uppercase',
              letterSpacing: '0.05em',
              fontSize: '0.8rem',
            }}
          >
            {displayName}
          </Typography>
          <Chip
            label={count}
            size="small"
            sx={{
              height: 20,
              minWidth: 28,
              fontSize: '0.7rem',
              fontWeight: 600,
              backgroundColor: theme.palette.mode === 'dark'
                ? 'rgba(255, 255, 255, 0.1)'
                : 'rgba(0, 0, 0, 0.06)',
              color: theme.palette.text.secondary,
            }}
          />
        </Box>
      </AccordionSummary>

      <AccordionDetails
        sx={{
          p: 0,
          pt: 2,
          pb: 3,
        }}
      >
        {children}
      </AccordionDetails>
    </Accordion>
  );
};

AppGroup.propTypes = {
  name: PropTypes.string.isRequired,
  count: PropTypes.number.isRequired,
  children: PropTypes.node.isRequired,
  defaultExpanded: PropTypes.bool,
};

export default AppGroup;
