import React, { Fragment, useState } from 'react';

import { Box, Heading, Button, Text } from 'rebass';

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

  const closeStreamState = () => {
    return stream ? stream.close() : false;
  }

  const [prefixes, setPrefixes] = useState(initialState.prefixes);
  const [selectedPrefix, setSelectedPrefix] = useState(initialState.selectedPrefix);
  const [output, setOutput] = useState(initialState.output);

  React.useEffect(() => {
    listPrefixes().then(({ prefixes }) => setPrefixes(prefixes))
  }, []);

  const selectPrefixHandler = (prefix) => {
    closeStreamState();
    setSelectedPrefix(prefix);
    setOutput(initialState.output)
  }

  const getKeysHandler = () => {
    closeStreamState();
    setOutput(initialState.output);
    return listPrefixKeys(selectedPrefix)
      .then(({ keys }) => setOutput(keys));
  }

  const getLogsHandler = () => {
    closeStreamState();
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
        <Heading color='primary' fontSize={5}>Prefixes</Heading>
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
        <Heading color='primary' fontSize={5}>Prefix: <Text as='span' color='#ed2dfd'>{selectedPrefix}</Text></Heading>
        {selectedPrefix && (
          <Box
            flex={true}
            justifyContent='end'
          >
            <Button variant='outline' onClick={getKeysHandler}>Get Keys</Button>
            <Button variant='outline' onClick={getLogsHandler}>Get Logs</Button>
            <Button variant='outline' onClick={tailLogsHandler}>Tail Logs</Button>
          </Box>
        )}
      </Box>
      <Box>
        <Heading color='primary' fontSize={5}>Output</Heading>
        <Box
          p={10}
          bg='#e6e6e6'
          flexGrow={true}
          border
          sx={{
            borderRadius: 10,
          }}
        >
          {output && output.map(item => (
            <Text key={item}><Text as='span' fontSize={12}>></Text> {item}</Text>
          ))}
        </Box>
      </Box>
    </Box>
  )
};

export default App;
