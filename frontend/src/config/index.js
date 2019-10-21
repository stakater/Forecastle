const dev = {
  API_URL: "http://localhost:3000/"
};

const prod = {
  API_URL: "http://localhost:3000/"
};

const config = process.env.REACT_APP_ENV === "production" ? prod : dev;

export default {
  ...config
};
