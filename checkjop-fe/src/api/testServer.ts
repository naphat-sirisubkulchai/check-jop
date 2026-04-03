import axios from "axios";

const apiClient = axios.create({
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000, // 10 sec
});

export const testService = {
  async testGet() {
    const res = await apiClient.get("https://jsonplaceholder.typicode.com/posts/1");
    return res.data;
  },
};

