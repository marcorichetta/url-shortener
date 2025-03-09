from locust import HttpUser, FastHttpUser, task, between
import gevent
import random
import string


class UrlShortenerUser(FastHttpUser):
    wait_time = between(1, 3)

    @task
    def shorten_urls(self):
        # Send request to shorten endpoint
        self.client.post("/shorten", data={"url": "https://example.com/"})
