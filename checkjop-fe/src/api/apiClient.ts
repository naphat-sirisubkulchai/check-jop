import axios from "axios";

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

const apiClient = axios.create({
  baseURL: BASE_URL,
  timeout: 10000, // 10 sec
});

// Interceptor: Set appropriate headers based on request data
apiClient.interceptors.request.use((config) => {
  //Content-Type for FormData - let browser handle it
  if (!(config.data instanceof FormData)) {
    config.headers['Content-Type'] = 'application/json';
  }
  return config;
});


export default apiClient;
