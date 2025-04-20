# Client

The client is the part of CookieFarm that is responsible for retrieving flags from the exploit manager and communicating with the server to give it the flags it has retrieved (which will later be sent to the flag checker).

## Exploit manager

The exploit manager is a python script but used as a library that provides a decorator to be used on top of a function so that:

- The exploit is run on all adversary machines.
- Every how many seconds to run an Exploit
- The handling of multithreading

### How it works

```python
#!/usr/bin/env python3

from utils.exploiter_manager import exploit_manager

import sys

@exploit_manager(server_ip=sys.argv[1])
def exploit(ip:str = "", port:int = 80):
    # Insert your exploit code here
    return  # Return flag

if __name__ == "__main__":
    port = 8081                 # Edit with port of the service to exploit
    exploit(port=port)

```

This approach allows the user to focus only on writing the 'exploit without worrying about multithreading or launching the exploit on all the adversary machines, making the exploitation process easier and more efficient. in addition, stdout handling with flush or script configuration is done automatically so as to reduce errors and avoid wasting time.
