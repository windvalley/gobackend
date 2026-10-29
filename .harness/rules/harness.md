# Harness 门禁

本项目的本地门禁由 harness 提供。
所有判定在本机完成，结果写入账本，不依赖云端。
本文的流程与命令由 Agent 执行，标明必须由人的步骤除外。

## 接入项目

- 先阅读 harness 仓库的 README 与参考手册，调查项目规则、Git 状态、已有钩子和真实检查命令，并运行 `init --dry-run` 预览。
- 按项目约束选择接入方式，准备任务分支与需要的 Worktree，不要求用户预先创建。
  接入默认走受控流程，是否新建 Linked Worktree 取决于隔离需要与项目规则。
- 保留已有有效规则和检查，把起始检查对应到项目真实命令，明确占位与未执行项。
  接管原有工具或放弃某条规则需要用户决定时，先列出影响与选择。
- 写入并提交接入改动，执行验证和独立评审，展示 `harness confirm --branch <接入分支>` 的摘要。
  带 `--yes` 的确认由人在自己的终端执行，自然语言同意不能替代，Agent 不代为执行。
- 人完成确认后，继续完成已获授权的本地合并，列出当前客户端所需的重启、信任或注册操作。
  用户完成后再运行诊断，区分配置就位和当前客户端实际调用已验证，不把初始化成功作为全部接入完成。
- 合并后保留接入分支和任务 Worktree，清理遵循“收尾”。

## 一次任务

先按 `core.md` 判断任务分级：直接实施的任务从第 1 步开始；先方案后实施的任务，方案经用户授权后再从第 1 步开始。

1. 按项目规则准备任务分支和工作环境，需要 Worktree 时由 Agent 创建或复用符合隔离要求的目录。
   在任务分支里实现并提交，提交信息按项目规范，`commit-msg` 钩子会校验。
2. `harness verify --branch <任务分支>`：在固定检出目录里按改动路径选择检查集，结果写成绑定该分支提交的签名证据。
   主分支自上次验证后新增的改动若未触及某项检查，且这项检查的定义没有变，它直接复用。
3. `harness review --branch <任务分支>`：独立评审只读验收清单、证据摘要、提交说明、文件统计和 diff，分需求符合度与工程质量两个维度。
   有 high 或 critical 级别的发现会被拒绝。
   提交说明要如实写清改了什么，评审会核对它与 diff 是否一致。
4. 在已有合并授权的范围内执行 `harness merge --branch <任务分支>`：核对证据与评审，合并到主分支。
   改动涉及门禁配置时必须先由人确认。
   主分支在评审之后改了本分支也改过的文件时，合并以 `review_stale` 拒绝：在任务分支上合并主分支并提交，再重新 verify、review。
5. 合并后任务 Worktree 保留，清理按下面的“收尾”由用户发起。
6. 汇报结果入口、逐条验收证据、验证与评审结论，以及已合并、待授权合并或阻塞的最终状态。
   同时说明任务分支、提交和工作区位置，未执行项不能算作通过。

## 远程交付

`.harness/config.json` 的 `delivery.mode` 为 `remote` 时，主分支在平台上，合并由团队走 PR/MR 与平台的保护规则。
交付路径与本地交付不同：

1. 第 1 到第 3 步与本地交付相同：在任务分支实现、提交、`harness verify`、`harness review`。
2. `harness submit --branch <任务分支>`：核对与 `merge` 相同的条件，把核对过的那个提交普通推送到远程同名分支，并创建或更新指向主分支的 PR/MR。
   描述里 harness 只维护 `<!-- harness:submit:start -->` 与 `<!-- harness:submit:end -->` 之间的区块，团队写在区块外的说明不会被覆盖。
   提交记录分“已推送”与“PR/MR 已就绪”两个阶段；推送之后创建 PR/MR 失败时重跑只会补 PR/MR，不会重复推送。
3. 之后由人在平台上合并。harness 不在平台上合并，也不强推；远程同名分支已分叉时以 `push_rejected` 拒绝。
4. 平台合入后运行 `harness sync`：把本地主分支快进到远程主分支，同步主 Worktree 的文件，需要时刷新本机安装并执行 `merge.postMergeCommand`。
   远程带来的普通代码改动按“信任平台主分支”导入，不在本机重复验证与评审；远程带来的门禁改动、交付设置变化或首次接入仍然要人确认，确认命令是 `harness confirm --branch <远程>/<主分支>`。
5. 汇报时区分“已提交到平台”“平台已合入”“本地已同步”，不把提交成功说成已合并。

- `harness submit` 要求本地主分支等于远程主分支，不一致时先 `harness sync`，否则以 `base_stale` 拒绝。
- `harness submit` 只接受本地任务分支，不接受主分支、远程跟踪引用或提交号。
- 平台命令（`gh`、`glab`）的安装与登录由人完成，harness 不读取、不保存也不传递平台令牌。
- 主分支还是本地交付时，把 `delivery.mode` 改成 `remote` 的分支经人工确认后可以 `submit`，这就是开启开关的路径；关闭时走普通提交流程，平台合入后 `sync` 导入。
- 本地交付下同样可以用 `harness sync` 把主分支快进到远程，取代由人用绕过开关手动快进的做法。

## 收尾：任务 Worktree 由用户发起清理

合并后保留任务 Worktree 和任务分支，不要自行删除：会话常常还停在里面，合并后也可能要回来补改，删除时还会连带删掉被忽略的本机文件。

- `harness merge` 的输出写明任务 Worktree 保留在哪里，以及已合并但还没清理的 Worktree 有几个；汇报时把这两点转告用户。
- 用户要求清理时，先运行 `harness worktree clean` 把预览交给用户：哪些可以删、会连带删掉哪些被忽略的文件，哪些跳过以及原因。
- 用户确认后运行 `harness worktree clean --yes`。
  预览里有含 `.env` 这类本机文件的 Worktree 时，用户明确同意删除这些文件才加 `--with-local-files`。
- 不要绕开这个命令自己删除：`git worktree remove --force` 与 `git branch -D` 会被拦截，它们会丢掉未提交的改动或未合并的提交。
- 不在会话里删除会话所在的 Worktree，`harness worktree clean` 也会跳过命令所在的 Worktree。

## 验收清单

放在 `.harness/acceptance/<任务分支名，/ 换成 ->.md`，按 `.harness/templates/acceptance.md` 写，每条写成能用测试或命令验证的结果。

- 清单可以随任务分支提交，评审会标注来源为“任务分支，尚未经人工确认”。
- 清单也可以先合入主分支，评审标注来源为“主分支已确认”。
- 没有清单时，需求符合度只能记为 `unknown`，评审只评工程质量。

## 实施授权

`.harness/policy.json` 的 `presets` 里有 `plan-first` 时，“先方案后实施”由门禁兜底：

- 任务分支写了实施方案（默认 `docs/changes/*/plan.md`，包括还没提交的）而还没有授权记录时，Agent 只能写 `docs/changes/` 与 `.harness/acceptance/` 下的文件，写测试和代码会被拦截。
- 交出方案时，把验收清单、行为说明与实施方案提交到任务分支，向用户说明方案后停下。
- 方案交付后，用户在对话里明确同意实施该方案，才是这一阶段的实施授权，不要求固定措辞。
  最初提出需求或仅认可某个澄清答案，不能替代对已交付方案的实施授权。
  Agent 运行 `harness approve --branch <任务分支> --yes --note "<用户原话>"` 记录下来，然后连续实施到完成，不需要用户执行任何命令。
- `--note` 必须照录用户表示同意的原话，合并、`status` 与评审都会显示它；用户没有同意时不得记录，也不得改写或编造原话。
- 授权绑定验收清单的内容：之后要改清单，先向用户说明改动并取得同意，再重新记录；实施方案里的进度记录可以随时更新。
- 合并时，写了实施方案、改了库表相关文件（`intent.plan.triggers`）或代码改动超过 `intent.plan.maxLines` 行的分支，必须有授权记录；没写方案就直接做了大改动的分支会在这里被拒绝。

## 被拦截时

拦截信息里带 `正确做法`，按它改道，不要尝试绕过。
实际生效的规则以 `.harness/policy.json` 为准，`harness rules` 可以列出全部规则。
预设规则的常见项：

- 密钥文件：只读同名 `.example`。
- 不可逆命令：列出具体对象请用户确认。
- 主分支：本地交付只通过 `harness merge` 更新；远程交付不本地合并，用 `harness submit` 提交、`harness sync` 同步。

项目自己加的规则（例如生成文件、迁移基线）同样按拦截信息里的正确做法处理。

## 门禁改动

`.harness/config.json`、`policy.json`、`verify.json`、`rules/**` 允许直接修改，改动会被记录并在合并前要求人工确认（`harness confirm --branch <分支>`）。
`.harness/bin/**` 是执行器入口，`.husky/**`、`.githooks/**`、`.claude/settings.json`、`.codex/hooks.json` 是客户端与 git 钩子的入口。
它们只由 `harness install`（或 `init`、`upgrade`）生成，Agent 不能手改，但可以在任务分支运行这些命令重新生成；改动同样在合并前要求人工确认。

## 多个 Worktree

入口 `.harness/bin/harness` 随仓库提交，Linked Worktree 检出后即可用。
git 钩子调用主 Worktree 的入口，任务分支改不了门禁代码自己的行为。
`.harness/config.json` 配置了 `git.hooksDir` 时，钩子整份生成在那个目录，`core.hooksPath` 指向主 Worktree 里的这一份，所有 Worktree 执行同一份钩子。
项目自己的检查写在 `git.projectHooks` 目录下的同名可执行脚本里，harness 检查通过后调用，Agent 可以直接维护；生成的钩子不手改，改配置后运行 `harness install` 重新生成。
`harness merge` 在主 Worktree 里更新主分支，要求它停在主分支、洁净、没有进行中的 git 操作；`harness status` 与 `harness doctor` 会提前列出不满足的项。
`.harness/policy.json` 的 `presets` 里有 `main-worktree` 时，主 Worktree 只用来合并：

- Agent 不能改主 Worktree 里被 git 跟踪或未忽略的文件，被 `.gitignore` 忽略的本机文件不受限。
- Agent 不能在主 Worktree 里切离主分支。
- 任何人都不能在主 Worktree 里提交。
- 新任务用 `git worktree add -b <任务分支> <路径> <主分支>` 建 Linked Worktree，在那里修改、提交。
  会话停在主 Worktree 时，改 Linked Worktree 里的文件按那个 Worktree 判定；命令按段跟着前面的 `cd` 与 `git -C` 推算所在目录。
  目录无法确定的段（`cd "$变量"`、`cd -`、管道里的 `cd`、`GIT_DIR`）按最严格处理，要在别的 Worktree 执行就写明绝对路径。

## 例外上报

需要人决定时用 `harness escalate open` 提交，内容包含现象、已尝试、阻塞点和至少两个选项。
人用 `harness escalate answer --id <编号> --choice <选项>` 答复，续跑前用 `harness escalate resume --id <编号>` 读取决定。

## 度量

`harness metrics` 输出按周统计：合并次数、验证首过率、复用检查数、评审拒绝率、策略拦截、门禁确认和黄金任务通过率。
`harness conformance` 逐项判定 L3 判据，样本不足的项标 `PEND`，不计为完成。
