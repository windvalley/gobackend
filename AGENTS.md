<!-- harness:start -->
## Harness

本项目的本地门禁由 harness 管。
下面几条对所有任务生效，不得跳过：

- 开始任何任务前（包括分析、设计与评审）先读 `.harness/rules/core.md`，再按场景读同目录下的方法文档；读不到时报告，不凭印象继续。
- 先判断任务分级。
  小优化、小迭代，不改库表与接口契约，代码改动少于 500 行（以 `.harness/config.json` 的 `intent.plan.maxLines` 为准）的直接实施；其余先探索、逐个澄清、写方案，交出方案后停下；方案交付后用户表达开始实施的意思才算授权，最初提出需求或回答澄清问题都不算；授权后连续实施到完成，不再等用户逐步安排。
- 改变行为的实施默认先红后绿：先把验收标准写成测试并确认失败，再实现到全部通过；不删测试、不弱化断言、不跳过检查。
- 改动放在任务分支里；交付到主分支只有两条路：本地交付用 `harness merge --branch <任务分支>`，远程交付用 `harness submit --branch <任务分支>`，两者都要求验证证据与独立评审，门禁改动还需要人工确认。
  用哪一条看 `.harness/config.json` 的 `delivery.mode`：`remote` 时本地不合并，由平台在人工评审后合并，之后用 `harness sync` 把本地主分支快进到远程主分支。
- 合并后保留任务 Worktree，不要自行删除；用户要求清理时先运行 `harness worktree clean` 给出预览，用户确认后再加 `--yes`，细则见 `harness.md` 的“收尾”。

细则见 `.harness/rules/harness.md`。
<!-- harness:end -->
