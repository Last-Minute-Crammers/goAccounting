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

---
### todo list
- ai服务
    - 财务报告（优先月报）
- api
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

func (s *ChatService) GetChatResponse(userInput string, userId uint, ctx context.Context) (string, error) {
    sessionId := uuid.New().String()
    requestId := uuid.New().String()

    log.Printf("=== 开始AI对话请求 ===")
    log.Printf("RequestID: %s", requestId)
    log.Printf("SessionID: %s", sessionId)
    log.Printf("用户输入: %s", userInput)

    // 定义系统提示词
    systemPrompt := "你的中文名字叫理财小汪，你是智能理财宠物，对用户的称呼是主人。你有着丰富的理财知识，活泼可爱，认真可靠，你需要协助用户进行个人理财规划。当回复问题时需要回复你的名字时，中文名必须回复理财小汪，此外回复和你的名字相关的问题时，也需要给出和你的名字对应的合理回复。"

    // 获取用户的理财数据
    financialData, err := s.getUserFinancialData(userId)
    if err != nil {
        return "", fmt.Errorf("获取用户理财数据失败: %v", err)
    }

    // 获取用户的对话历史
    history := s.ConversationHistory[userId]

    // 动态生成 Prompt
    fullPrompt := fmt.Sprintf("%s\n\n用户的理财数据:\n%s\n\n对话历史:\n%s\n\n用户提问:\n%s",
        systemPrompt, financialData, history, userInput)

    // 构建请求体
    reqBody := blueLMRequest{
        Prompt:       fullPrompt,
        Model:        "vivo-BlueLM-TB-Pro",
        SessionId:    sessionId,
        SystemPrompt: systemPrompt,
        Extra: map[string]interface{}{
            "temperature":  0.9,
            "systemPrompt": systemPrompt,
        },
    }

    bodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("请求体序列化失败: %v", err)
    }

    log.Printf("请求体JSON: %s", string(bodyBytes))

    // 准备查询参数
    queryParams := map[string]string{
        "requestId": requestId,
    }

    // 构建URL
    fullUrl := blueLMApiUrl
    log.Printf("请求URL: %s", fullUrl)

    httpReq, err := http.NewRequestWithContext(ctx, "POST", fullUrl, bytes.NewReader(bodyBytes))
    if err != nil {
        return "", fmt.Errorf("创建HTTP请求失败: %v", err)
    }

    // 添加查询参数
    q := httpReq.URL.Query()
    for k, v := range queryParams {
        q.Add(k, v)
    }
    httpReq.URL.RawQuery = q.Encode()

    // 设置请求头
    httpReq.Header.Set("Content-Type", "application/json")

    // 生成认证头
    authHeaders := GenerateAuthHeaders("POST", "/vivogpt/completions", queryParams, blueLMAppID, blueLMAppKey)
    for key, value := range authHeaders {
        httpReq.Header.Set(key, value)
    }

    // 发送请求
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(httpReq)
    if err != nil {
        return "", fmt.Errorf("HTTP请求失败: %v", err)
    }
    defer resp.Body.Close()

    respBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("读取响应失败: %v", err)
    }

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("HTTP状态码错误: %d, 响应: %s", resp.StatusCode, string(respBytes))
    }

    var apiResp blueLMResponse
    if err := json.Unmarshal(respBytes, &apiResp); err != nil {
        return "", fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(respBytes))
    }

    if apiResp.Code != 0 {
        return "", fmt.Errorf("API调用失败: %s (code: %d)", apiResp.Msg, apiResp.Code)
    }

    // 更新对话历史
    s.ConversationHistory[userId] = fmt.Sprintf("%s\n用户: %s\nAI: %s", history, userInput, apiResp.Data.Content)

    // 保存对话记录
    if err := s.saveChatRecord(userId, userInput, apiResp.Data.Content); err != nil {
        log.Printf("保存对话记录时出错: %v", err)
    }

    return apiResp.Data.Content, nil
}