# 问题求解实战项目开发 -- goAccounting server端

## 注意事项
- docker compose up -d出错时，可以试一下docker compose down随后再docker compose up -d --build完全重构
- 这时有时候会出现网络问题，那就可以先通过docker pull golang:1.24-alpine和docker pull alpine:latest手动拉取镜像

## 前后端协作说明

- 前端项目（ai-finance-frontend）通过 HTTP 请求本服务的 API（默认端口8080）。
- 启动本服务后，确保前端的 API 地址指向本服务（可在前端 `uni.scss` 文件中配置）。
- 推荐开发时前后端分别运行，生产环境可通过 Nginx 等方式反向代理。

## 更新list
---
### 已do list
- 好友申请，添加和拒绝
- 登录与注册（无邮箱验证）
- 测试登录与注册，修改了JWT令牌的ID和Subject
- 基础记账功能+单日静态统计
- ai服务大体框架
- category服务
- api

---
### todo list
- ai服务
    - 财务报告（优先月报）
- 数据统计
    - 单日统计
    - 周、月、年 收支统计
    - 周、月、年 细分类别收支统计
    - 周、月、年 与上一单位时间收支相比变化百分比


---
### consider todo list（暂时优先级低）
- 单日统计分为收入和支出
- 单日统计接入LLM
- 设定头像
- 邮件验证
- user_log
- 创建新记账本，多人账本
- 用户邮箱，密码，昵称的更改
- 好友排名
- 称号成就系统

---
### 暂时不考虑实现的功能
- NATS消息队列
- 不同类型的设备适配（暂时只以安卓来考虑）
- 无邮箱本地注册
