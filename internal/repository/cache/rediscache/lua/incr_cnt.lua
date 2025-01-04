-- biz
local key = KEYS[1]

-- read,like or collect
local cntKey = ARGV[1]

local delta = tonumber(ARGV[2])

local exist = redis.call("EXISTS", key, cntKey, delta)

if exist == 1 then
	redis.call("HINCRBY", key, cntKey, delta)
	return 1
else
	return 0
end
