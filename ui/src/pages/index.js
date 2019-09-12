import React from 'react';
import { ThemeProvider } from 'emotion-theming';
import theme from '@rebass/preset';

import App from '../components/App';

const IndexPage = () => (
  <ThemeProvider theme={theme}>
    <App />
  </ThemeProvider>
);

export default IndexPage;
