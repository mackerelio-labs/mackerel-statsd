# mackerelstatsd

## Overview

**mackerelstatsd** is a small StatsD server that calculates and posts to [Mackerel](https://mackerel.io/) the sum, or minimum/maximum/average of metrics accepted by the StatsD protocol.

StatsD is a simple UDP protocol and processor for sending metrics (see [Etsy's blog](https://www.etsy.com/codeascraft/measure-anything-measure-everything/) and [original StasD GitHub](https://github.com/statsd/statsd)). It is easy for application developers to implement metric transmission according to StasD.

## Why mackerelstatsd is needed?

Currently Mackerel receives a metric every **minute**. This means that even if the value of what you are monitoring changes significantly within a minute, the only metric stored in Mackerel is the value at the time it is retrieved.

mackerelstatsd helps solve this problem. mackerelstatsd stores the metrics that your application posts with the StatsD protocol. Their sum, or minimum/maximum/average values, are posted to Mackerel every minute. This allows you to record information that occurred in less than a minute.

mackerelstatsd supports the following two of the StatsD protocols:

- Counting: receives a counter value. This will be a metric of the sum of values over a period of time.
- Timing: receives milliseconds. These will be a metric of the minimum/maximum/average values.

## Usage

Prepare [Go's development environment](https://go.dev/dl/) and install mackerelstatsd with the following command:

```
go install github.com/mackerelio-labs/mackerel-statsd@latest
```

This will install the `mackerelstatsd` file in `$GOPATH/bin`.

To execute, do the following:
```
MACKEREL_APIKEY=<Mackerel_API_Key> $GOPATH/bin/mackerelstatsd -host <host_ID>
```

- `MACKEREL_APIKEY`: Specify the API key for your Mackerel organization. Please generate a write-enabled API key on Mackerel.
- `-host`: Specify the host ID (must be a standard host) to which the metric is posted.

### Example

There are sample implementations in `example` folder.

- `sample-client` : It posts a metric for the total number of rolls of the dice 10 times (`custom.statsd.sample.dice`) and a metric for the average/maximum/minimum latency of 10 requests to Hatena's top page (`custom.statsd.sample.http.hatena.average/max/min`). Try to execute it several times.
- `sample-http-server` : The web server listens on `localhost:8080` and posts average/maximum/minimum request processing latency metrics (`custom.statsd.sample.http.request_.average/max/min` and `custom.statsd. sample.http.request_favicon.ico.average/max/min`). Open `http://localhost:8080` in your browser and reload it several times.

![Execution example of sample-http-server](images/latency.png)

### Making it a startup service

Service files for Linux systemd are available.

1. Copy each fileto the system.

   ```
   sudo cp $GOPATH/bin/mackerelstatsd /usr/local/bin
   sudo cp example/systemd/mackerelstatsd.service /lib/systemd/system
   sudo mkdir -p /etc/sysconfig
   sudo cp example/systemd/mackerelstatsd.sysconfig /etc/sysconfig/mackerelstatsd
   ```

2. Edit `/etc/sysconfig/mackerelstatsd` file and replace the Mackerel API key (`MACKEREL_APIKEY`) and host ID (`HOSTID`) with your own.

   ```
   sudo vi /etc/sysconfig/mackerelstatsd
   ```

3. Enable systemd service.

   ```
   sudo systemctl enable mackerelstatsd
   ```

## License

Copyright 2023 Hatena Co, Ltd.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
