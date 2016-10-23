<p align="center">
  <img src="https://raw.githubusercontent.com/d3estudio/digest/master/docs/logo.png" /><br/>
  <a href="https://goreportcard.com/report/github.com/d3estudio/digest"><img src="https://goreportcard.com/badge/github.com/d3estudio/digest"></a>
  <img alt="Language" src="https://img.shields.io/badge/language-Go-blue.svg" />
  <img alt="License" src="https://img.shields.io/badge/license-MIT-blue.svg" />
  <a href="https://microbadger.com/images/d3estudio/digest" title="Get your own image badge on microbadger.com"><img src="https://images.microbadger.com/badges/image/d3estudio/digest.svg"></a>
  <a href="https://microbadger.com/images/d3estudio/digest" title="Get your own version badge on microbadger.com"><img src="https://images.microbadger.com/badges/version/d3estudio/digest.svg"></a>A
</p>
> **Warning**: Digest is a work in progress.

Here at [D3 Estudio](http://d3.do), Slack is our primary communication channel. Found a nice article? Slack. Found a cool gif? Slack. Want to spread the word about an event? Well, guess what? Slack.

We have been overusing the (relatively) new [Reactions](http://slackhq.com/post/123561085920/reactions) feature lately, and we stumbled on a nice idea: _Why not create a digest based on these reactions_?

Well, this is what **Digest** does: it watches channels through a bot and stores messages and their respective reactions.

## Installing

Digest is distributed through a single Docker Image containing all required binaries:

```
$ docker pull d3estudio/digest:latest
```

In order to work, Digest needs a Mongo and Redis instance to connect to. Those servers addresses can be defined in the environment variables of the containers running each process:

 - `DIGEST_TOKEN`: Slack token used to authenticate and connect to your Slack team. This token can be obtained by [Creating a New Slack Bot](https://my.slack.com/apps/A0F7YS25R-bots).
 - `DIGEST_SILENCER_EMOJIS`: A comma-separated string of reactions that will prevent a given message that contains it from appearing publicly. Please do note that this does not prevent a message from being stored. Defaults to `no_entry_sign`.
 - `DIGEST_REDIS_SERVER`: Address of the Redis server used for IPC. Defaults to `redis`.
 - `DIGEST_REDIS_PORT`: Port of the Redis server used for IPC. Defaults to `6379`.
 - `DIGEST_MONGO_SERVER`: URL of the Mongo server used for storing data and keeping indexes. It must be in following format `[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]`, which allows it from being simple as `localhost` to something more involved, like `mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb`. Assumes `27017` as the default port, if not provided. Defaults to `mongo`.
 - `DIGEST_MONGO_DATABASE`: Name of the database that will hold collections created by the application. Defaults to `digest`.
 - `DIGEST_BASE_QUEUE_NAME`: Name of the base name of all Redis queues. Useful if working with more than one instance on the same Redis server. Defaults to `Digest`.
 - `DIGEST_TWITTER_KEY` and `DIGEST_TWITTER_SECRET`: Used by the `prefetcher` process to acquire oembed data from Twitter. Those keys are optional, but not providing them will prevent Digest from embeding tweets into the results. Defaults to `""`.

## Developing
The development environment only requires a functional installation of the Go compiler and [Godep](https://github.com/tools/godep), that is used to manage our dependencies. That said, everything you have to do is clone this repository into your `$GOPATH`, head over to it and run `godep restore`. This repository also contains a basic `docker-compose.yml` that will spin up a Redis and Mongo instance.

## License
```
MIT License

Copyright (c) 2016 Victor Gama

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN D3 Estudio

Newspaper Icon by ✦ Shmidt Sergey ✦ from the Noun Project
```
