# Cluster Balance Report

## Process

1. Extract the topic details

```bash
rpk topic describe -p -r '*' > input.txt
```

2. Process the details into CSV
```bash
go run . -input input.txt > data.csv 
```

3. Load the data into sqllite

```bash
$ sqlite3

sqlite> .mode csv

sqlite> create table replicas (
    "topic" text,
    "partition" integer,
    "leader" boolean,
    "replica" integer,
    "start" long,
    "end" long
);

sqlite> .import data.csv replicas
```

# Queries

First, generate a unique list of node IDs:

```sql
select distinct replica from replicas order by replica;
```
## Leadership Distribution per Topic

Modify the query below to include the correct pivot columns (one column per node id):

```sql
-- Show leader balance by topic
select topic,
       sum(replicas_on_2), -- modify this list to match your node ids
       sum(replicas_on_3),
       sum(replicas_on_5),
       sum(replicas_on_6),
       sum(replicas_on_8),
       sum(replicas_on_9),
       sum(replicas_on_10),
       sum(replicas_on_12),
       sum(replicas_on_13)
from (select topic,
             case when replica = 2 then 1 else 0 end  replicas_on_2, -- modify this list to match your node ids
             case when replica = 3 then 1 else 0 end  replicas_on_3,
             case when replica = 5 then 1 else 0 end  replicas_on_5,
             case when replica = 6 then 1 else 0 end  replicas_on_6,
             case when replica = 8 then 1 else 0 end  replicas_on_8,
             case when replica = 9 then 1 else 0 end  replicas_on_9,
             case when replica = 10 then 1 else 0 end replicas_on_10,
             case when replica = 12 then 1 else 0 end replicas_on_12,
             case when replica = 13 then 1 else 0 end replicas_on_13
      from replicas where leader = 'true') i
group by topic;
```

## Replica Distribution per Topic

Modify the query below to include the correct pivot columns (one column per node id):

```sql
-- Show leader balance by topic
select topic,
       sum(replicas_on_2), -- modify this list to match your node ids
       sum(replicas_on_3),
       sum(replicas_on_5),
       sum(replicas_on_6),
       sum(replicas_on_8),
       sum(replicas_on_9),
       sum(replicas_on_10),
       sum(replicas_on_12),
       sum(replicas_on_13)
from (select topic,
             case when replica = 2 then 1 else 0 end  replicas_on_2, -- modify this list to match your node ids
             case when replica = 3 then 1 else 0 end  replicas_on_3,
             case when replica = 5 then 1 else 0 end  replicas_on_5,
             case when replica = 6 then 1 else 0 end  replicas_on_6,
             case when replica = 8 then 1 else 0 end  replicas_on_8,
             case when replica = 9 then 1 else 0 end  replicas_on_9,
             case when replica = 10 then 1 else 0 end replicas_on_10,
             case when replica = 12 then 1 else 0 end replicas_on_12,
             case when replica = 13 then 1 else 0 end replicas_on_13
      from replicas) i
group by topic;
```

## Cluster-level Leadership Distribution

Run the following SQL:

```sql
-- Show leader balance
select replica, count(*) from replicas where leader = "true" group by replica;
```

## Cluster-level Replica Distribution

Run the following SQL:

```sql
-- Show leader balance
select replica, count(*) from replicas where leader = "true" group by replica;
```
