import secrets
import random
from locust import task, between, FastHttpUser, HttpUser
from locust.exception import StopUser

class NormalUser(HttpUser):
    wait_time = between(5, 15)

    def on_start(self):
        self.username = secrets.token_urlsafe(8)
        self.password = secrets.token_urlsafe(16)
        with self.client.post("/api/v1/register", json={"username": self.username, "password": self.password}, catch_response=True) as resp:
            if resp.status_code != 200:
                raise StopUser()
    
    def login(self):
        with self.client.post("/api/v1/login", json={"username": self.username, "password": self.password}, catch_response=True) as resp:
            if resp.status_code != 200:
                raise StopUser()
            return resp.json()["token"]
    
    @task
    def get_me(self):
        token = self.login()
        self.client.get("/me", headers={"X-Api-Token": token})
    
    @task
    def get_userlist(self):
        token = self.login()
        self.client.get("/api/v1/userlist", headers={"X-Api-Token": token})
    
    @task
    def get_scoreboard(self):
        self.client.get("/scoreboard")
    
    @task
    def get_currentmap(self):
        self.client.get("/map", params={"start": 0, "end": 10})
    
    @task(20)
    def restful_get_user_conquer_fields(self):
        token = self.login()
        self.client.get("/api/v1/fields", headers={"X-Api-Token": token})

    @task(20)
    def restful_conquer_field(self):
        token = self.login()
        with self.client.rename_request("/api/v1/conquer/[id]"):
            for _ in range(5000):
                id = random.randint(1, 1000000)
                self.client.post(f"/api/v1/conquer/{id}", headers={"X-Api-Token": token})
