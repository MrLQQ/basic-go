local key = KEYS[1]
--拼接字符串，cntKey表示可重试的次数
local cntKey = key..":cnt"
-- 用户输入的验证码
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    -- 验证次数耗尽了
    return -1
end

if code == expectedCode then
    redis.call("set",cntKey,0)
    return 0
else
    -- 不相等，用户输入验证码错误
    redis.call("decr",cntKey)
    return -2
end