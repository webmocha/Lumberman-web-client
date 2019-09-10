const API_URL = 'http://localhost:8080/api';

const apiCall = uri => fetch(API_URL + uri);

// listPrefixes /list-prefixes
const listPrefixes = () => apiCall('/list-prefixes');

// listKeys /list-keys?prefix=''
const listKeys = prefix => apiCall(`/list-keys?=${prefix}`);

// getLogsStream /get-logs-stream?prefix=''
const getLogsStream = prefix => apiCall(`/get-logs-stream?=${prefix}`);
