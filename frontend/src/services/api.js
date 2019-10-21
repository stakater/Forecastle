import axios from "axios";
import config from "../config";

// Basic configuration for HTTP calls
const instance = axios.create({
  baseURL: config.API_URL,
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
