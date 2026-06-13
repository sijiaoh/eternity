# 图像素材管理

免费素材（含版权记录与`.source.md`格式）规则见上层`assets/AGENTS.md`。

## 图像生成

使用AI生成图像时，创建`[素材名].prompt.md`文件记录生成参数。

### `.prompt.md` 格式要求

```markdown
# 素材名称

## 生成工具

- **API**: [使用的API，如 fal.ai FLUX dev]
- **模型**: [模型名称]

## Prompt

[生成时使用的prompt]

## 参数

- **尺寸**: [宽x高]
- **其他参数**: [如有]

## 素材内容

- `文件名.png` - 说明
```

## 整理

素材应该按照领域模型组织，而非按下载来源组织。
每个图片文件必须有同名的`.source.md`或`.prompt.md`（不允许目录单位的source文件）。

### 目录结构示例

```
assets/images/
└── characters/
    ├── hero/
    │   ├── sprite.png
    │   ├── sprite.source.md
    │   ├── portrait.png
    │   └── portrait.source.md
    └── enemy/
        ├── sprite.png
        ├── sprite.source.md
        ├── portrait.png
        └── portrait.source.md
```

### 命名规则

- 目录按领域模型命名：`characters/hero/`、`characters/enemy/`、`ui/buttons/`
- 图片按用途命名：`sprite.png`、`portrait.png`、`icon.png`