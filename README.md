# vultrdata

Service to return sanitized data about the requesting instance.

## Installation

For systemd Linux systems:

- Run `make deps build-linux`
- Create the directory /opt/vultrdata
- Copy `vultrdata-linux` to /opt/vultrdata/vultrdata
  (or copy the file named `vultrdata` if make run on a Linux system)
- Copy the file vultrdata.service to /etc/systemd/vultrdata.service
- Edit that vultrdata.service file changing the API_KEY value as well as listen addr/port and/or removing --userdata option
- On that system: `sudo systemctl daemon-reload`
- On that system: `sudo systemctl enable vultrdata`
- On that system: `sudo systemctl start vultrdata`

It is recommended that the listen address is chosen to be on the internal network for security. When instances are created with internal networking enabled it usually has an additional address starting with 10. or possibly some other non-routing range.

## Testing

- `curl -si 'http://10.1.1.1:8888/'`
  (Using whichever listen addr/port was chosen in vultrdata.service)
- Run the same command from another instance which should return data for that instance different than the first
