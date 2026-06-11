# Ebiten代码规则

**`internal/component/` 禁止 import ebiten。**

游戏逻辑代码应放在不依赖ebiten的包中，确保可直接`go test`。

## 验证方法

```bash
go list -f '{{.Imports}}' ./internal/component/... ./internal/config/... | grep ebiten
# 无输出 = 通过
```

## 各包职责

| 包 | 职责 | ebiten依赖 |
|---|---|---|
| `component` | 数据结构、游戏逻辑、状态计算 | 禁止 |
| `config` | 屏幕尺寸等全局常量 | 禁止 |
| `entity` | 游戏对象，组合component与render | 允许 |
| `scene` | 场景生命周期、entity管理 | 允许 |
| `render` | 精灵、动画等渲染结构 | 允许 |
| `input` | 键盘/鼠标输入封装 | 允许 |
| `game` | ebiten.Game实现、场景切换 | 允许 |

## 设计原则

1. **逻辑与渲染分离**：游戏逻辑放component，可独立测试
2. **组合优于继承**：entity通过嵌入component组合功能
3. **单向依赖**：高层依赖低层（scene→entity→component），禁止反向

## 测试示例

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
