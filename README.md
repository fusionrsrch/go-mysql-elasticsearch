go-mysql-elasticsearch is a service syncing your MySQL data into Elasticsearch automatically. 

It uses `mysqldump` to fetch the origin data at first, then syncs data incrementally with binlog.

## Dependencies
+ Go and $GOPATH, [GOPATH](https://golang.org/doc/code.html#GOPATH)
+ Godep, `go get github.com/tools/godep`

## Install 
+ `go get github.com/fusionrsrch/go-mysql-elasticsearch`, it will print some messages in console, skip it. :-)
+ `cd $GOPATH/src/github.com/fusionrsrch/go-mysql-elasticsearch`
+ `make`

## How to use?

+ Create table in MySQL.
+ Create the associated Elasticsearch index, document type and mappings if possible, if not, Elasticsearch will create these automatically.
+ Config base, see the example config [river.toml](./etc/river.toml).
+ Set MySQL source in config file, see [Source](#source) below.
+ Customize MySQL and Elasticsearch mapping rule in config file, see [Rule](#rule) below.
+ Start `./bin/go-mysql-elasticsearch -config=./etc/river.toml` and enjoy it.

## Notes

+ binlog_format must be set to **row** on master.

```
mysql> show variables like 'binlog_format';
+---------------+-------+
| Variable_name | Value |
+---------------+-------+
| binlog_format | ROW   |
+---------------+-------+
1 row in set (0.00 sec)
```

+ binlog_row_image must be set to **full** on master for MySQL. You may lose field data when binlog_row_image is set to minimal or noblob. MariaDB only supports full row image. 
```
mysql> show variables like 'binlog_row_image';
+------------------+-------+
| Variable_name    | Value |
+------------------+-------+
| binlog_row_image | FULL  |
+------------------+-------+
1 row in set (0.00 sec)
```

## Source

In go-mysql-elasticsearch, you must decide which tables you want to sync into elasticsearch in the source config.

The format in config file is below:

```
[[source]]
schema = "test"
tables = ["t1", t2]

[[source]]
schema = "test_1"
tables = ["t3", t4]
```

`schema` is the database name, and `tables` includes the table need to be synced. 

## Rule

By default, go-mysql-elasticsearch will use MySQL table name as the Elasticserach's index and type name, use MySQL table field name as the Elasticserach's field name. 
e.g, if a table named blog, the default index and type in Elasticserach are both named blog, if the table field named title, 
the default field name is also named title.

Rule can let you change this name mapping. Rule format in config file is below:

```
[[rule]]
schema = "test"
table = "t1"
index = "t"
type = "t"

    [rule.field]
    title = "my_title"
```

In the example above, we will use a new index and type both named "t" instead of default "t1", and use "my_title" instead of field name "title".

## Wildcard table

go-mysql-elasticsearch only allows you determind which table to be synced, but sometimes, if you split a big table into multi sub tables, like 1024, table_0000, table_0001, ... table_1023, it is very hard to write rules for every table. 

go-mysql-elasticserach supports using wildcard table, e.g:

```
[[source]]
schema = "test"
tables = ["test_river_[0-9]{4}"]

[[rule]]
schema = "test"
table = "test_river_[0-9]{4}"
index = "river"
type = "river"
```

"test_river_[0-9]{4}" is a wildcard table definition, which represents "test_river_0000" to "test_river_9999", at the same time, the table in the rule must be same as it. 

At the above example, if you have 1024 sub tables, all tables will be synced into Elasticsearch with index "river" and type "river". 



