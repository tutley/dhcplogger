## dhcplogger

Listens on a specified interface for DHCP packets, takes the Reply messages and uses them to store IP assignments into a mongo database.

#### Usage
dhcplogger -i=eth0

There is also the ability to add additional filters (the default is just port 67 and port 68). You can do this with the -filter flag

dhcplogger -i=eth0 -filter="host 192.168.0.1"



One option is to use the systemd process to make this a auto-run program, make one for each interface:

cat /etc/systemd/system/dhcplogger.service

    [Unit]
    Description=dhcplogger
    
    [Service]
    Environment="FILTER=host 192.168.0.1"
    ExecStart=/root/go/bin/dhcplogger -i=eth0 -filter=${FILTER}
    WorkingDirectory=/usr/local/bin
    Restart=always
    StandardOutput=syslog
    StandardError=syslog
    SyslogIdentifier=dhcplogger

    [Install]
    WantedBy=multi-user.target