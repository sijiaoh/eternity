# Eternity

Vampire Survivors 风格的生存游戏。玩家控制角色在地图上击败敌人、收集道具和升级，目标是尽可能长时间生存并获得高分。

## 技术栈

- Go
- Ebiten v2.9.9

## 开发

### 环境初始化

```bash
./scripts/setup.sh
```

自动安装Go依赖和平台特定工具（Ubuntu上安装Xvfb用于headless测试）。

### 运行

```bash
go run ./cmd/game
```

### 测试

```bash
./scripts/test.sh                     # 运行所有测试
./scripts/test.sh ./internal/entity/  # 指定测试目标
```

Ubuntu上自动使用Xvfb进行headless执行。

### 代码检查

```bash
./scripts/lint.sh
```

运行gofmt、go vet及golangci-lint（如已安装）。

## 图画生成

使用 fal.ai FLUX dev API 生成游戏图画。

1. 创建 prompt 文件：`[path-to-image].prompt.md`
2. 设置 FAL_KEY（二选一）：
   - 环境变量：`export FAL_KEY=your-api-key`
   - .env 文件：在项目根目录创建 `.env` 文件，内容为 `FAL_KEY=your-api-key`
3. 生成图画：`go run ./cmd/image-generator [path-to-image].prompt.md`

生成的 PNG 文件将保存至 `[path-to-image].png`。

> 注：环境变量优先级高于 .env 文件。.env 文件已在 .gitignore 中，不会被提交。
