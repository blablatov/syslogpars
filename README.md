### UDP syslog-сервер  
Если данные содержат строки `High temperature` или `Maximum temperature`, вызывается метод `beep` для подачи звукового сигнала системным бипером.   
Сервер тестировался с использованием устройств для мониторинга окружающей среды, таких как `NetBotz Rack Monitor 200` и т.д.  
Также будет работать с любыми хостами, отправляющими данные системного журнала по UDP-протоколу на любой свободный порт (здесь, "51444`)  
  
#### Использование  
	$ syslogpars  
	
#### Тестирование   
	$ echo -e "High temperature" | nc -vutw 4 127.0.0.1 51444 	
В примере бипер должен звучать четыре раза, в диапазоне: "бип, бип, пауза, пауза, бип, бип". :sos:  	


####US  
Simple syslog server, working to UDP-protocol.  
If data contains need the string of `High temperature` or `Maximum temperature`, calls method `beep` to sound the system beeper.       
Server was tested with device for environmental monitoring like `NetBotz Rack Monitor 200` and etc.   
Also is will be working with anythings hosts sending syslog data on UDP-protocol on any free port (in example the `51444`).  
 
#### Usage  
    $ syslogpars  
	
#### Testing  
    $ echo -e "High temperature" | nc -vutw 4 127.0.0.1 51444  	
In the example, beeper should sound four times, in the range: `beep, beep, pause, pause, beep, beep`. :sos: 
  

