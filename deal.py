import time
print("this python stdout!!!")


import socket
import socks
import requests

socks.set_default_proxy(socks.SOCKS5, "127.0.0.1", 4562)
socket.socket = socks.socksocket
print(requests.get('http://10.100.220.1:5000/').text)


time.sleep(2)
