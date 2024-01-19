#!/usr/bin/env python3

import requests

# This is the heirarchy of the test data
# users:
#   - username: lcrown
#     email: lcrown@test.org
#     firstname: Lucas
#     lastname: Crownover
#   - username: marka
#     email: marka@test.org
#     firstname: Mark
#     lastname: Allen
#   - username: mollman
#     email: mollman@test.org
#     firstname: Patrick
#     lastname: Mollman
#   - username: craigs
#     email: craigs@test.org
#     firstname: Craig
#     lastname: Sorensen
# pirgs:
#   - name: racs
#     owner: marka 
#     users:
#       - lcrown
#       - marka
#     admins:
#       - lcrown
#       - marka
#   - name: systems
#     owner: craigs
#     users:
#       - craigs
#       - mollman
#     admins: 
#       - craigs

TEST_API_KEY = 'testkey1'

users = {
    "lcrown": {
        "username": "lcrown",
        "email": "lcrown@localhost",
        "firstname": "Lucas",
        "lastname": "Crownover"
    },
    "marka": {
        "username": "marka",
        "email": "marka@localhost",
        "firstname": "Mark",
        "lastname": "Allen"
    },
    "mollman": {
        "username": "mollman",
        "email": "mollman@localhost",
        "firstname": "Patrick",
        "lastname": "Mollman"
    },
    "craigs": {
        "username": "craigs",
        "email": "craigs@localhost",
        "firstname": "Craig",
        "lastname": "Sorensen"
    }
}

for user in users:
    # try to get the user
    usersearchURL = 'http://localhost:3333/api/v1/users?username={}'.format(users[user]['username'])
    r = requests.get(usersearchURL, headers={'X-API-Key': TEST_API_KEY})
    print(r.text)
    if r.status_code == 200:
        # user exists, update the id
        users[user]['id'] = r.json()['id']
        continue
    # user doesn't exist, create it
    r = requests.post('http://localhost:3333/api/v1/users', json=users[user], headers={'X-API-Key': TEST_API_KEY})
    if r.status_code == 201:
        users[user]['id'] = r.json()['id']

pirgs = [
    {
        "name": "racs",
        "owner_id": users['marka']['id'],
        "user_ids": [
            users['lcrown']['id'],
            users['marka']['id']
        ],
        "admin_ids": [
            users['marka']['id'],
        ]
    },
    {
        "name": "systems",
        "owner_id": users['craigs']['id'],
        "user_ids": [
            users['craigs']['id'],
            users['mollman']['id']
        ],
        "admin_ids": [
            users['craigs']['id'],
        ]
    }
]

for pirg in pirgs:
    r = requests.post('http://localhost:3333/api/v1/pirgs?name={}'.format(pirg['name']), json=pirg, headers={'X-API-Key': TEST_API_KEY})
    if r.status_code == 200:
        # pirg exists, update the id
        pirg['id'] = r.json()['id']
        continue
    # pirg doesn't exist, create it
    r = requests.post('http://localhost:3333/api/v1/pirgs', json=pirg, headers={'X-API-Key': TEST_API_KEY})
    if r.status_code == 201:
        pirg['id'] = r.json()['id']


print(users)
print(pirgs)
