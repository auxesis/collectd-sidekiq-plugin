# collectd-sidekiq-plugin

A [collectd](http://collectd.org/) [exec plugin](https://collectd.org/wiki/index.php/Plugin:Exec) to query [Sidekiq](http://sidekiq.org/) status.

## Using


```
wget github -O /usr/local/bin/check_sidekiq
./collectd_sidekiq --help
```

Then add collectd configuration:

```
LoadPlugin exec
<Plugin exec>
  Exec deploy "/usr/local/bin/check_sidekiq"
</Plugin>
```

(Change `deploy` to whatever user you want to run the check as.)

Then restart collectd:

```
sudo service collectd restart
```

You should soon see Sidekiq stats showing up in your graphs:

![image](https://cloud.githubusercontent.com/assets/12306/22453501/6a865b6a-e7d3-11e6-9220-c9240c2284ef.png)

![image](https://cloud.githubusercontent.com/assets/12306/22453520/88b5a172-e7d3-11e6-8894-95b7087532a5.png)

(the above is from Grafana + InfluxDB with a query like `SELECT mean("value") FROM "sidekiq_value" WHERE "host" = 'li123-45.members.linode.com' AND "type" = 'queue_depth' AND "instance" = 'scraper' AND $timeFilter GROUP BY time($interval) fill(null)`)

### Customising the check

To change the Redis instance to query:

```
./collectd_sidekiq --redis-server="192.168.1.111:6380"
```

To change the Redis database to query:

```
./collectd_sidekiq --redis-database=7
```

As of collectd 4.9, the exec plugin exports two environment variables:

 - `COLLECTD_HOSTNAME` - the global hostname
 - `COLLECTD_INTERVAL` - the global interval setting

To change the hostname the check reports as:

```
./collectd_sidekiq --hostname=$(hostname -f)
```

To change the interval the check runs at (this is useful for debugging):

```
./collectd_sidekiq --interval=5
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
