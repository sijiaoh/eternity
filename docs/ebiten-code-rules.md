# Ebiten代码规则

## 各包职责

| 包 | 职责 |
|---|---|
| `component` | 游戏对象组件（位置、移动、生命等） |
| `config` | 屏幕尺寸等全局常量 |
| `entity` | 游戏对象，组合component |
| `scene` | 场景生命周期、entity管理 |
| `game` | ebiten.Game实现、场景切换 |

## 设计原则

1. **组合优于继承**：entity通过嵌入component组合功能
2. **单向依赖**：高层依赖低层（scene→entity→component），禁止反向

## 单位系统

使用世界单位（units）而非像素，便于物理计算和数值调整。

| 常量/函数 | 说明 |
|-----------|------|
| `PixelsPerUnit = 48` | 1单位 = 48像素（与角色/tile尺寸一致） |
| `UnitsToPixels()` | 绘制时转换位置 |
| `PixelsToUnits()` | 初始化时转换像素值 |

### 单位边界

| 层级 | 单位 | 示例 |
|------|------|------|
| 游戏逻辑 | 单位 | Position、Speed（5 units/s） |
| Ebiten接口 | 像素 | ScreenWidth、Draw坐标 |
| 视觉效果 | 像素 | Camera.GetOffset()、渲染偏移量 |

```go
// ✅ 正确：位置和速度使用单位
pos := component.NewPosition(2.0, 3.0) // 世界坐标 (2, 3) 单位
speed := 5.0                            // 5 units/s

// ✅ 绘制时转换为像素
x := config.UnitsToPixels(pos.X) // 用于 DrawImage
```

## 时间系统

**禁止依赖固定帧率**。所有游戏逻辑必须基于时间（秒），而非帧数。

### deltaTime规范

```go
// ✅ 正确
func (p *Player) Update(dt float64) {
    p.Position.X += p.velocity * dt  // velocity: units/s
}

// ❌ 错误
func (p *Player) Update() {
    p.Position.X += p.velocity / 60.0  // 硬编码帧率
}
frameCount++; if frameCount >= 60 { ... }  // 帧计数器
cooldown := 120  // 用帧数表示2秒
```

**scene层获取deltaTime**：
```go
s.clock.Update(1.0 / float64(ebiten.TPS()))
dt := s.clock.DeltaTime()
player.Update(dt)
```

### 时间缩放

使用`Clock`实现暂停、慢动作等效果。

| Scale值 | 效果 |
|---------|------|
| 0 | 暂停 |
| 0.5 | 慢动作 |
| 1.0 | 正常 |
| 2.0 | 快进 |

**暂停与恢复**：`Pause()`保存当前scale，`Resume()`恢复之前的scale。
```go
clock.SetScale(0.5)  // 慢动作
clock.Pause()        // 暂停，保存scale=0.5
clock.Resume()       // 恢复scale=0.5（而非默认1.0）
```

**嵌套时钟**用于独立控制不同系统的时间：
```go
gameClock := component.NewChildClock(root)  // 游戏逻辑
uiClock := component.NewChildClock(root)    // UI动画

gameClock.SetScale(0.5) // 游戏慢动作，UI正常
```

