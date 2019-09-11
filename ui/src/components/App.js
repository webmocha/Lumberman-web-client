import React, { Fragment, useState } from 'react';

import { Box, Heading, Button } from 'rebass';

import {
  listPrefixes,
  listPrefixKeys,
  getLogsStream,
} from '../api';

const App = () => {
  let stream;
  const [prefixes, setPrefixes] = useState([]);
  const [selectedPrefix, setSelectedPrefix] = useState('');
  const [output, setOutput] = useState([]);

  React.useEffect(() => {
    listPrefixes().then(({ prefixes }) => setPrefixes(prefixes))
  }, []);

  const selectPrefixHandler = (prefix) => () => setSelectedPrefix(prefix);

  const getKeysHandler = () => {
    console.log(output)
    listPrefixKeys(selectedPrefix)
      .then(({ keys }) => setOutput(keys));
  }

  const getLogsHandler = () => {
    console.log(output)
    getLogsStream(selectedPrefix)
      .then(setOutput)
  };

  const tailLogsHandler = () => {
    if (!!window.EventSource) {
      stream = new EventSource(`/api/tail-logs-stream?prefix=${selectedPrefix}`);
    } else {
      console.warn('Streaming is not available on your client');
      return null;
    }

    stream.addEventListener('log', (e) => {
      setOutput(s => [...s, e.data]);
    })

    stream.addEventListener('error', (e) => {
      console.log(e.readyState);
    }, false)
  };

  return (
    <Box>
      <Box>
        <Heading fontSize={5}>Prefixes</Heading>
        <ul>
          {prefixes.map(prefix => (
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
