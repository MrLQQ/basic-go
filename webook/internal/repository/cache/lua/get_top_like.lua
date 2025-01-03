---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by louqiangqiang.
--- DateTime: 2024/2/29 09:47
---

-- 使用ZRANGE命令获取按like_cnt排序的前topN个文章及其对应的like_cnt（WITHSCORES选项会返回成员及其分数）
local topArticlesWithScores = redis.call("ZREVRANGE", "topLike", 0, topN-1, 'WITHSCORES')

-- 将结果转换为键值对的形式，便于处理（这里返回的是一个二维数组，每一项是[ID, like_cnt]）
local formattedTopArticles = {}
for i = 1, #topArticlesWithScores, 2 do
    table.insert(formattedTopArticles, {topArticlesWithScores[i], tonumber(topArticlesWithScores[i+1])})
end

-- 返回排序后的键（ID）及其like_cnt列表
return formattedTopArticles
