# statuslight

Simple web service to set Mi-Light lamp colour depending on received status. Uses [milightd](https://github.com/sgrzywna/milightd) to control the lamp.

I use it together with CI based on Jenkins to show status of tests - green means all modules passed tests, yellow means some of the modules didn't make it, and red means catastrophic failure :-)
