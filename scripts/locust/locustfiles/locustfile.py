from locust import FastHttpUser, task, between
import logging
import time
import random
import uuid
import string

import pandas
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

def generate_random_page_visit(pdf):
    import random
    import bisect
    cdf = [sum(pdf[:i+1]) for i in range(len(pdf))]
    r = random.random() * cdf[-1]
    return bisect.bisect(cdf, r)
  
pdf = [1/i for i in range(1,11)]  

class GeneralUser(FastHttpUser):    
    weight = 1
    wait_time = between(0, 0.5)
    
    def reservation(self):
      concert_list_page_token = ""
      n_iter = generate_random_page_visit(pdf) + 1
      limit = 20
      res = None      
      for _ in range(n_iter):        
        res = self.client.get(f"/api/course/v1/courses?pageToken={concert_list_page_token}&pageSze={limit}", name="/api/course/v1/courses")
        concert_list_page_token = res.json()["nextPageToken"]
        time.sleep(0.1)
      
      items = res.json()["courses"]
      len_items = len(items)
      idx = random.randint(0, len_items - 1)
      course = items[idx]["courseId"]
      course_res = self.client.get(f"/api/course/v1/courses/{course}", name="/api/course/v1/courses/{course}") 
      
      course_data = course_res.json()
      batch = course_data["batches"][random.randint(0, len(course_data["batches"]) - 1)]["batchId"]

      booking_res = self.client.post("/api/course/v1/bookings", json={
        "course": course,
        "batch": batch,
        "customer": {
          "name": "Foo",
          "email": "foo@bar.com"  
        }
      }, name="/api/course/v1/bookings")
      
      if booking_res.status_code != 200:
        logging.error("booking: %s is not created", booking_res.json())
        return None
      
      booking = booking_res.json()["number"]
      self.client.post(f"/api/course/v1/bookings/{booking}:reserve", json={}, name="/api/course/v1/bookings/{booking}:reserve")      
      logging.info("booking: %s is reserved", booking)      
      return booking
    
    @task(3)
    def make_reservation(self):
      self.reservation()
      
    @task(1)
    def make_reservation_with_expiration(self):
      booking = self.reservation()
      if booking is None:
        return
      self.client.post(f"/api/course/v1/bookings/{booking}:expire", json={}, name="/api/course/v1/bookings/{booking}:expire")
    
    @task(1)
    def make_reservation_with_two_expiration(self):
      booking = self.reservation()      
      if booking is None:
        return
      
      self.client.post(f"/api/course/v1/bookings/{booking}:expire", json={}, name="/api/course/v1/bookings/{booking}:expire")
      self.client.post(f"/api/course/v1/bookings/{booking}:expire", json={}, name="/api/course/v1/bookings/{booking}:expire")
          
    @task(1)
    def get_course(self):
      course = str(uuid.uuid4())
      self.client.get(f"/api/course/v1/courses/{course}", name="/api/course/v1/courses/{course}") 

    def random_string(self):      
      return ''.join(random.choice(string.ascii_letters) for i in range(random.randint(10, 20)))

    @task(1)
    def get_courses_with_random_string(self):      
      course = self.random_string()
      self.client.get(f"/api/course/v1/courses/{course}", name="/api/course/v1/courses/{course}") 

class CompetingUser(FastHttpUser):    
    @task(1)
    def reservation(self):
      concert_list_page_token = ""
      limit = 1
      res = self.client.get(f"/api/course/v1/courses?pageToken={concert_list_page_token}&pageSze={limit}", name="/api/course/v1/courses")
      concert_list_page_token = res.json()["nextPageToken"]
      
      items = res.json()["courses"]      
      idx = 0
      course = items[idx]["courseId"]
      course_res = self.client.get(f"/api/course/v1/courses/{course}", name="/api/course/v1/courses/{course}") 
      
      course_data = course_res.json()
      batch = course_data["batches"][random.randint(0, len(course_data["batches"]) - 1)]["batchId"]

      booking_res = self.client.post("/api/course/v1/bookings", json={
        "course": course,
        "batch": batch,
        "customer": {
          "name": "Foo",
          "email": "foo@bar.com"  
        }
      }, name="/api/course/v1/bookings")
      
      if booking_res.status_code != 200:
        logging.error("booking: %s is not created", booking_res.json())
        return None
      
      booking = booking_res.json()["number"]
      self.client.post(f"/api/course/v1/bookings/{booking}:reserve", json={}, name="/api/course/v1/bookings/{booking}:reserve")      
      logging.info("booking: %s is reserved", booking)      
      return booking
    
    @task(1)
    def list_preload(self):
      concert_list_page_token = ""
      n_iter = generate_random_page_visit(pdf) + 1
      limit = 20
      res = None      
      for _ in range(n_iter):        
        res = self.client.get(f"/api/course/v1/courses?pageToken={concert_list_page_token}&pageSze={limit}&listMask=courses.batches", name="/api/course/v1/courses")
        concert_list_page_token = res.json()["nextPageToken"]
        time.sleep(0.1)
        
    @task(1)
    def list(self):
      concert_list_page_token = ""
      n_iter = generate_random_page_visit(pdf) + 1
      limit = 20
      res = None      
      for _ in range(n_iter):        
        res = self.client.get(f"/api/course/v1/courses?pageToken={concert_list_page_token}&pageSze={limit}", name="/api/course/v1/courses")
        concert_list_page_token = res.json()["nextPageToken"]
        time.sleep(0.1)
    