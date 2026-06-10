- 注重可读性
- 注重一致性
- 注重最佳实践
- 注重《程序员修炼之道》中的，横跨代码、测试、文档的DRY
- 注重测试，遵循测试金字塔，不盲目追求代码的 coverage，注重规格 coverage
- 注重代码整理，所有逻辑都该放在它应该在的地方
- 注重YAGNI，AI世代开发成本降低，无效代码的危害变大
- 注重文档维护，多余文档以及无效文档都是有害的

## 技術

使用Go Ebiten v2.9.9。
注意要参照最新文档。

## 游戏

类似于Vampire Survivors的游戏，玩家控制一个角色，在一个地图上不断地击败敌人，收集道具和升级，目标是生存尽可能长的时间并获得高分。

## 图画

图画使用AI生成（fal.ai FLUX dev API）。

### 使用方法

1. 创建prompt文件：`[path-to-image].prompt.md`
2. 设置环境变量：`export FAL_KEY=your-api-key`
3. 生成图画：`go run ./cmd/image-generator [path-to-image].prompt.md`

生成的PNG文件将保存至 `[path-to-image].png`。
