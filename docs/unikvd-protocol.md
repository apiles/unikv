# UniKVd Protocol

UniKVd is a simple daemon for unikv.

Based on HTTP, the protocol is simpler than simple.

## Describe

Path to the data is described in URI in the following format:

```uri
/v1/NAMESPACE/BUCKET/KEY
```

When you need to access `count` under `web` bucket under `public` namespace, you need to:

```uri
http://server-ip:port/v1/public/web/count
```

When you GET this address, you will get the result in json.

When you POST this address, the body will be the data.

When you DELETE this address, the key will be deleted.

When KEY is set to `_list`, a list of all the keys will be returned in JSON format.

## System Namespace

When you get `/v1/_system/*`, you are visiting system namespace.

| endpoint                           | desciption                    |
| ---------------------------------- | ----------------------------- |
| `/v1/_system/lifecycle/shutdown`   | Shutdown server               |
| `/v1/_system/lifecycle/restart`    | Restart server                |
| `/v1/_system/stats/total_requests` | Get total request number      |
| `/v1/_system/stats/configure`      | Get configure structure       |
| `/v1/_system/stats/start_time`     | Get start time                |
| `/v1/_system/stats/uptime`         | Get program uptime            |
| `/v1/_system/stats/all`            | Get all statistic information |

## Running the server

Just run `unikvd SERVER_ADDR`.

Example: `unikvd 0.0.0.0:12345`.

Place the `unikv.yml` in the same directory, or specify one using the environment variable.
