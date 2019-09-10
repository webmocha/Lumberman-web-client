const API_PREFIX = '/api';

const apiCall = uri => fetch(API_PREFIX + uri).then(response => response.json());

// listPrefixes /list-prefixes
export const listPrefixes = () => apiCall('/list-prefixes');

// listKeys /list-keys?prefix=''
export const listPrefixKeys = prefix => apiCall(`/list-keys?=${prefix}`);

// getLogsStream /get-logs-stream?prefix=''
export const getLogsStream = prefix => apiCall(`/get-logs-stream?=${prefix}`);
