# 运维AI助手 - Git工作流

> 版本: v1.0
> 更新日期: 2026-03-24

---

## 1. 分支策略

### 1.1 主要分支

- `main`: 生产环境代码，受保护
- `develop`: 开发主分支
- `feature/*`: 功能开发分支
- `bugfix/*`: Bug修复分支
- `hotfix/*`: 紧急修复分支
- `release/*`: 发布准备分支

### 1.2 分支命名规范

```bash
feature/US-1-chat-diagnostic
bugfix/fix-metrics-query-timeout
hotfix/critical-memory-leak
release/v1.0.0
```

---

## 2. 提交规范

### 2.1 Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 2.2 Type 类型

- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

### 2.3 示例

```bash
feat(orchestration): 实现告警分诊Graph编排

- 添加parse_input节点
- 添加route_domain分支节点
- 实现证据收集并行节点

Closes #123
```

---

## 3. 开发流程

```bash
# 1. 从develop创建功能分支
git checkout develop
git pull origin develop
git checkout -b feature/US-1-chat-diagnostic

# 2. 开发并提交
git add .
git commit -m "feat(api): 实现chat接口"

# 3. 推送到远程
git push origin feature/US-1-chat-diagnostic

# 4. 创建Pull Request到develop

# 5. Code Review通过后合并

# 6. 删除功能分支
git branch -d feature/US-1-chat-diagnostic
```

---

## 4. Pull Request规范

### 4.1 PR标题

```
[Feature] 实现对话式诊断功能
[Bugfix] 修复指标查询超时问题
[Hotfix] 紧急修复内存泄漏
```

### 4.2 PR描述模板

```markdown
## 变更说明
简要描述本次变更的目的和内容

## 变更类型
- [ ] 新功能
- [ ] Bug修复
- [ ] 重构
- [ ] 文档更新

## 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手工测试通过

## 相关Issue
Closes #123

## 截图（如适用）
```

---

## 5. 发布流程

```bash
# 1. 创建release分支
git checkout develop
git checkout -b release/v1.0.0

# 2. 更新版本号和CHANGELOG
# 编辑version文件和CHANGELOG.md

# 3. 提交版本变更
git commit -am "chore: 准备v1.0.0发布"

# 4. 合并到main
git checkout main
git merge --no-ff release/v1.0.0

# 5. 打标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 6. 推送
git push origin main --tags

# 7. 合并回develop
git checkout develop
git merge --no-ff release/v1.0.0

# 8. 删除release分支
git branch -d release/v1.0.0
```

---

## 6. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
