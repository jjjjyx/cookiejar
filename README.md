# cookiejar
完全拷贝 golang 的原生库 `net/http/cookiejar`   
改动范围:   
1. 增加了对全部域的cookie的导出
2. 导出cookie 包含 `SameSite` `Secure` `HttpOnly` `HostOnly` `Expires` 属性
3. 修改 `entry` 结构体 增加 `CanonicalHost` 缓存计算结果，减少在排序时的重复计算
4. 保留原始提交的cookie 的 maxAge 属性


> 2023年10月17日
