local key = KEYS[1]
--拼接字符串，cntKey表示可重试的次数
local cntKey = key..":cnt"
-- 准备存储的验证码
local val = ARGV[1]

-- ttl的返回值如果是-2 表示目标key不存在
-- ttl的返回值如果是-1 表示目标key存在，但是没有设置剩余生存时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 目标key不存在，获取目标key的剩余时间小于9分钟(60 * 9 = 540)
    -- 可以发送验证码
    -- 将验证码set到redis
    redis.call("set",key, val)
    -- 设置验证码的过期时间10分钟(60 * 10 = 600)
    redis.call("expire", key, 600)
    -- 设置验证码验证次数为3次,已经对应的过期时间
    redis.call("set",cntKey,3)
    redis.call("expire", cntKey, 600)
else
    -- 发送太频繁
    return -1
end