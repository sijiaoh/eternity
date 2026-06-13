# LPC Goblin Sprite Sheet

## 来源

- **URL**: https://opengameart.org/content/lpc-goblin
- **下载文件**: https://opengameart.org/sites/default/files/goblinsword_0.png
- **作者**: Stephen "Redshrike" Challener (graphic artist)、William.Thompsonj (contributor)、Evert (sprite sheet 重排) (https://opengameart.org/users/redshrike)
- **许可证**: CC-BY 4.0 (Creative Commons Attribution 4.0)

## 许可证说明

CC-BY 4.0 允许自由使用、修改和商业化，但需要署名。
使用时需注明: `Attribute Stephen "Redshrike" Challener as graphic artist and William.Thompsonj as contributor. If reasonable link to this page or the OGA homepage.`（建议附带本页链接 https://opengameart.org/content/lpc-goblin ）

## 素材内容

- `sprite.png` - 持匕首哥布林 sprite sheet（LPC 风格，俯视 3/4 视角，真 4 方向动画）

## Sprite Sheet 布局

- 整图尺寸：**704 × 320 像素**
- **FrameWidth = 64, FrameHeight = 64, Columns = 11**（11 列 × 5 行，共 64×64 帧）
- 行 = 方向；每行 11 帧 = 行走(前 8 帧) + 匕首突刺攻击(后 3 帧)
- **真 4 方向素材**：每个方向（含 right）各占独立一行，无需翻转

| 行 | 帧区间 | 方向 | walk 帧 | attack 帧 |
|---|---|---|---|---|
| 0 | 0–10  | down（正面） | 0–7   | 8–10  |
| 1 | 11–21 | left         | 11–18 | 19–21 |
| 2 | 22–32 | up（背面）   | 22–29 | 30–32 |
| 3 | 33–43 | right        | 33–40 | 41–43 |
| 4 | 44–48 | —（death，不循环，5 帧） | — | — |

> 动画状态（idle/walk 的 StartFrame、FrameCount 等）以 `internal/entity/goblin_factory.go` 为准，此处仅描述素材本身。

### 渲染提示

角色本体约占 64×64 帧中部 ~32px，四周透明留白较多；按帧宽换算 SizeInUnits 时角色偏小，必要时调大 SizeInUnits 或调整 Anchor（脚底约在帧内 y≈58）。

## 原始 README

Attribute Stephen "Redshrike" Challener as graphic artist and William.Thompsonj as contributor. If reasonable link to this page or the OGA homepage. Sprite sheet realignment (rear-left-front-right + attack 独立成行) by Evert.
