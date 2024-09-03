local json = require ("dkjson")

local data = require(arg[1])

local str = json.encode (data, { indent = true })

print(str)
