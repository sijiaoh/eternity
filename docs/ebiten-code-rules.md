# Ebiten代码规则

## 各包职责

| 包 | 职责 |
|---|---|
| `component` | 纯数据组件（Position、Velocity、Animation、Faction等） |
| `config` | 屏幕尺寸等全局常量 |
| `ecs` | ECS核心：World、Storage、System接口 |
| `ecs/system` | 具体System实现（Input、Movement、Render等） |
| `entity` | 实体工厂函数（组合组件创建实体） |
| `scene` | 场景接口与生命周期、`Manager`场景切换、System调度 |
| `game` | ebiten.Game实现、组装场景与切换回调 |
| `i18n` | 多语言文本：key→文本查表、locale 数据文件内嵌、默认语言回退 |
| `scenario` | 调试启动配置的格式与解析（纯数据，无 ebiten）：选哪个场景 + 场景内初始情境 |

## 调试启动

仅供人工调试/测试，不影响正常玩家启动（不传参即走标题页）。两种入口互斥（同传报错）：

- `-scene <name>`：直达某场景，套用其默认情境。可用场景名由 `internal/game` 注册表登记（场景名 const 为单一来源），当前为 `title`、`battle`。
- `-scenario <path>`：加载 JSON 文件，自包含地复现「场景 + 场景内具体情境」（参考 VS Code launch.json）。`scene` 取值同上；顶层 `locale` 与场景无关，对任意被选中的场景（含 title）生效，可按指定语言复现任意场景；情境字段含义与默认值见 `internal/scenario` 的 Go doc，省略即回退正常游玩默认值。

最小示例（英文、进场即对话、慢动作、不生成 goblin）：

```json
{"scene": "battle", "locale": "en", "battle": {"dialogue": true, "timeScale": 0.5, "goblin": false}}
```

未知场景名、未知字段、非法 locale、以及为非 battle 场景设置 battle 情境，都会清晰报错。`locale` 在 `game.New` 构建被选中场景前应用到共享 bundle（与场景无关，故 title/对话等都用对语言）；场景构建经 `BattleSceneConfig.Situation` 接收情境，应用逻辑（位置、对话、goblin 开关、time scale）放在 `scene` 包无构建标签的 `battle_situation.go`，故可无头测试。

## 场景切换

`Manager`持有当前`Scene`，`SetScene`切换。场景间切换用注入回调解耦——场景不直接依赖`Manager`或目标场景类型，由`game`在组装时注入切换动作（如`TitleScene`的`onStart`换入`BattleScene`）。

## ECS架构

采用Entity-Component-System模式：
- **Entity**：仅为ID（uint32 + Generation），无逻辑。ID从1开始，零值`Entity{}`永远无效
- **Component**：纯数据，存储于`Storage[T]`
- **System**：纯逻辑，遍历组件执行行为

### World API

| 方法 | 说明 |
|------|------|
| `Spawn()` | 创建单个实体 |
| `Despawn(e)` | 移除实体，回收槽位 |
| `Alive(e)` | 检查实体是否存活且Generation匹配 |
| `AllAlive()` | 返回所有存活实体的切片 |
| `SpawnBatch(n)` | 批量创建n个实体 |
| `DespawnBatch(es)` | 批量移除实体 |
| `Count()` | 返回存活实体数量 |

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
aiTargetingSystem.Update(world, dt)    // 2. 按阵营选最近敌人，写入AIFollow.Target
aiFollowSystem.Update(world, dt)       // 3. AI跟随（朝目标更新敌人速度）
// 4. 友军按状态行动：普通态尾随，战斗态各自AI（敌人进入主角guardRange即战斗）
updateParty(combatStateSystem.InCombat(world), trailSystem, rangedAISystem, world, dt)
movementSystem.Update(world, dt)       // 5. 应用移动
leashSystem.Update(world, dt)          // 6. 约束：把友军钳回主角guardRange内（仅战斗态生效）
facingSystem.Update(world, dt)         // 7. 更新朝向
animationStateSystem.Update(world, dt) // 8. 设置动画状态
animationSystem.Update(world, dt)      // 9. 更新动画帧
cameraSystem.Update(world, dt)         // 10. 跟随摄像机
```

### 创建实体

工厂函数只组合出通用角色（位置、移动、朝向、动画、精灵），不绑定「谁被玩家操控/被摄像机跟随/属哪个阵营」——这些是场景层决策。场景在组装时按需附加 `InputControlled`（移动速度随控制语义放在 `InputControlled.Speed`）、`CameraTarget` 与 `Faction`（敌我归属，决定 AI 目标锁定）。同一 `CreateMage` 通用角色，场景附 `InputControlled` 即成玩家，改附 `Trail`（普通态尾随）/ `RangedAI`+`Leash`（战斗态远程）即成 AI 驱动的友军（两态切换见上文 System 执行顺序）：

```go
mage := entity.CreateMage(world, components, entity.MageFactoryConfig{
    X: 2.0, Y: 3.0,
    // ...
})
// 由场景决定操控者、摄像机目标与阵营
inputs.Set(mage, component.InputControlled{Speed: 5.0})
cameraTargets.Set(mage, component.CameraTarget{})
factions.Set(mage, component.FactionPlayer)
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
| 视觉效果 | 像素 | system.CameraGetOffset()、渲染偏移量 |

```go
// ✅ 正确：位置和速度使用单位
pos := component.NewPosition(2.0, 3.0) // 世界坐标 (2, 3) 单位
speed := 5.0                            // 5 units/s

// ✅ 绘制时转换为像素
x := config.UnitsToPixels(pos.X) // 用于 DrawImage
```

### Sprite缩放

`SpriteSheet.SizeInUnits`指定sprite的目标宽度（世界单位），渲染时自动计算缩放比例。

| SizeInUnits | 效果 |
|-------------|------|
| 0（默认） | 原尺寸（scale=1.0） |
| 1.0 | 宽度为1单位（48像素） |
| 2.0 | 宽度为2单位（96像素） |

```go
// 创建实体时指定目标尺寸
entity.CreateMage(world, components, entity.MageFactoryConfig{
    SizeInUnits: 1.5, // 显示为1.5个世界单位宽
})
```

缩放比例计算：`scale = UnitsToPixels(SizeInUnits) / FrameWidth`

## 动画系统

三层分离，使用方只设朝向，无需知道自己有几个方向的素材：

| 层 | 类型 | 职责 |
|---|---|---|
| `Animation` | component | 纯帧播放：按 `CurrentState` 顺序播帧，无任何方向语义 |
| `DirectionalAnimation` | component | 解析层：按角色的方向集，把「移动状态 + 朝向」解析成状态名 |
| `AnimationStateSystem` | system | 接线：`ResolveState` → `SetState`，自身不做方向判断 |
| `AnimationSystem` | system | 按 `dt` 推进当前状态的帧 |

### 数据流

`FacingSystem` 写 `Facing{Direction, Walking}` → `AnimationStateSystem` 用 `DirectionalAnimation.ResolveState(walking, facing)` 得状态名 → `Animation.SetState(name)` 切换 → `AnimationSystem` 推进帧 → 渲染取 `Animation.Frame()`。

### 状态命名（单一来源）

状态名为 `{idle|walk}_{down|left|up|right}`，由 `component.AnimationStateName(walking, dir)` 统一产出——`DirectionalSheetSpec` 据此造状态、`ResolveState` 据此查状态，二者永不失配。

### 方向集与回退

角色声明一个 `DirectionSet`，`ResolveState` 保证解析结果必落在集合内：

| DirectionSet | 语义 |
|---|---|
| `DirectionsFour` | 上下左右各有独立素材，朝向直接采用 |
| `DirectionsHorizontal` | 仅左右素材；上下移动回退到上一次水平朝向（初始向右） |

### 工厂：单 spec 双产出

工厂用一个 `DirectionalSheetSpec`（方向集 + 每方向行首帧 `Rows` + `IdleFrames`/`WalkFrames` + `FPS`）同时产出动画状态集与方向集，二者同源不会失配。`States()` 为集合内每个方向生成 idle/walk 两态（共享该行首帧，idle 即首帧）。具体行布局以各工厂为准，sprite sheet 见对应 `sprite.source.md`。

```go
spec := component.DirectionalSheetSpec{
    Directions: component.DirectionsFour,
    Rows:       map[component.FacingDirection]int{ /* 方向 → 行首帧 */ },
    IdleFrames: 1, WalkFrames: 6, FPS: cfg.AnimFPS,
}
animations.Set(e, *component.NewAnimation(spec.States()))
directionals.Set(e, *component.NewDirectionalAnimation(spec.Directions))
```

### fail-fast

无静默兜底：`SetState` 遇未知状态、`States()` 遇方向缺 `Rows` 行都直接 panic，让配置错误在开发期暴露。正常配置下 `ResolveState` 只产出集合内状态、`States()` 必为其生成对应态，故 panic 不可达。

## 时间系统

**禁止依赖固定帧率**。所有游戏逻辑必须基于时间（秒），而非帧数。

### Clock API

| 方法 | 用途 |
|------|------|
| `NewClock()` | 创建根时钟 |
| `NewChildClock(parent)` | 创建子时钟，继承父时钟的deltaTime |
| `Update(rawDelta)` | **仅根时钟**：设置原始帧时间并推进TotalTime |
| `Tick()` | **仅子时钟**：推进TotalTime |
| `DeltaTime()` | 获取缩放后的帧时间 |
| `TotalTime()` | 获取累计游戏时间 |
| `SetScale(s)` | 设置时间缩放（负值钳位为0） |
| `Pause()` / `Resume()` | 暂停/恢复（保留暂停前的scale） |

### 根时钟 vs 子时钟

- **根时钟**：时间源头，用`Update(rawDelta)`接收帧时间
- **子时钟**：从父时钟继承deltaTime，用`Tick()`推进时间，scale与父时钟叠加

```go
root := component.NewClock()
game := component.NewChildClock(root)
ui := component.NewChildClock(root)

// 每帧更新
root.Update(1.0 / float64(ebiten.TPS()))  // 根时钟：传入帧时间
game.Tick()                                // 子时钟：用Tick()
ui.Tick()

// Scale叠加：game.DeltaTime() = root.DeltaTime() * 0.5
game.SetScale(0.5)  // 游戏慢动作，UI正常
```

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

### 时间缩放

| Scale值 | 效果 |
|---------|------|
| 0 | 暂停 |
| 0.5 | 慢动作 |
| 1.0 | 正常 |
| 2.0 | 快进 |

`Pause()`保存当前scale，`Resume()`恢复之前的scale：
```go
clock.SetScale(0.5)  // 慢动作
clock.Pause()        // 暂停，保存scale=0.5
clock.Resume()       // 恢复scale=0.5（而非默认1.0）
```

## 多语言（i18n）

所有用户可见文本（窗口标题、对话、UI文案）必须经`i18n.Bundle`取得，代码中不留硬编码字符串。单一`Bundle`实例由`main`通过`i18n.New()`创建，经`game.New`注入Scene，是文本唯一来源。locale数据文件的格式与新增语言规则见`internal/i18n/locales/AGENTS.md`。

### Bundle API

| 方法 | 说明 |
|------|------|
| `i18n.New()` | 加载内嵌locales，默认语言`DefaultLocale="zh"`（生产入口） |
| `i18n.Load(fsys, locale)` | 从自定义`fs.FS`加载（测试用） |
| `Get(key)` | 取当前语言文本，缺失按回退链处理 |
| `SetLocale(locale)` | 切换语言；未加载的语言返回false并保持当前 |
| `Locale()` / `Locales()` | 当前语言 / 可用语言（已排序） |

### 回退规则

`Get`三级回退：当前语言 → 默认语言（`zh`）→ key本身（让缺失翻译可见而非空白）。

