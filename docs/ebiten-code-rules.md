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

## 时间系统

### deltaTime规范

所有与时间相关的计算必须使用deltaTime，不能依赖固定帧率。

```go
// ✅ 正确：速度单位为"像素/秒"
speed := 240.0 // pixels per second
dx := speed * deltaTime

// ❌ 错误：假设固定帧率
dx := 4.0 // 假设60fps
```

**scene层获取deltaTime**：
```go
s.clock.Update(1.0 / float64(ebiten.TPS()))
dt := s.clock.DeltaTime()
player.Update(dt)
```

注：Ebitengine使用固定时间步长，使用`1.0/ebiten.TPS()`而非`ActualTPS()`。

### 时间缩放

使用`Clock`实现暂停、慢动作等效果。

| Scale值 | 效果 |
|---------|------|
| 0 | 暂停 |
| 0.5 | 慢动作 |
| 1.0 | 正常 |
| 2.0 | 快进 |

**嵌套时钟**用于独立控制不同系统的时间：
```go
gameClock := component.NewChildClock(root)  // 游戏逻辑
uiClock := component.NewChildClock(root)    // UI动画

gameClock.SetScale(0.5) // 游戏慢动作，UI正常
```

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
