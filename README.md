# LaunchPad

[![Build Status](https://travis-ci.com/FreakyGranny/launchpad-api.svg?branch=master)](https://travis-ci.com/FreakyGranny/launchpad-api) [![Go Report Card](https://goreportcard.com/badge/github.com/FreakyGranny/launchpad-api)](https://goreportcard.com/report/github.com/FreakyGranny/launchpad-api) [![codecov](https://codecov.io/gh/FreakyGranny/launchpad-api/branch/master/graph/badge.svg)](https://codecov.io/gh/FreakyGranny/launchpad-api)

The task that the project solves is raising money to buy something for joint ownership, or finding people for some kind of activity.

The service is suitable for small communities, such as an office, for example, it is assumed that the service is deployed on an internal network, thereby limiting the audience of projects.

LaunchPad project consist of:
=============================

* launchpad-api (this repo)
* [launchpad-gui](https://github.com/FreakyGranny/launchpad-gui)

AUTH
===========================

* [Create Create VK auth application](https://vk.com/editapp?act=create)
* Use app credentionals with [docker-compose](./deployments/docker-compose.yml)

Prepare first start
===================

init migrations
```
lpad migrate init
```

apply migrations
```
lpad migrate up
```

Run API
=======

```
lpad api
```
