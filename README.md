# Ebiten Agent Example

Vampire Survivors 风格的生存游戏。玩家控制角色在地图上击败敌人、收集道具和升级，目标是尽可能长时间生存并获得高分。

## 技术栈

- Go
- Ebiten v2.9.9

## 运行

```bash
go run ./cmd/game
```

## 图画生成

使用 fal.ai FLUX dev API 生成游戏图画。

1. 创建 prompt 文件：`[path-to-image].prompt.md`
2. 设置环境变量：`export FAL_KEY=your-api-key`
3. 生成图画：`go run ./cmd/image-generator [path-to-image].prompt.md`

生成的 PNG 文件将保存至 `[path-to-image].png`。
