local key = KEYS[1]
local releaseAmount = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])
local capacity = tonumber(ARGV[3])

-- {{{ Get now

local redis_time = redis.call("TIME")

local seconds = tonumber(redis_time[1])
local microseconds = tonumber(redis_time[2])

local now = math.floor((seconds * 1000) + (microseconds / 1000))

-- }}}

local rateInfo = redis.call("HMGET", key, "lastReleaseTime", "tokens")
local lastReleaseTime = tonumber(rateInfo[1])
local tokens = tonumber(rateInfo[2])

if not lastReleaseTime then
	tokens = capacity - 1
	lastReleaseTime = now
	redis.call("HMSET", key, "lastReleaseTime", lastReleaseTime, "tokens", tokens)
	return 0
end

-- {{{ release

local timePassed = now - lastReleaseTime
local intervalPassed = math.floor(timePassed / interval)
local shouldRelease = intervalPassed * releaseAmount

tokens = tokens + shouldRelease

if tokens >= capacity then
	tokens = capacity
end

lastReleaseTime = now

-- }}}

tokens = tokens - 1

if tokens < 0 then
	tokens = 0
	redis.call("HMSET", key, "lastReleaseTime", lastReleaseTime, "tokens", tokens)
	return -1
end

redis.call("HMSET", key, "lastReleaseTime", lastReleaseTime, "tokens", tokens)
return 0
