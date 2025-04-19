#!/usr/bin/env python3

from flask import Flask
import random
import string

app = Flask(__name__)

def generate_random_flag():
    allowed_chars = string.ascii_uppercase + string.digits
    random_part = ''.join(random.choice(allowed_chars) for _ in range(31))
    return random_part + "="

@app.route('/get-flag/')
def get_flag():
    flag = generate_random_flag()
    return flag

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081, debug=False)
