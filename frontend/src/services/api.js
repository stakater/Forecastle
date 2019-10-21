import axios from "axios";

// Basic configuration for HTTP calls
const instance = axios.create({
  baseURL: '/',
  responseType: "json"
});

// --- Apps --- //
export async function getApps() {
  return instance.get("api/apps");
}

// --- Config --- //
export async function getConfig() {
  return instance.get("api/config");
}

// --- Export default instance of axios API --- //
export default instance;
