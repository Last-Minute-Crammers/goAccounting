# 问题求解实战项目开发 -- goAccounting server端

## 注意事项
- Gong的WSL中MySQL端口配置为3307，假如push到云端需要修改为3306

## 更新list
---
### 已do list
- 好友申请，添加和拒绝
- 登录与注册（无邮箱验证）
- 测试登录与注册，修改了JWT令牌的ID和Subject
- 基础记账功能+单日静态统计
- ai服务大体框架

---
### todo list
- ai服务
- category服务（要ICON吗还）
- api

---
### consider todo list（暂时优先级低）
- 邮件验证
- user_log
- 创建新记账本，多人账本

---
### 暂时不考虑实现的功能
- NATS消息队列
- 不同类型的设备适配（暂时只以安卓来考虑）
- 无邮箱本地注册
