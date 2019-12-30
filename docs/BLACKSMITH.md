# Blacksmith 

[Blacksmith](https://github.com/blacksmith-community) is an interesting community project that enables brokers for a number of components.

Configuration of Blacksmith is done by taking the base manifest and adding service brokers. The default manifest actually _fails_ to deploy
without any form of broker configured.

> Note 1: Look at the "tinsmith" releases for a shared database/cluster experience. The tinsmith service brokers need to be installed (`cf push` style).

> Note 2: The supplied `blacksmith.yml` file sets up all the forges and some random plans. Please set up accordingly. Pay attention to the various configuration options as well, in particular the VM and disk size chosen.

Deployment:

```
bosh -d blacksmith -n deploy manifests/blacksmith.yml \
     -v blacksmith_ip=192.168.123.250 \
     -v bosh_username=admin \
     -v bosh_password=${BOSH_CLIENT_SECRET} \
     -v bosh_ip=${BOSH_ENVIRONMENT}
```

To interact with the Open Service-Broker API, I used [Eden](https://github.com/starkandwayne/eden). Install via their directions.

To setup the Eden CLI:

```
$ export SB_BROKER_URL=http://192.168.123.250:3000
$ export SB_BROKER_USERNAME=blacksmith
$ export SB_BROKER_PASSWORD=$(credhub get -n /libvirt/blacksmith/broker_password --quiet)
```

With Eden, the service broker operations are easily available from the command-line.

For instance, the service catalog is avilable:

```
$ eden catalog
Service     Plan           Free         Description
=======     ====           ====         ===========
mariadb     small-4G       unspecified  no description provided
postgresql  clustered-4G   unspecified  no description provided
~           small-4G       unspecified  no description provided
rabbitmq    cluster-3      unspecified  no description provided
~           single         unspecified  no description provided
redis       cache          unspecified  no description provided
~           clustered-1x1  unspecified  no description provided
~           clustered-2x1  unspecified  no description provided
~           large          unspecified  no description provided
~           small          unspecified  no description provided
```

> Note: Postgresql-Forge _assumes_ that Cloud Foundry is deployed, which includes a Postgres BOSH release. If you do not have Postgres already loaded, pull in the latest.  Grab it [here](https://bosh.io/releases/github.com/cloudfoundry/postgres-release?all=1).

Once a Postgres release is avilable, a cluster may be provisioned:

```
$ eden -i postgres-cluster2 provision -s postgresql -p clustered-4G
provision:   postgresql/clustered-4G - name: postgres-cluster2
provision:   in-progress
provision:   in progress - 
<snip>
provision:   in progress - 
provision:   succeeded - 
provision:   done
$ eden -i postgres-cluster2 bind
Success

Run 'eden credentials -i postgres-cluster2 -b postgresql-3968ee1f-b609-4d04-892a-b7e11789c85a' to see credentials
$ eden -i postgres-cluster2 credentials
{
  "db_host": "192.168.123.21",
  "db_name": "postgres",
  "db_port": 5432,
  "host": "192.168.123.21",
  "hostname": "192.168.123.21",
  "hosts": [
    "192.168.123.23",
    "192.168.123.21",
    "192.168.123.22"
  ],
  "jdbc_read_uri": "jdbc:postgresql://192.168.123.21:5432/postgres",
  "jdbc_uri": "jdbc:postgresql://192.168.123.21:5432/postgres",
  "password": "sekr3t",
  "port": 5432,
  "read_host": "192.168.123.21",
  "read_port": 5432,
  "read_uri": "postgresql://randomuser:sekr3t@192.168.123.21:5432/postgres",
  "uri": "postgresql://randomuser:sekr3t@192.168.123.21:5432/postgres",
  "username": "randomuser"
}
```

Then you can log into the database with the credentials for the binding. Note that I connected to all three. Created a Google'd table and then inserted data in each of the 3 nodes. The data was replicated. Cool beans!

```
$ psql -h 192.168.123.23 -U randomuser -W
Password for user randomuser: 
psql (10.10 (Ubuntu 10.10-0ubuntu0.18.04.1), server 9.5.1)
Type "help" for help.

randomuser=# \dt
              List of relations
 Schema |  Name   | Type  |      Owner       
--------+---------+-------+------------------
 public | company | table | randomuser
(1 row)

randomuser=# select * from public.company;
 id | name | age | address | salary 
----+------+-----+---------+--------
(0 rows)

randomuser=# select * from public.company;
 id | name | age |          address          | salary 
----+------+-----+---------------------------+--------
  1 | Rob  |  99 | somewhere                 |   1.23
(1 row)
```
