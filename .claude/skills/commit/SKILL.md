---
description: 查看当前变更并提交代码。使用时机：用户要求提交、commit、或保存变更时。
disable-model-invocation: true
allowed-tools: Bash(git *)
---

# Commit

查看当前变更并提交。

1. 运行 `git status` 查看变更文件
2. 运行 `git diff` 查看具体变更内容
3. 运行 `git log -3 --oneline` 查看最近提交风格
4. 根据变更内容撰写简洁的提交信息
5. 使用 `git add` 添加相关文件
6. 使用 `git commit` 提交

## 注意

- 不要修改任何代码
- 不要使用 `--amend`
- 不要使用 `git add -A`
