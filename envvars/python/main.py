import os

username = os.environ['DB_USER']
password = os.environ['DB_PASS']

def F(request):
    return f'{username}:{password}'
