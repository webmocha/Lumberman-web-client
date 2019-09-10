const API_PREFIX = '/api';

const apiCall = (uri, headers = {}) => fetch(API_PREFIX + uri, {
  headers,
}).then(response => response.json());

// listPrefixes /list-prefixes
export const listPrefixes = () => apiCall('/list-prefixes');

// listKeys /list-keys?prefix=''
export const listPrefixKeys = prefix => apiCall(`/list-keys?prefix=${prefix}`);

// getLogsStream /get-logs-stream?prefix=''
export const getLogsStream = prefix => apiCall(`/get-logs-stream?prefix=${prefix}`);

// tailLogsStream /tail-logs-stream
export const tailLogsStream = prefix => apiCall(`/get-logs-stream?prefix=${prefix}`, {
  'Content-Type': 'application/x-ndjson',
});
