local key = KEYS[1]
local limit = ARGV[1]
local windowSize = ARGV[2]
-- Pass nanosecond for precision. Otherwise when testing with go routine,
-- there are chance that some requests arrive at the same time and be counted only once.
local now = ARGV[3]

local cutTime = now - tonumber(windowSize)

-- Remove expired entries
redis.call("ZREMRANGEBYSCORE", key, "-inf", cutTime)

-- Count current requests in the window
local count = redis.call("ZCARD", key)

if count >= tonumber(limit) then
	redis.log(redis.LOG_WARNING, "cutTime:" .. cutTime .. " count:" .. count)
	return -1
end

redis.log(redis.LOG_WARNING, "cutTime:" .. cutTime .. " count:" .. count)

redis.call("ZADD", key, now, now)

redis.call("EXPIRE", key, math.ceil(windowSize / 1000) + 1)

return 0
