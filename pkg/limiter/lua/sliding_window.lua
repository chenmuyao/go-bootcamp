local key = KEYS[1]
local limit = ARGV[1]
local windowSize = ARGV[2]

-- {{{ Get now
local redis_time = redis.call("TIME")

local seconds = tonumber(redis_time[1])
local microseconds = tonumber(redis_time[2])

local now = math.floor((seconds * 1000) + (microseconds / 1000))
-- }}}

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
