# vultrdata

Service to return santized metadata about the requesting instance itself.

## Building

- Run `make deps build-linux`

The v1.0.0 release has a Linux 64-bit binary compiled with Go 1.11.2

## Installation

For systemd Linux systems:

- Create the directory /opt/vultrdata
- Copy `vultrdata-linux` to /opt/vultrdata/vultrdata
  (or copy the file named `vultrdata` if make run on a Linux system)
- Copy the file vultrdata.service to /etc/systemd/vultrdata.service
- Edit that vultrdata.service file changing the API_KEY value as well as listen addr/port and/or removing --userdata option
- On that system: `sudo systemctl daemon-reload`
- On that system: `sudo systemctl enable vultrdata`
- On that system: `sudo systemctl start vultrdata`

It is recommended that the listen address is chosen to be on the internal network for security. When instances are created with internal networking enabled it usually has an additional address starting with 10. or possibly some other non-routing range.

You may have to configure both the Vultr Firewall configuration as well as any firewall configured on the operating-system.

## Testing

- `curl -si 'http://10.1.1.1:8888/'`
  (Using whichever listen addr/port was chosen in vultrdata.service)
- Run the same command from another instance which should return data for that instance different than the first
- Log messages from the service go to `/var/log/syslog`
