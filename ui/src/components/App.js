import React, { useState, useEffect } from 'react';

import { Box, Heading } from 'rebass';

import { listPrefixes, listPrefixKeys } from '../api';

const initialState = {
  prefixes: [],
  selectedPrefix: '',
}

const App = () => {
  const [state, setState] = useState(initialState);

  React.useEffect(() => {
    listPrefixes().then(({ prefixes }) => setState({
      prefixes,
    }))
  }, []);

  const selectPrefixHandler = (prefix) => () => setState({
    ...state,
    selectedPrefix: prefix,
  });

  return (
    <Box>
      <Box>
        <Heading fontSize={5}>Prefixes</Heading>
        <ul>
          {state.prefixes.map(prefix => (
            <li
              key={prefix}
              onClick={selectPrefixHandler(prefix)}
            >
              {prefix}
            </li>
          ))}
        </ul>
      </Box>
      <Box>
        <Heading fontSize={5}>Prefix: {state.selectedPrefix}</Heading>
      </Box>
    </Box>
  )
};

export default App;
