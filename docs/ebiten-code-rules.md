# Ebiten代码规则

## 各包职责

| 包 | 职责 |
|---|---|
| `component` | 纯数据组件（Position、Velocity、Animation等） |
| `config` | 屏幕尺寸等全局常量 |
| `ecs` | ECS核心：World、Storage、System接口 |
| `ecs/system` | 具体System实现（Input、Movement、Render等） |
| `entity` | 实体工厂函数（组合组件创建实体） |
| `scene` | 场景生命周期、System调度 |
| `game` | ebiten.Game实现、场景切换 |

## ECS架构

采用Entity-Component-System模式：
- **Entity**：仅为ID（uint32 + Generation），无逻辑
- **Component**：纯数据，存储于`Storage[T]`
- **System**：纯逻辑，遍历组件执行行为

### 依赖关系

```
scene → ecs/system → ecs → component
  ↓                   ↑
entity ───────────────┘
```

### System执行顺序

Scene负责按顺序调用System：

```go
// Update顺序决定逻辑正确性
inputSystem.Update(world, dt)          // 1. 读取输入
movementSystem.Update(world, dt)       // 2. 应用移动
facingSystem.Update(world, dt)         // 3. 更新朝向
animationStateSystem.Update(world, dt) // 4. 设置动画状态
animationSystem.Update(world, dt)      // 5. 更新动画帧
cameraSystem.Update(world, dt)         // 6. 跟随摄像机
```

### 创建实体

使用工厂函数组合组件：

```go
entity.CreatePlayer(world, components, entity.PlayerFactoryConfig{
    X: 2.0, Y: 3.0,
    Speed: 5.0,
    // ...
})
```

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
// ✅ 正确：System通过参数接收dt
func (s *MovementSystem) Update(w *ecs.World, dt float64) {
    pos.X += vel.X * dt  // velocity: units/s
}

// ❌ 错误
func (s *MovementSystem) Update(w *ecs.World) {
    pos.X += vel.X / 60.0  // 硬编码帧率
}
frameCount++; if frameCount >= 60 { ... }  // 帧计数器
cooldown := 120  // 用帧数表示2秒
```

**scene层获取deltaTime**：
```go
s.clock.Update(1.0 / float64(ebiten.TPS()))
dt := s.clock.DeltaTime()
s.movementSystem.Update(s.world, dt)
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

