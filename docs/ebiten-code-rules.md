# Ebiten代码规则

## 各包职责

| 包 | 职责 |
|---|---|
| `component` | 游戏对象组件（位置、移动、生命等） |
| `config` | 屏幕尺寸等全局常量 |
| `entity` | 游戏对象，组合component |
| `scene` | 场景生命周期、entity管理 |
| `render` | 精灵、动画等渲染结构 |
| `game` | ebiten.Game实现、场景切换 |

## 设计原则

1. **组合优于继承**：entity通过嵌入component组合功能
2. **单向依赖**：高层依赖低层（scene→entity→component），禁止反向

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

