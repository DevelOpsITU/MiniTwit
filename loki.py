import requests
import json
import datetime
import pytz


def log( msg : str):
  host = 'Thor'
  curr_datetime = datetime.datetime.now(pytz.timezone('Europe/Copenhagen'))
  curr_datetime = curr_datetime.isoformat('T')

  # push msg log into grafana-loki
  url = 'http://loki:3100/api/prom/push'
  headers = {
      'Content-type': 'application/json'
  }
  payload = {
      'streams': [
          {
              'labels': '{source=\"MiniTwit\",app=\"MiniTwit\", host=\"' + host + '\"}',
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