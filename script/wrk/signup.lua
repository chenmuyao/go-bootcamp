wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local random = math.random
local function uuid()
	local template = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	return string.gsub(template, "[xy]", function(c)
		local v = (c == "x") and random(0, 0xf) or random(8, 0xb)
		return string.format("%x", v)
	end)
end

function init(args)
	-- The `tostring({}):sub(8)` trick works because:
	--
	-- 1. `{}` creates a new table object
	-- 2. `tostring(table)` generates a unique string representation of that table's memory address
	-- 3. This string typically looks like "table: 0x123456"
	-- 4. `:sub(8)` extracts the numeric part of this address
	-- 5. When converted to a number, this provides a unique identifier for each thread/call
	--
	-- So it's essentially using the table's memory address as a source of randomness, which is different each time it's called.
	math.randomseed(tonumber(tostring({}):sub(8)) + os.time())
	cnt = 0
	prefix = uuid()
end

function request()
	body = string.format(
		'{"email":"%s%d@test.com", "password":"hello#world123", "confirm_password": "hello#world123"}',
		prefix,
		cnt
	)
	cnt = cnt + 1
	return wrk.format("POST", wrk.path, wrk.headers, body)
end

function response() end

