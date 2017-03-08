# esnctl

[![Build Status](https://travis-ci.org/dtan4/esnctl.svg?branch=master)](https://travis-ci.org/dtan4/esnctl)
[![codecov](https://codecov.io/gh/dtan4/esnctl/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/esnctl)

Elasticsearch Node Controller with AWS Auto Scaling Group

## Why

Graceful Elasticsearch node addition/removal requires several steps.

### Add node

1. Disable shard reallocation
    - see https://www.elastic.co/guide/en/elasticsearch/reference/1.4/cluster-nodes-shutdown.html
2. Add node
3. Enable shard reallocation

### Remove node

1. Remove node from load balancer
2. Wait for connection draining
3. Remove node from shard allocation targets
    - see http://stackoverflow.com/a/23905040
4. Wait for that shards on target node escape to other nodes
5. (Es 1.x only) Shut down node

So far we have conducted this by hand. However, it sometimes causes operation errors.
We realize that these operations should be automated and conducted by ONE action.

## Required environment

- Elasticsearch 1.x / 2.x
- Elasticsearch cluster is running on __AWS EC2 instances__
  - Using [EC2 Discovery](https://www.elastic.co/guide/en/elasticsearch/plugins/current/discovery-ec2-discovery.html)
- EC2 instances are managed by __AWS Auto Scaling Groups__
  - Instances (= Nodes) can be added/removed by modifying DesiredCapacity
- EC2 instances and Auto Scaling Group are attached to __Target Group__
  - Cluster can be accessed through Application Load Balancer

(TODO: architecture image here)

## Installation

### Precompiled binary

Precompiled binaries for Windows, OS X, Linux are available at [Releases](https://github.com/dtan4/esnctl/releases).

### From source

```bash
$ go get -d github.com/dtan4/esnctl
$ cd $GOPATH/src/github.com/dtan4/esnctl
$ make deps
$ make install
```

## Usage

To run `esnctl add` or `esnctl remove`, you need to set valid AWS credentials beforehand.

```bash
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=xx-yyyy-0
```

### `esnctl list`

List nodes

```bash
$ esnctl add \
  --cluster-url http://elasticsearch.example.com
ip-10-0-1-21.ap-northeast-1.compute.internal
ip-10-0-1-22.ap-northeast-1.compute.internal
ip-10-0-1-23.ap-northeast-1.compute.internal
```

|Option|Description|
|---------|-----------|
|`--cluster-url=CLUSTERURL`|Elasticsearch cluster URL|

### `esnctl add`

Add nodes

```bash
$ esnctl add \
  --cluster-url http://elasticsearch.example.com \
  --group elasticsearch \
  -n 2
===> Disabling shard reallocation...
===> Launching 2 instances on elasticsearch...
===> Waiting for nodes join to Elasticsearch cluster...
........................
===> Enabling shard reallocation...
===> Finished!
```

|Option|Description|
|---------|-----------|
|`--group=GROUP`|Auto Scaling Group|
|`--cluster-url=CLUSTERURL`|Elasticsearch cluster URL|
|`-n`, `--number=NUMBER`|Number to add instances|
|`--region=REGION`|AWS region|

### `esnctl remove`

Remove a node

Only 1 node can be removed at the same time.

```bash
$ esnctl remove \
  --cluster-url http://elasticsearch.example.com \
  --group elasticsearch \
  --node-name ip-10-0-1-21.ap-northeast-1.compute.internal
===> Retrieving target instance ID...
===> Retrieving target group...
===> Detaching instance from target group...
............................................................
===> Excluding target node from shard allocation group...
===> Waiting for shards escape from target node...
..................
===> Shutting down target node...
===> Detaching target instance...
===> Finished!
```

|Option|Description|
|---------|-----------|
|`--group=GROUP`|Auto Scaling Group|
|`--cluster-url=CLUSTERURL`|Elasticsearch cluster URL|
|`--node-name=NODENAME`|Elasticsearch node name to remove|
|`--region=REGION`|AWS region|

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
