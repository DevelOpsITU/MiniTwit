# -*- coding: utf-8 -*-
"""
    MiniTwit Tests
    ~~~~~~~~~~~~~~

    Tests the MiniTwit application.

    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
"""
import unittest
import tempfile
from datetime import datetime
import urllib.parse

import requests



class MiniTwitTestCase(unittest.TestCase):

    def setUp(self):
        """Before each test"""
        self.hostname = "localhost"
        self.port = 5000
        timestamp = str(datetime.now())
        self.test_username0 = timestamp+"_minitwit_client_test0"
        self.test_username1 = timestamp+"_minitwit_client_test1"
        self.test_email0 = self.test_username0 + '@example.com'
        self.test_email1 = self.test_username1 + '@example.com'

    def tearDown(self):
        """After each test"""
        print("Niceness!")


    # helper functions

    def register(self, username:str, password:str, password2=None, email=None):
        """Helper function to register a user"""
        if password2 is None:
            password2 = password
        if email is None:
            email = username + '@example.com'
        url = "http://"+self.hostname+":"+str(self.port)+"/register"
        
        params = {'username': username, 'email': email, 'password': password, "password2": password2 }
        
        payload=urllib.parse.urlencode(params)
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
        
        response = requests.request("POST", url, headers=headers, data=payload)
            
        return response

    def login(self, username, password):
        """Helper function to login"""
        url = "http://"+self.hostname+":"+str(self.port)+"/login"

        params = {'username': username, 'password': password}

        payload=urllib.parse.urlencode(params)
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }

        return requests.request("POST", url, headers=headers, data=payload)

    def register_and_login(self, username, password):
        """Registers and logs in in one go"""
        self.register(username, password)
        return self.login(username, password)

    def logout(self, cookieJar):
        """Helper function to logout"""
        url = "http://"+self.hostname+":"+str(self.port)+"/logout"
        return requests.request("GET", url, cookies=cookieJar)


    def add_message(self, text, cookieJar):
        """Records a message"""
        url = "http://"+self.hostname+":"+str(self.port)+"/add_message"
        params = {'text': text }

        payload=urllib.parse.urlencode(params)
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }

        rv = requests.request("POST", url, headers=headers, data=payload, cookies=cookieJar)

        assert rv.status_code == 200
        if text :
            assert 'Your message was recorded' in rv.text
        return rv

    def get_public_timeline(self):
        url = "http://"+self.hostname+":"+str(self.port)+"/public"
        return requests.request("GET", url)

    def get_private_timeline(self,cookieJar):
        url = "http://"+self.hostname+":"+str(self.port)+"/"
        return requests.request("GET", url, cookies=cookieJar)


    def get_follow_user(self, username, cookieJar):
        username = username.replace(" ","%20")

        url = "http://"+self.hostname+":"+str(self.port)+"/"+username+'/follow'
        return requests.request("GET", url, cookies=cookieJar)

    def get_personal_timeline(self, username):
        username = username.replace(" ","%20")

        url = "http://"+self.hostname+":"+str(self.port)+"/"+username
        return requests.request("GET", url)

    def unfollow_user(self, username, cookieJar):
        username = username.replace(" ","%20")

        url = "http://"+self.hostname+":"+str(self.port)+"/"+username+'/unfollow'
        return requests.request("GET", url, cookies=cookieJar)

    # testing functions

    def test_register(self):
        """Make sure registering works"""
        rv = self.register(self.test_username0, "default")
        assert 'You were successfully registered ' \
               'and can login now' in rv.text
        rv = self.register(self.test_username0, 'default')
        assert 'The username is already taken' in rv.text
        rv = self.register('', 'default')
        assert 'You have to enter a username' in rv.text
        rv = self.register('meh', '')
        assert 'You have to enter a password' in rv.text
        rv = self.register('meh', 'x', 'y')
        assert 'The two passwords do not match' in rv.text
        rv = self.register('meh', 'foo', email='broken')
        assert 'You have to enter a valid email address' in rv.text

    def test_login_logout(self):
        """Make sure logging in and logging out works"""
        rv = self.register_and_login(self.test_username0, 'default')
        cookie = rv.cookies
        assert 'You were logged in' in rv.text
        rv = self.logout(cookie)
        assert 'You were logged out' in rv.text
        rv = self.login(self.test_username0, 'wrongpassword')
        assert 'Invalid password' in rv.text
        rv = self.login(self.test_username1, 'wrongpassword')
        assert 'Invalid username' in rv.text

    def test_message_recording(self):
        """Check if adding messages works"""
        rv = self.register_and_login(self.test_username0, 'default')
        cookie = rv.cookies
        self.add_message('test message 1',cookie)
        self.add_message('<test message 2>',cookie)
        rv = self.get_private_timeline(cookie)
        assert 'test message 1' in rv.text
        assert '&lt;test message 2&gt;' in rv.text

    def test_timelines(self):
        """Make sure that timelines work"""
        rv = self.register_and_login(self.test_username0, 'default')
        cookie = rv.cookies

        self.add_message('the message by foo',cookie)
        self.logout(cookie)
        rv = self.register_and_login(self.test_username1, 'default')
        cookie = rv.cookies
        self.add_message('the message by bar',cookie)
        rv = self.get_public_timeline()
        assert 'the message by foo' in rv.text
        assert 'the message by bar' in rv.text

        # bar's timeline should just show bar's message
        rv = self.get_private_timeline(cookie)
        assert 'the message by foo' not in rv.text
        assert 'the message by bar' in rv.text

        # now let's follow foo
        rv = self.get_follow_user(self.test_username0,cookie)
        assert 'You are now following &#34;'+self.test_username0+'&#34;' in rv.text

        # we should now see foo's message
        rv = self.get_private_timeline(cookie)
        assert 'the message by foo' in rv.text
        assert 'the message by bar' in rv.text

        # but on the user's page we only want the user's message
        rv = self.get_personal_timeline(self.test_username1)
        assert 'the message by foo' not in rv.text
        assert 'the message by bar' in rv.text
        rv = self.get_personal_timeline(self.test_username0)
        assert 'the message by foo' in rv.text
        assert 'the message by bar' not in rv.text

        # now unfollow and check if that worked
        rv = self.unfollow_user(self.test_username0,cookie)

        assert 'You are no longer following &#34;'+self.test_username0+'&#34;' in rv.text
        rv = self.get_private_timeline(cookie)
        assert 'the message by foo' not in rv.text
        assert 'the message by bar' in rv.text




if __name__ == '__main__':
    unittest.main()
