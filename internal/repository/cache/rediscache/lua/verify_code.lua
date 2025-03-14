local key = KEYS[1]

local cntKey = key .. ":cnt"

local expectedCode = ARGV[1]

local cnt = redis.call("get", cntKey)
local code = redis.call("get", key)

if cnt == nil or tonumber(cnt) <= 0 then
	-- no more reries
	return -1
end

if code == expectedCode then
	redis.call("set", cntKey, 0)
	return 0
else
	-- Wrong code, decrease retry count
	redis.call("decr", cntKey)
	return -2
end
