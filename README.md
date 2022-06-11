## syslogparse
Simple syslog server, working to UDP-protocol.  
 
Server was tested with device for environmental monitoring like `NetBotz Rack Monitor 200` and etc.   
Also is will be working with anythings hosts sending syslog data on UDP-protocol on any free port (in example the `51444`).  
If data contains need the string of `High temperature`,  we call method beeper via goroutine for alarm to be beep.
#### Usage
    $ syslogpars  
#### Testing
    $ echo -e "High temperature" | nc -vutw 4 127.0.0.1 51444  
We will do it like that diapason `beep, beep, pause, pause, beep, beep`. :sos:   
And so four times like the example.
