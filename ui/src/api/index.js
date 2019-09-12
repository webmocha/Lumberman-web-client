const API_PREFIX = '/api';

const apiCall = (uri, headers = {}) => fetch(API_PREFIX + uri, {
  headers,
});

// listPrefixes /list-prefixes
export const listPrefixes = () => apiCall('/list-prefixes').then(response => response.json());

// listKeys /list-keys?prefix=''
export const listPrefixKeys = prefix => apiCall(`/list-keys?prefix=${prefix}`)
  .then(response => response.json());

// getLogsStream /get-logs-stream?prefix=''
export const getLogsStream = prefix => apiCall(`/get-logs-stream?prefix=${prefix}`)
  .then(response => response.text())
  .then(data => data.split('\n'));

// tailLogsStream /tail-logs-stream
// export const tailLogsStream = prefix => apiCall(`/get-logs-stream?prefix=${prefix}`, {
//   'Content-Type': 'application/x-ndjson',
// })
//   .then(response => response.text())
//   .then(data => data);
