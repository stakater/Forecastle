import React, { useEffect, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Box, Container, Skeleton } from '@mui/material';
import { useTheme } from '@mui/material/styles';

import { loadApps, refreshApps } from '../../redux/app/appsModule';
import selectApps from '../../redux/app/appsSelector';
import { selectViewMode } from '../../redux/slices/uiSlice';
import { AppGridView, AppListView } from '../../components/views';
import { EmptyState, ErrorState } from '../../components/feedback';

// Skeleton loader for cards
const CardSkeleton = () => {
  const theme = useTheme();
  return (
    <Box
      sx={{
        borderRadius: 1.5,
        border: `1px solid ${theme.palette.divider}`,
        backgroundColor: theme.palette.background.paper,
        p: 2,
        height: 200,
      }}
    >
      <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Skeleton variant="rounded" width={56} height={56} sx={{ mb: 2 }} />
        <Skeleton variant="text" width="60%" height={24} />
        <Skeleton variant="text" width="80%" height={16} sx={{ mt: 0.5 }} />
      </Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 'auto', pt: 2 }}>
        <Skeleton variant="rounded" width={60} height={22} />
        <Skeleton variant="circular" width={24} height={24} />
      </Box>
    </Box>
  );
};

// Loading grid
const LoadingGrid = () => (
  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
    <Box>
      <Skeleton variant="rounded" width={200} height={52} sx={{ mb: 2 }} />
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
            lg: 'repeat(4, 1fr)',
          },
          gap: 2,
        }}
      >
        {[1, 2, 3, 4].map((i) => (
          <CardSkeleton key={i} />
        ))}
      </Box>
    </Box>
    <Box>
      <Skeleton variant="rounded" width={180} height={52} sx={{ mb: 2 }} />
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
            lg: 'repeat(4, 1fr)',
          },
          gap: 2,
        }}
      >
        {[1, 2].map((i) => (
          <CardSkeleton key={i} />
        ))}
      </Box>
    </Box>
  </Box>
);

const AppList = () => {
  const theme = useTheme();
  const dispatch = useDispatch();

  // Redux state
  const appsData = useSelector((state) => state.apps.data);
  const filters = useSelector((state) => state.filters);
  const isLoading = useSelector((state) => state.apps.isLoading);
  const isLoaded = useSelector((state) => state.apps.isLoaded);
  const error = useSelector((state) => state.apps.error);
  const viewMode = useSelector(selectViewMode);

  // Filter apps based on search query
  const apps = selectApps(appsData, filters);
  const hasApps = Object.keys(apps).length > 0;
  const hasQuery = filters.query && filters.query.trim().length > 0;

  // Load apps on mount and auto-refresh every 30 seconds
  useEffect(() => {
    // Initial load (shows loading state)
    dispatch(loadApps());

    // Background refresh every 30 seconds (silent, only updates if data changed)
    const refreshInterval = setInterval(() => {
      dispatch(refreshApps());
    }, 30000);

    // Cleanup on unmount
    return () => clearInterval(refreshInterval);
  }, [dispatch]);

  // Retry handler
  const handleRetry = useCallback(() => {
    dispatch(loadApps());
  }, [dispatch]);

  return (
    <Box
      component="main"
      sx={{
        minHeight: 'calc(100vh - 130px)',
        backgroundColor: theme.palette.background.default,
        py: 3,
      }}
    >
      <Container maxWidth="xl">
        {/* Loading State */}
        {isLoading && !isLoaded && <LoadingGrid />}

        {/* Error State */}
        {error && !isLoaded && (
          <ErrorState
            title="Failed to load applications"
            description="We couldn't fetch the applications from the server."
            error={error}
            onRetry={handleRetry}
          />
        )}

        {/* Empty State - No Results */}
        {!isLoading && isLoaded && !hasApps && hasQuery && (
          <EmptyState
            title="No matching applications"
            description={`No applications found matching "${filters.query}". Try a different search term.`}
          />
        )}

        {/* Empty State - No Apps */}
        {!isLoading && isLoaded && !hasApps && !hasQuery && (
          <EmptyState
            title="No applications discovered"
            description="No applications have been discovered yet. Make sure your ingresses have the forecastle annotation enabled."
          />
        )}

        {/* Apps View */}
        {!isLoading && isLoaded && hasApps && (
          <>
            {viewMode === 'grid' ? (
              <AppGridView apps={apps} />
            ) : (
              <AppListView apps={apps} />
            )}
          </>
        )}
      </Container>
    </Box>
  );
};

export default AppList;
