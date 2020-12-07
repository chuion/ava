from flask import Flask
app = Flask(__name__)


@app.route('/')
def hello_world():
    print("有请求连接上来")
    return "hello_world"


if __name__ == '__main__':
    app.run(host="0.0.0.0")
