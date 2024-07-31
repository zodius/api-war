// Utilities
import axios from "axios";
import { defineStore } from "pinia";

export const useAppStore = defineStore("app", {
  state: () => ({
    token: localStorage.getItem("token") || null,
    username: null,
    currentType: localStorage.getItem("currentType") || "restful",
  }),
  getters: {
    isLoggedIn() {
      return this.token !== null && this.username !== null;
    },
  },
  actions: {
    async verifyToken() {
      if (!this.token) {
        return;
      }
      try {
        let res = await axios.get("/me", {
          headers: {
            "X-API-TOKEN": this.token,
          },
        });
        let username = res.data.username;
        this.username = username;
        console.log("Logged in as", username);
      } catch (error) {
        console.error(error);
        localStorage.removeItem("token");
        this.token = null;
        this.username = null;
      }
    },
    async login(username, password) {
      switch (this.currentType) {
        case "restful":
          this.restLogin(username, password);
          break;
        case "graphql":
          this.graphqlLogin(username, password);
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
        case "graphql":
          this.graphqlRegister(username, password);
          break;
        default:
          console.error("Invalid type");
      }
    },
    async conquer(index) {
      switch (this.currentType) {
        case "restful":
          this.restConquer(index);
          break;
        case "graphql":
          this.graphqlConquer(index);
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
      this.token = res.data.token;
      this.username = username;
      console.log("Set token to", this.token);
      console.log("Set username to", this.username);
      localStorage.setItem("token", this.token);
    },
    async restRegister(username, password) {
      let res = await axios.post("/api/v1/register", {
        username: username,
        password: password,
      });
    },
    async restConquer(index) {
      let res = await axios.post(
        `/api/v1/conquer/${index}`,
        {},
        {
          headers: {
            "X-API-TOKEN": this.token,
          },
        }
      );
    },
    async graphqlRegister(username, password) {
      let res = await axios.post("/graphql", {
        query: `
          mutation {
            register(username: "${username}", password: "${password}")
          }
        `,
      });
    },
    async graphqlLogin(username, password) {
      let headers = {
        "content-type": "application/json",
      };
      let query = {
        "query": `mutation {
          login(username: "${username}", password: "${password}")
        }`,
      }
      let res = await axios({
        method: "post",
        url: "/graphql",
        data: query,
        headers: headers,
      });
      let token = res.data.data.login;
      this.token = token;
      this.username = username;
      console.log("Set token to", this.token);
      console.log("Set username to", this.username);
      localStorage.setItem("token", this.token);
    },
    async graphqlConquer(index) {
      let headers = {
        "content-type": "application/json",
        "X-Api-Token": this.token,
      };
      let query = {
        "query": `mutation {
          conquerField(FieldID: ${index})
        }`,
      }
      let res = await axios({
        method: "post",
        url: "/graphql",
        data: query,
        headers: headers,
      });
      console.log(res.data);
    }
  },
});
