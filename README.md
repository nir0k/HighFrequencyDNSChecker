# High-Frequency DNS Resolution Program for Prometheus Integration
Program Overview:
Our cutting-edge program is designed to provide high-frequency DNS name resolution for server names from a predefined list. This data is then efficiently recorded in Prometheus, empowering users with real-time insights into their server infrastructure's performance and availability.

Key Program Features:
- High-Frequency Polling: Fast and frequent DNS polling (default: every 150ms).
- Dynamic Server Naming: Random server names with timestamps.
- IP Address Verification: Ensures that IP address 1.1.1.1 is resolved for each server name.
- Continuous Monitoring: Operates endlessly in a loop for uninterrupted monitoring.
- Records resolved server names, resolve time, rcode, protocol, truncated flag and timestamps in Prometheus for historical analysis.
- Runs in an infinite loop for continuous monitoring.
- Web access to edit configuration
- Log rotate

### Install
**requred go 1.21.1**
```bash
git clone https://github.com/nir0k/HighFrequencyDNSChecker.git
cd HighFrequencyDNSChecker
make build
```

### Prepare
- Create and fill out the file `.env` in the same folder with the program. Complited example we found in the file `.env-example` in the project
- Create and fill out csv-file with DNS servers information. Complited example we found in the file `dns_servers.csv` in the project. If you using diferen filename, you are need change it in .env file


### Use
```bash
./High_Frequency_DNS_Monitoring-linux-amd64
```

### Custom Rcode:
- 30 - Resolved IP-address not equals 1.1.1.1
- 50 - DNS Server not answer


## Available parameters in configurations:

- DNS settings:
  - `DNS_RESOLVERPATH` - Path to file with list of DNS servers
  - `DNS_TIMEOUT` - DNS answer timeout in seconds
  - `DELIMETER` - Delimeter for CSV-file fields 
  - `DELIMETER_FOR_ADDITIONAL_PARAM` - Delimeter for value in fiels 'zonename', 'query_count_rps', 'zonename_with_recursion' and 'query_count_with_recursion_rps' in CSV file

- Prometheus settings:
  - `PROM_URL` - Prometheus remote write url. example: http://prometheus:8428/api/v1/write
  - `PROM_METRIC` - Prometheus metric name
  - `PROM_AUTH` - Prometheus authentication. false or true. If true, values PROM_USER and PROM_PASS are required
  - `PROM_USER` - Prometheus username
  - `PROM_PASS` - Prometheus password
  - `PROM_RETRIES` - Count retries for post data in prometheus
  - Labels: 
    - `OPCODES` - OpCodes. Possible value: true or false
    - `AUTHORITATIVE` - Authoritative. Possible value: true or false
    - `TRUNCATED` - Truncated. Possible value: true or false
    - `RCODE` - Rcode. Possible value: true or false
    - `RECURSION_DESIRED` - RecursionDesired. Possible value: true or false
    - `RECURSION_AVAILABLE` - RecursionAvailable. Possible value: true or false
    - `AUTHENTICATE_DATA` - AuthenticatedData. Possible value: true or false
    - `CHECKING_DISABLED` - CheckingDisabled. Possible value: true or false
    - `POLLING_RATE` - Polling rate. Possible value: true or false
    - `RECURSION` - Recursion. Request with reqursion or not. Possible value: true or false

- Watcher settings:
  - `CONF_CHECK_INTERVAL` - Interval check changes in config in minutes
  - `BUFFER_SIZE` - Timeseries buffer size for sent to prometheus
  - `WATCHER_LOCATION` - Watcher location
  - `WATCHER_SECURITYZONE` - Watcher security zone

- Web-Server settings:
  - `WATCHER_WEB_PORT` - Lisenning port
  - `WATCHER_WEB_USER` - Username to acces into web
  - `WATCHER_WEB_PASSWORD` - Password to acces into web

- Logging Watcher:
  - `LOG_FILE` - Path to application log file
  - `LOG_LEVEL` - Minimal severity level for application logging. Possible values: debug, info, warning, error, fatal (default: warning)
  - `LOG_MAX_AGE` - The maximum age of a log file to retain. (in days)
  - `LOG_MAX_SIZE` - The maximum size of a log file before it gets rotated. (in Megabytes (MB))
  - `LOG_MAX_FILES` - The maximum number of old log files to retain.
- Audit log Watcher:
  - `WATCHER_WEB_AUTH_LOG_FILE` - Path to auth log file
  - `WATCHER_WEB_AUTH_LOG_LEVEL` - Minimal severity level for auth logging. Possible values: debug, info, warning, error, fatal (default: warning)
  - `WATCHER_WEB_AUTH_LOG_MAX_AGE` - The maximum age of a log file to retain. (in days)
  - `WATCHER_WEB_AUTH_LOG_MAX_SIZE` - The maximum size of a log file before it gets rotated. (in Megabytes (MB))
  - `WATCHER_WEB_AUTH_LOG_MAX_FILES` - The maximum number of old log files to retain.


## CSV structure:

**Delimeter: `,` (comma)**


Example CSV:
```csv
server,server_ip,service_mode,domain,prefix,location,site,server_security_zone,protocol,zonename,query_count_rps,zonename_with_recursion,query_count_with_recursion_rps
google_dns_1,8.8.8.8,false,google.com1,dnsmon,testloc1,testsite1,srv_sec-zone1,udp,testzone1,5,testzone1_r,2
google_dns_2,8.8.4.4,false,google.com2,dnsmon,testloc2,testsite2,srv_sec-zone2,udp,testzone7&testzone10,5&4,testzone7_r&testzone100_4,2&3
```

### Field descriptiom:

 - `server` - DNS Server name. Value type: String
 - `server_ip` - DNS Server IP address. Value type: String
 - `service_mode` - Enable service mode. Possible value: true or false
 - `domain` - Domain. Value type: String
 - `prefix` - Suffix for create dunamic hostname fo resolve. Hostname create by this rule: `<unixtime with nanoseconds>.<suffix>.<zonename>`
 - `location` - DNS Server Location. Value type: String
 - `site` - DNS Server Site. Value type: String
 - `server_security_zone` - DNS Server secyrity zone. Value type: string
 - `protocol` - Protocol Used for polling. Value type: String. Possible value: tcp, udp, udp4, udp6, tcp4, tcp6
 - `zonename` - DNS Zonename without recusrion. Value type: String
 - `query_count_rps` - Count request per secconds for DSN server polling without recursion. Value type: Intenger
 - `zonename_with_recursion` - DNS Zonename with recusrion. Value type: String
 - `query_count_with_recursion_rps` - Count request per secconds for DSN server polling with recursion. Value type: Intenger
