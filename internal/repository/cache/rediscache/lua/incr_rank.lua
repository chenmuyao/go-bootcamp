local key = KEYS[1]
local bizId = ARGV[1]
local delta = tonumber(ARGV[2])

-- Update the score only when the key exists in the sorted set
if redis.call("ZSCORE", key, bizId) then
	redis.call("ZINCRBY", key, delta, bizId)
end

return 0
