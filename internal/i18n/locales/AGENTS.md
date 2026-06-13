# locale 数据文件

游戏的翻译源数据（一方数据，非`assets/`素材），经`//go:embed`内嵌进二进制。

- 每种语言一份`<locale>.json`，默认语言为`zh`。
- 格式：扁平JSON对象`key→文本`，key用点号命名空间（如`window.title`、`dialogue.intro.ready`）。

## 新增语言

增一份`<locale>.json`，**补全默认语言的所有key**。`TestEmbeddedLocalesComplete`双向校验各locale的key与默认语言完全一致——漏译或多余key都会fail。

非默认语言漏译会静默回退到默认语言而非报错，故须靠该完整性测试锁死。
