## Configuration

### Rationale

The early version of this program didn't have configuration file and all
parameters are being hard-coded. At first I think that's enough. However, I
think giving an example of how to implement server who adheres to [12 Factor
App](https://12factor.net/) rules are important.

### Implementation

See `config.go`. This file contains code to parse configuration files from two
sources: environment variables and configuration files. Configuration files
takes precedence. These are the environment variables, configuration file key
path and default value.

| Environment Variable | YAML keypath  | Default value | Description           |
|----------------------|---------------|---------------|-----------------------|
| `KAD_LISTEN_HOST`    | `listen.host` | "127.0.0.1"   | Server Listen Address |
| `KAD_LISTEN_PORT`    | `listen.port` | 8080          | Server Port Address   |
| `KAD_DB_USER`        | `db.user`     | postgres      | Postgres User         |
| `KAD_DB_PASSWORD`    | `db.password` | password      | Postgres Password     |
| `KAD_DB_HOST`        | `db.host`     | "127.0.0.1"   | Postgres Host         |
| `KAD_DB_PORT`        | `db.port`     | 5432          | Postgres Port         |
| `KAD_DB_NAME`        | `db.db_name`  | "todo"        | Database Name         |
| `KAD_DB_SSL`         | `db.ssl_mode` | "disable"     | SSL Mode              |

The default values, if we express it in configuration file is as follows.

```yaml
listen:
  host: 127.0.0.1
  port: 8080

db:
  user: postgres
  password: password
  db_name: todo
  host: 127.0.0.1
  port: 5432 
  ssl_mode: disable
```

### Configuration file location

The program will search for `config.yaml` on current working directory, or you
can pass `-c` flag to force the program to use your own configuration file name.
For example you can run it using something like this:

```
./mda -c someconfig.yml
```

## Summary

This project is a heuristic, not a guide or a 'framework' of structure. It's to
show that you can have a sensible code structure and architecture by sticking
to the simplicity of Go.

## LICENSE

```
Copyright (c) 2023 Didiet Noor

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
