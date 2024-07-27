// Utilities
import axios from "axios";
import { defineStore } from "pinia";

export const useAppStore = defineStore("app", {
  state: () => ({
    token: localStorage.getItem("token") || null,
    currentType: localStorage.getItem("currentType") || "restful",
  }),
  actions: {
    async login(username, password) {
      switch (this.currentType) {
        case "restful":
          this.restLogin(username, password);
          break;
        default:
          console.error("Invalid type");
      }
    },
    async register(username, password) {
      switch (this.currentType) {
        case "restful":
          this.restRegister(username, password);
          break;
        default:
          console.error("Invalid type");
      }
    },
    setMode(mode) {
      this.currentType = mode;
      localStorage.setItem("currentType", mode);
    },
    async restLogin(username, password) {
      let res = await axios.post("/api/v1/login", {
        username: username,
        password: password,
      });
      console.log(res);
    },
    async restRegister(username, password) {
      let res = await axios.post("/api/v1/register", {
        username: username,
        password: password,
      });
      console.log(res);
    }
  },
});
