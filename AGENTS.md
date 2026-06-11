# 开发原则

- 注重可读性
- 注重一致性
- 注重最佳实践
- 注重《程序员修炼之道》中的，横跨代码、测试、文档的DRY
- 注重测试，遵循测试金字塔，不盲目追求代码的 coverage，注重规格 coverage
- 注重代码整理，所有逻辑都该放在它应该在的地方
- 注重YAGNI，AI世代开发成本降低，无效代码的危害变大
- 注重文档维护，多余文档以及无效文档都是有害的

## 技术栈

- Ebiten v2.9.9
- 图画生成：fal.ai FLUX dev API

参照最新文档。

## Ebiten代码规则

**`internal/component/` 禁止 import ebiten。**

游戏逻辑代码应放在不依赖ebiten的包中，确保可直接`go test`。

### 验证方法

```bash
go list -f '{{.Imports}}' ./internal/component/... ./internal/config/... | grep ebiten
# 无输出 = 通过
```

### 架构

```
internal/
├── component/    # 禁止ebiten，纯数据结构
├── config/       # 禁止ebiten，全局配置常量
├── game/         # 可用ebiten
├── input/        # 可用ebiten
├── render/       # 可用ebiten
├── entity/       # 可用ebiten
└── scene/        # 可用ebiten
```

### 测试示例

```go
// internal/component/health_test.go
func TestHealth_TakeDamage(t *testing.T) {
    h := &Health{Current: 100, Max: 100}
    h.TakeDamage(30)
    if h.Current != 70 {
        t.Errorf("expected 70, got %d", h.Current)
    }
}
```

直接 `go test ./internal/component/...`，无需图形环境。
