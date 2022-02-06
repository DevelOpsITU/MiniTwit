import os

import requests
import json
import datetime
import pytz


def log( msg : str):
  container_name = 'Thor'
  if os.getenv("HOSTNAME"):
      container_name = os.getenv("HOSTNAME")

  curr_datetime = datetime.datetime.now(pytz.timezone('Europe/Copenhagen'))
  curr_datetime = curr_datetime.isoformat('T')

  loki_hostname = 'localhost'

  if os.getenv("LOKI_HOSTNAME"):
      loki_hostname = os.getenv("LOKI_HOSTNAME")

  # push msg log into grafana-loki
  url = 'http://'+loki_hostname+':3100/api/prom/push'
  headers = {
      'Content-type': 'application/json'
  }
  payload = {
      'streams': [
          {
              'labels': '{source=\"MiniTwit\",app=\"MiniTwit\", hostname=\"' + container_name + '\"}',
              'entries': [
                  {
                      'ts': curr_datetime,
                      'line': '[INFO] ' + msg
                  }
              ]
          }
      ]
  }
  payload = json.dumps(payload)
  answer = requests.post(url, data=payload, headers=headers)
  print(answer)
  response = answer
  print(response)
  # end pushing