import React, { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { Box } from '@mui/material';
import { useTheme } from '@mui/material/styles';

import AppList from '../AppList/AppList';
import { Header, Footer } from '../../components';
import { loadConfig } from '../../redux/app/configModule';

const App = () => {
  const theme = useTheme();
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(loadConfig());
  }, [dispatch]);

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: theme.palette.background.default,
      }}
    >
      <Header />
      <Box sx={{ flexGrow: 1 }}>
        <AppList />
      </Box>
      <Footer />
    </Box>
  );
};

export default App;
