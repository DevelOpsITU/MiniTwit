# -*- coding: utf-8 -*-
"""
    MiniTwit
    ~~~~~~~~

    A microblogging application written with Flask and sqlite3.

    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
"""
import os
import re
import time
import sqlite3
from loki import *
from hashlib import md5
from datetime import datetime
from contextlib import closing
from flask import Flask, request, session, url_for, redirect, \
     render_template, abort, g, flash
from werkzeug.security import check_password_hash, generate_password_hash

import logging
import time
from jaeger_client import Config


log_level = logging.DEBUG
logging.getLogger('').handlers = []
logging.basicConfig(format='%(asctime)s %(message)s', level=log_level)

jaeger_hostname = 'localhost'

if os.getenv("JAEGER_HOSTNAME"):
    jaeger_hostname = os.getenv("JAEGER_HOSTNAME")

config = Config(
    config={ # usually read from some yaml config
        'sampler': {
            'type': 'const',
            'param': 1,
        },
        'local_agent': {
            'reporting_host': jaeger_hostname,
            'reporting_port': '6831',
        },
        'logging': True,
    },
    service_name='MiniTwit',
    validate=True,
)
# this call also sets opentracing.tracer
tracer = config.initialize_tracer()

with tracer.start_span('TestSpan') as span:
    span.log_kv({'event': 'test message', 'life': 42})

    with tracer.start_span('ChildSpan', child_of=span) as child_span:
        child_span.log_kv({'event': 'down below'})

#time.sleep(2)   # yield to IOLoop to flush the spans - https://github.com/jaegertracing/jaeger-client-python/issues/50
#tracer.close()  # flush any buffered spans

# configuration
DATABASE = '/tmp/minitwit.db'
PER_PAGE = 30
DEBUG = True
SECRET_KEY = 'development key'

# create our little application :)
app = Flask(__name__)


def connect_db():
    """Returns a new connection to the database."""
    return sqlite3.connect(DATABASE)


def init_db():
    """Creates the database tables."""
    with closing(connect_db()) as db:
        with app.open_resource('schema.sql') as f:
            db.cursor().executescript(f.read().decode("UTF-8"))
        db.commit()


def query_db(query, args=(), one=False):
    """Queries the database and returns a list of dictionaries."""
    cur = g.db.execute(query, args)
    rv = [dict((cur.description[idx][0], value)
               for idx, value in enumerate(row)) for row in cur.fetchall()]
    return (rv[0] if rv else None) if one else rv


def get_user_id(span, username):
    """Convenience method to look up the id for a username."""
    with tracer.start_span("Find user", child_of=span) as span:
        rv = g.db.execute('select user_id from user where username = ?',
                        [username]).fetchone()
        return rv[0] if rv else None


def format_datetime(timestamp):
    """Format a timestamp for display."""
    return datetime.utcfromtimestamp(timestamp).strftime('%Y-%m-%d @ %H:%M')


def gravatar_url(email, size=80):
    """Return the gravatar image for the given email address."""
    return 'http://www.gravatar.com/avatar/%s?d=identicon&s=%d' % \
        (md5(email.strip().lower().encode('utf-8')).hexdigest(), size)


@app.before_request
def before_request():
    """Make sure we are connected to the database each request and look
    up the current user so that we know he's there.
    """
    g.db = connect_db()
    g.user = None
    if 'user_id' in session:
        g.user = query_db('select * from user where user_id = ?',
                          [session['user_id']], one=True)


@app.after_request
def after_request(response):
    """Closes the database again at the end of the request."""
    g.db.close()
    return response


@app.route('/')
def timeline():
    """Shows a users timeline or if no user is logged in it will
    redirect to the public timeline.  This timeline shows the user's
    messages as well as all the messages of followed users.
    """
    with tracer.start_span('GET /') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/ traceID={:x}".format(span.trace_id))
        print("We got a visitor from: " + str(request.remote_addr))
        if not g.user:
            return redirect(url_for('public_timeline'))
        offset = request.args.get('offset', type=int)
        return render_template('timeline.html', messages=query_db('''
            select message.*, user.* from message, user
            where message.flagged = 0 and message.author_id = user.user_id and (
                user.user_id = ? or
                user.user_id in (select whom_id from follower
                                        where who_id = ?))
            order by message.pub_date desc limit ?''',
            [session['user_id'], session['user_id'], PER_PAGE]))


@app.route('/public')
def public_timeline():
    """Displays the latest messages of all users."""
    with tracer.start_span('GET /public') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/public traceID={:x}".format(span.trace_id))
        return render_template('timeline.html', messages=query_db('''
            select message.*, user.* from message, user
            where message.flagged = 0 and message.author_id = user.user_id
            order by message.pub_date desc limit ?''', [PER_PAGE]))


@app.route('/<username>')
def user_timeline(username):
    """Display's a users tweets."""
    with tracer.start_span('GET /<username>') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/" + username +" traceID={:x}".format(span.trace_id))
        profile_user = query_db('select * from user where username = ?',
                                [username], one=True)
        if profile_user is None:
            abort(404)
        followed = False
        if g.user:
            followed = query_db('''select 1 from follower where
                follower.who_id = ? and follower.whom_id = ?''',
                [session['user_id'], profile_user['user_id']], one=True) is not None
        return render_template('timeline.html', messages=query_db('''
                select message.*, user.* from message, user where
                user.user_id = message.author_id and user.user_id = ?
                order by message.pub_date desc limit ?''',
                [profile_user['user_id'], PER_PAGE]), followed=followed,
                profile_user=profile_user)


@app.route('/<username>/follow')
def follow_user(username):
    """Adds the current user as follower of the given user."""
    with tracer.start_span('GET /<username>/follow') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/" + username + "/follow traceID={:x}".format(span.trace_id))
        if not g.user:
            abort(401)
        whom_id = get_user_id(span, username)
        if whom_id is None:
            abort(404)
        g.db.execute('insert into follower (who_id, whom_id) values (?, ?)',
                    [session['user_id'], whom_id])
        g.db.commit()
        flash('You are now following "%s"' % username)
        return redirect(url_for('user_timeline', username=username))


@app.route('/<username>/unfollow')
def unfollow_user(username):
    """Removes the current user as follower of the given user."""
    with tracer.start_span('GET /<username>/unfollow') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/" + username + "/unfollow traceID={:x}".format(span.trace_id))
        if not g.user:
            abort(401)
        whom_id = get_user_id(span, username)
        if whom_id is None:
            abort(404)
        g.db.execute('delete from follower where who_id=? and whom_id=?',
                    [session['user_id'], whom_id])
        g.db.commit()
        flash('You are no longer following "%s"' % username)
        return redirect(url_for('user_timeline', username=username))


@app.route('/add_message', methods=['POST'])
def add_message():
    """Registers a new message for the user."""
    with tracer.start_span('POST /add_message') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/add_message traceID={:x}".format(span.trace_id))
        if 'user_id' not in session:
            abort(401)
        if request.form['text']:
            g.db.execute('''insert into message (author_id, text, pub_date, flagged)
                values (?, ?, ?, 0)''', (session['user_id'], request.form['text'],
                                    int(time.time())))
            g.db.commit()
            flash('Your message was recorded')
        return redirect(url_for('timeline'))


@app.route('/login', methods=['GET', 'POST'])
def login():
    """Logs the user in."""
    with tracer.start_span(request.method + ' /login') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/login traceID={:x}".format(span.trace_id))
        if g.user:
            return redirect(url_for('timeline'))
        error = None
        if request.method == 'POST':
            with tracer.start_span("Find user", child_of=span) as span:
                user = query_db('''select * from user where
                    username = ?''', [request.form['username']], one=True)
                if user is None:
                    error = 'Invalid username'
                    return render_template('login.html', error=error)
                with tracer.start_span("Verify password", child_of=span) as span:   
                    if check_password_hash(user['pw_hash'],
                                                request.form['password']):
                        error = 'Invalid password'
                        return render_template('login.html', error=error)
                    else:
                        flash('You were logged in')
                        session['user_id'] = user['user_id']
                        return redirect(url_for('timeline'))
        return render_template('login.html', error=error)


@app.route('/register', methods=['GET', 'POST'])
def register():
    """Registers the user."""
    with tracer.start_span(request.method + ' /register') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/register traceID={:x}".format(span.trace_id))
        if g.user:
            return redirect(url_for('timeline'))
        error = None
        if request.method == 'POST':
            if not request.form['username']:
                error = 'You have to enter a username'
            elif not request.form['email'] or \
                    '@' not in request.form['email']:
                error = 'You have to enter a valid email address'
            elif not request.form['password']:
                error = 'You have to enter a password'
            elif request.form['password'] != request.form['password2']:
                error = 'The two passwords do not match'
            elif get_user_id(span, request.form['username']) is not None:
                error = 'The username is already taken'
            else:
                with tracer.start_span("Create user", child_of=span) as span:
                    g.db.execute('''insert into user (
                        username, email, pw_hash) values (?, ?, ?)''',
                        [request.form['username'], request.form['email'],
                        generate_password_hash(request.form['password'])])
                    g.db.commit()
                    flash('You were successfully registered and can login now')
                    return redirect(url_for('login'))
        return render_template('register.html', error=error)


@app.route('/logout')
def logout():
    """Logs the user out"""
    with tracer.start_span('GET /logout') as span:
        log(str(request.remote_addr) + " " + str(request.method) + " path=/logout traceID={:x}".format(span.trace_id))
        flash('You were logged out')
        session.pop('user_id', None)
        return redirect(url_for('public_timeline'))


# add some filters to jinja and set the secret key and debug mode
# from the configuration.
app.jinja_env.filters['datetimeformat'] = format_datetime
app.jinja_env.filters['gravatar'] = gravatar_url
app.secret_key = SECRET_KEY
app.debug = DEBUG


if __name__ == '__main__':
    app.run(host="0.0.0.0")
