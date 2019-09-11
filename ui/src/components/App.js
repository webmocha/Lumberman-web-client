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
};

const App = () => {
  let stream;

  const checkStreamState = () => {
    return stream && stream.close();
  }

  const [prefixes, setPrefixes] = useState(initialState.prefixes);
  const [selectedPrefix, setSelectedPrefix] = useState(initialState.selectedPrefix);
  const [output, setOutput] = useState(initialState.output);

  React.useEffect(() => {
    listPrefixes().then(({ prefixes }) => setPrefixes(prefixes))
  }, []);

  const selectPrefixHandler = (prefix) => {
    checkStreamState();
    setSelectedPrefix(prefix);
  }

  const getKeysHandler = () => {
    checkStreamState();
    setOutput(initialState.output);
    return listPrefixKeys(selectedPrefix)
      .then(({ keys }) => setOutput(keys));
  }

  const getLogsHandler = () => {
    checkStreamState();
    setOutput(initialState.output);
    return getLogsStream(selectedPrefix)
      .then(setOutput)
  };

  const tailLogsHandler = () => {
    setOutput(initialState.output);
    if (!!window.EventSource) {
      stream = new EventSource(`/api/tail-logs-stream?prefix=${selectedPrefix}`);
    } else {
      console.warn('Streaming is not available on your client');
      return null;
    }

    stream.addEventListener('log', (e) => {
      setOutput(s => [...s, e.data]);
    }, false);

    stream.addEventListener('error', (e) => {
      console.log(e.readyState);
    }, false);
  };

  return (
    <Box>
      <Box>
        <Heading fontSize={5}>Prefixes</Heading>
        <ul>
          {prefixes.map(prefix => (
            <li
              key={prefix}
              onClick={() =>selectPrefixHandler(prefix)}
            >
              {prefix}
            </li>
          ))}
        </ul>
      </Box>
      <Box>
        <Heading fontSize={5}>Prefix: {selectedPrefix}</Heading>
        {selectedPrefix && (
          <Fragment>
            <Button color='black' onClick={getKeysHandler}>Get Keys</Button>
            <Button color='black' onClick={getLogsHandler}>Get Logs</Button>
            <Button color='black' onClick={tailLogsHandler}>Tail Logs</Button>
          </Fragment>
        )}
      </Box>
      <Box>
        <Heading fontSize={5}>Output</Heading>
        <Box pad={5}>
          {output && output.map(item => (
            <li key={item}>{item}</li>
          ))}
        </Box>
      </Box>
    </Box>
  )
};

export default App;
