# collectd-sidekiq-plugin

A [collectd](http://collectd.org/) [exec plugin](https://collectd.org/wiki/index.php/Plugin:Exec) to query [Sidekiq](http://sidekiq.org/) general and queue statistics.

## Using

Install the plugin in `/usr/local/bin` by running:

```
wget https://github.com/auxesis/collectd-sidekiq-plugin/releases/download/0.2.0/collectd_sidekiq.linux_amd64 -O /usr/local/bin/collectd_sidekiq
chmod +x /usr/local/bin/collectd_sidekiq
```

Then run the plugin to test it works:

```
/usr/local/bin/collectd_sidekiq --help
```

Then add collectd configuration:

```
LoadPlugin exec
<Plugin exec>
  Exec deploy "/usr/local/bin/collectd_sidekiq"
</Plugin>
```

(Change `deploy` to whatever user you want to run the check as)

Add these to your `types.db` (probably in `/usr/share/collectd/types.db`)

```
# Sidekiq types
processed               value:DERIVE:0:U
failed                  value:DERIVE:0:U
retries                 value:DERIVE:0:U
queue_depth             value:GAUGE:0:U
```

Then restart collectd:

```
sudo service collectd restart
```

You should soon see Sidekiq stats showing up in your graphs:

![image](https://cloud.githubusercontent.com/assets/12306/22453501/6a865b6a-e7d3-11e6-9220-c9240c2284ef.png)

![image](https://cloud.githubusercontent.com/assets/12306/22453520/88b5a172-e7d3-11e6-8894-95b7087532a5.png)

(the above is from [Grafana](http://grafana.org/) + [InfluxDB](https://www.influxdata.com/) with a query like `SELECT mean("value") FROM "sidekiq_value" WHERE "host" = 'li123-45.members.linode.com' AND "type" = 'queue_depth' AND "instance" = 'default' AND $timeFilter GROUP BY time($interval) fill(null)`)

### Customising the check

You almost certainly aren't using the default Sidekiq queues, so to change the queues to query:

```
collectd_sidekiq --queues="scraper,worker"
```

To change the Redis instance to query:

```
collectd_sidekiq --redis-server="192.168.1.111:6380"
```

To change the Redis database to query:

```
collectd_sidekiq --redis-database=7
```

As of collectd 4.9, the exec plugin exports two environment variables:

 - `COLLECTD_HOSTNAME` - the global hostname
 - `COLLECTD_INTERVAL` - the global interval setting

To change the hostname the check reports as:

```
collectd_sidekiq --hostname=$(hostname -f)
```

To change the interval the check runs at (this is useful for debugging):

```
collectd_sidekiq --interval=5
```

For these settings to take effect, make sure you update your collectd configuration appropriately:

```
LoadPlugin exec
<Plugin exec>
  Exec deploy "/usr/local/bin/collectd_sidekiq --queues='scraper,worker' --redis-server='localhost:6780' --redis-database=7"
</Plugin>
```

## Developing

Ensure you have Git, Go, and Redis installed, then run:

```
git clone https://github.com/auxesis/collectd-sidekiq-plugin.git
cd collectd-sidekiq-plugin
make # runs a local copy of the check
```


### Releasing

Follow the above developing steps, then run:

```
make build
```

This will produce a binary at ./collectd_sidekiq.linux_amd64
