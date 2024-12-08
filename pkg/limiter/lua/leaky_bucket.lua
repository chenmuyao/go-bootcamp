local key = KEYS[1]
local limit = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])
local capacity = tonumber(ARGV[3])

-- {{{ Get now

local redis_time = redis.call("TIME")

local seconds = tonumber(redis_time[1])
local microseconds = tonumber(redis_time[2])

local now = math.floor((seconds * 1000) + (microseconds / 1000))

-- }}}

local rateInfo = redis.call("HMGET", key, "lastLeakTime", "water")
local lastLeakTime = tonumber(rateInfo[1])
local water = tonumber(rateInfo[2])

if not lastLeakTime then
	water = 1
	lastLeakTime = now
	redis.call("HMSET", key, "lastLeakTime", lastLeakTime, "water", water)
	return 0
end

-- {{{ leak

local timePassed = now - lastLeakTime
local intervalPassed = math.floor(timePassed / interval)
local shouldLeak = intervalPassed * limit

water = water - shouldLeak

if water < 0 then
	water = 0
end

lastLeakTime = now

-- }}}

water = water + 1

if water > capacity then
	water = capacity
	redis.call("HMSET", key, "lastLeakTime", lastLeakTime, "water", water)
	return -1
end

redis.call("HMSET", key, "lastLeakTime", lastLeakTime, "water", water)
return 0
