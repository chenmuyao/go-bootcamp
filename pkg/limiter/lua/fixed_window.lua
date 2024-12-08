local key = KEYS[1]
local limit = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])

local timeBeginKey = key .. ":time"
local cntKey = key .. ":cnt"

-- {{{ Get now
local redis_time = redis.call("TIME")

local seconds = tonumber(redis_time[1])
local microseconds = tonumber(redis_time[2])

local now = math.floor((seconds * 1000) + (microseconds / 1000))
-- }}}

local cnt = redis.call("get", cntKey)
local timeBegin = redis.call("get", timeBeginKey)

if not cnt or tonumber(cnt) <= 0 or not timeBegin or tonumber(timeBegin) <= 0 then
	-- Key does not exist
	redis.call("set", cntKey, 1)
	redis.call("set", timeBeginKey, now)
	return 0 -- Accept connection
end

if now - tonumber(timeBegin) > interval then
	redis.call("set", cntKey, 1)
	redis.call("set", timeBeginKey, now)
	return 0 -- Accept connection
end

if cnt + 1 > limit then
	return -1 -- Rate limit
end

redis.call("incr", cntKey)

return 0 -- Accept
