package utils // FIXME: Yeah, utils. I know. It's late and I can't think of something better.

import "github.com/d3estudio/digest/shared/redis"

// Redis is a shared instance of the redis.Client instance, used by
// the entrypoint and processors.
var Redis redis.Client
