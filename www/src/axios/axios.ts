import axios from "axios";
import router from '../router';

const instance = axios.create({
  baseURL: "http://localhost:8081",
  withCredentials: true
});

// Function to refresh the token
const refreshToken = async () => {
  const refresh_token = localStorage.getItem("refresh_token");
  if (!refresh_token) throw new Error("No refresh token available");

  try {
    const response = await axios.get("http://localhost:8081/user/refresh_token", {
      headers: { Authorization: `Bearer ${refresh_token}` },
      withCredentials: true
    });
    const newToken = response.headers["x-jwt-token"];
    const newRefreshToken = response.headers["x-refresh-token"];

    // Update tokens in localStorage
    if (newToken) {
      localStorage.setItem("token", newToken);
    }
    if (newRefreshToken) {
      localStorage.setItem("refresh_token", newRefreshToken);
    }

    return newToken;
  } catch (error) {
    console.error("Failed to refresh token", error);
    throw error;
  }
};

// Response interceptor to handle token refresh and retry
instance.interceptors.response.use(
  function (resp) {
    const newToken = resp.headers["x-jwt-token"];
    const newRefreshToken = resp.headers["x-refresh-token"];

    if (newToken) {
      localStorage.setItem("token", newToken);
    }
    if (newRefreshToken) {
      localStorage.setItem("refresh_token", newRefreshToken);
    }
    if (resp.status === 401) {
      router.push({ path: "/login" });
    }
    return resp;
  },
  async (err) => {
    const originalRequest = err.config;

    if (err.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true; // Prevent infinite retry loops

      try {
        const newToken = await refreshToken();
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return instance(originalRequest); // Retry the original request
      } catch (refreshError) {
        console.error("Token refresh failed", refreshError);
        router.push({ path: "/login" });
        return Promise.reject(refreshError);
      }
    }

    if (err.response.status === 401) {
      router.push({ path: "/login" });
    }

    return Promise.reject(err);
  }
);

// Request interceptor to attach JWT token
instance.interceptors.request.use(
  (req) => {
    const token = localStorage.getItem("token");
    if (token) {
      req.headers.Authorization = `Bearer ${token}`;
    }
    return req;
  },
  (err) => {
    console.error("Request error", err);
    return Promise.reject(err);
  }
);

export default instance;

