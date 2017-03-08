# esnctl

[![Build Status](https://travis-ci.org/dtan4/esnctl.svg?branch=master)](https://travis-ci.org/dtan4/esnctl)
[![codecov](https://codecov.io/gh/dtan4/esnctl/branch/master/graph/badge.svg)](https://codecov.io/gh/dtan4/esnctl)

Elasticsearch Node Controller with AWS AutoScaling Group

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

- Elasticsearch 1.x / 2.x / 5.x
- Elasticsearch cluster is running on __AWS EC2 instances__
  - Using [EC2 Discovery](https://www.elastic.co/guide/en/elasticsearch/plugins/current/discovery-ec2-discovery.html)
- EC2 instances are managed by __AWS Auto Scaling Groups__
  - Instances (= Nodes) can be added/removed by modifying DesiredCapacity

## Installation

TBD

## Usage

### `esnctl list`

### `esnctl add`

### `esnctl remove`

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
