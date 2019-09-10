import React, { useState, useEffect } from 'react';

import { Box, Heading } from 'rebass';

import { listPrefixes } from '../api';

const initialState = {
  prefixes: [],
  key: '',
}

const App = () => {
  const [state, setState] = useState(initialState);

  React.useEffect(() => {
    listPrefixes().then(({ prefixes }) => setState({
      prefixes,
    }))
  }, []);

  return (
    <Box>
      <Box>
        <Heading fontSize={5}>Prefixes</Heading>
        <ul>
          {state.prefixes.map(prefix => <li key={prefix}>{prefix}</li>)}
        </ul>
      </Box>
    </Box>
  )
};

export default App;
