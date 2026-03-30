# System information agent
Provides data about CPU, RAM, Disks and Network devices and makes it accessible via API endpoint

## Data
Data model is described in "models.go" file. If error occurs, CPU and Mem fileds are zeroed out, while Disk and Net become NULL. Agent is designed to stay online even if no data is gathered for some reason.

## Access
By default information is accessible by "localhost:8080/sysinfo" route. Port can be changed by specifing PORT env variable.

## Tools
Agent uses Gin for handling API requests and gopsutil(v4) for gathering information.
