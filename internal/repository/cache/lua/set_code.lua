local key = KEYS[1]

local cntKey = key .. ":cnt"

local val = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
	-- key exists but no associated expiration time
	return -2
elseif ttl == -2 or ttl < 540 then
	-- send verification code
	redis.call("set", key, val)
	-- 600 secs
	redis.call("expire", key, 600)
	redis.call("set", cntKey, 3)
	redis.call("expire", cntKey, 600)
else
	-- too frequent requests
	return -1
end
