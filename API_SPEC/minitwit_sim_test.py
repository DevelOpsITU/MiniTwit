import os
import json
import base64
import sqlite3
import unittest

import requests
from datetime import datetime
from contextlib import closing


BASE_URL = 'http://localhost:8080/sim'
DATABASE = "../minitwit.db"
USERNAME = 'simulator'
PWD = 'super_safe!'
CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {'Connection': 'close',
           'Content-Type': 'application/json',
           f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}


class MiniTwitTestCase(unittest.TestCase):


    ''' asd

    def init_db(self):
        """Creates the database tables."""
        with closing(sqlite3.connect(DATABASE)) as db:
            with open("../database/schema.sql") as fp:
                db.cursor().executescript(fp.read())
            db.commit()
    '''
    def setUp(self) -> None:
        timestamp = str(datetime.now())
        self.USERNAME_A = 'a_' + timestamp
        self.USERNAME_B = 'b_' + timestamp
        self.USERNAME_C = 'c_' + timestamp
        self.EMAIL_A = 'a_@' + timestamp
        self.EMAIL_B = 'b_@' + timestamp
        self.EMAIL_C = 'c_@' + timestamp


    def test_latest(self):
        # post something to update LATEST
        url = f"{BASE_URL}/register"
        data = {'username': self.USERNAME_A, 'email': self.EMAIL_A, 'pwd': 'foo'}
        params = {'latest': 1337}
        response = requests.post(url, data=json.dumps(data),
                                 params=params, headers=HEADERS)
        assert response.ok

        # verify that latest was updated
        url = f'{BASE_URL}/latest'
        response = requests.get(url, headers=HEADERS)
        assert response.ok
        assert response.json()['latest'] == 1337


    def test_register(self):
        username = self.USERNAME_A
        email = self.EMAIL_A
        pwd = 'a'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 1}
        response = requests.post(f'{BASE_URL}/register',
                                 data=json.dumps(data), headers=HEADERS, params=params)
        assert response.ok
        # TODO: add another assertion that it is really there

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 1


    def test_create_msg(self):
        self.test_register()
        username = self.USERNAME_A
        data = {'content': 'Blub!'}
        url = f'{BASE_URL}/msgs/{username}'
        params = {'latest': 2}
        response = requests.post(url, data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 2


    def test_get_latest_user_msgs(self):
        self.test_register()
        self.test_create_msg()

        username = self.USERNAME_A

        query = {'no': 20, 'latest': 3}
        url = f'{BASE_URL}/msgs/{username}'
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.status_code == 200

        got_it_earlier = False
        for msg in response.json():
            if msg['content'] == 'Blub!' and msg['user'] == username:
                got_it_earlier = True

        assert got_it_earlier

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 3


    def test_get_latest_msgs(self):
        self.test_register()
        self.test_create_msg()

        username = self.USERNAME_A
        query = {'no': 20, 'latest': 4}
        url = f'{BASE_URL}/msgs'
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.status_code == 200

        got_it_earlier = False
        for msg in response.json():
            if msg['content'] == 'Blub!' and msg['user'] == username:
                got_it_earlier = True

        assert got_it_earlier

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 4


    def test_register_b(self):
        username = self.USERNAME_B
        email = self.EMAIL_B
        pwd = 'b'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 5}
        response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok
        # TODO: add another assertion that it is really there

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 5


    def test_register_c(self):
        username = self.USERNAME_C
        email = self.EMAIL_C
        pwd = 'c'
        data = {'username': username, 'email': email, 'pwd': pwd}
        params = {'latest': 6}
        response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 6


    def test_follow_user(self):
        self.test_register()
        self.test_register_b()
        self.test_register_c()
        username = self.USERNAME_A
        url = f'{BASE_URL}/fllws/{username}'
        data = {'follow': self.USERNAME_B}
        params = {'latest': 7}
        response = requests.post(url, data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok

        data = {'follow': self.USERNAME_C}
        params = {'latest': 8}
        response = requests.post(url, data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok

        query = {'no': 20, 'latest': 9}
        response = requests.get(url, headers=HEADERS, params=query)
        assert response.ok

        json_data = response.json()
        assert self.USERNAME_B in json_data["follows"]
        assert self.USERNAME_C in json_data["follows"]

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 9


    def test_a_unfollows_b(self):
        self.test_follow_user()
        username = self.USERNAME_A
        url = f'{BASE_URL}/fllws/{username}'

        #  first send unfollow command
        data = {'unfollow': self.USERNAME_B}
        params = {'latest': 10}
        response = requests.post(url, data=json.dumps(data),
                                 headers=HEADERS, params=params)
        assert response.ok

        # then verify that b is no longer in follows list
        query = {'no': 20, 'latest': 11}
        response = requests.get(url, params=query, headers=HEADERS)
        assert response.ok
        assert self.USERNAME_B not in response.json()['follows']

        # verify that latest was updated
        response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
        assert response.json()['latest'] == 11

