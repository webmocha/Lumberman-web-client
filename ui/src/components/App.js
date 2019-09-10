import React, { Fragment, useState } from 'react';

import { Box, Heading, Button } from 'rebass';

import {
  listPrefixes,
  listPrefixKeys,
  getLogsStream,
} from '../api';

const initialState = {
  prefixes: [],
  selectedPrefix: '',
  output: [],
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

  const getKeysHandler = () => listPrefixKeys(state.selectedPrefix)
    .then(({ keys }) => setState({
      ...state,
      output: keys,
    }))

  const getLogsHandler = () => getLogsStream(state.selectedPrefix)
    .then((response) => setState({ ...state, output: response }))

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
        {state.selectedPrefix && (
          <Fragment>
            <Button color='black' onClick={getKeysHandler}>Get Keys</Button>
            <Button color='black' onClick={getLogsHandler}>Get Logs</Button>
            <Button color='black' onClick={() => {}}>Tail Logs</Button>
          </Fragment>
        )}
      </Box>
      <Box>
        <Heading fontSize={5}>Output</Heading>
        <Box pad={5}>
          {state.output && state.output.map(item => (
            <li key={item}>{item}</li>
          ))}
        </Box>
      </Box>
    </Box>
  )
};

export default App;
