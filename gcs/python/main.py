import os
import json
from google.cloud import storage

blob = storage.Client() \
    .get_bucket(os.environ['STORAGE_BUCKET']) \
    .get_blob('app1') \
    .download_as_string()

parsed = json.loads(blob)

username = parsed['username']
password = parsed['password']

def F(request):
    return f'{username}:{password}'
