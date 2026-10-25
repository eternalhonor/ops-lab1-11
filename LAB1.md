# 实验 1 记录：第 11 组

实验日期：2026-10-25 至 2026-10-26。本文记录 v1.0.0、v1.0.1 及此前的实验过程；报告本身通过后续 PR 提交。时间统一为北京时间（UTC+8），精确时间从 GitHub API 导出，附页面链接供复核。CI 时长统一采用 ci 作业的结束时间减开始时间，不与工作流总时长混用。

## 1. 基本信息

| 项 | 内容 |
|---|---|
| 小组 | 第 11 组 |
| 仓库地址 | <https://github.com/eternalhonor/ops-lab1-11>（公开） |
| 代码载体 | 统一示例 enroll-conflict，Go |
| runner 与语言版本 | Actions：ubuntu-24.04、Go 1.26.9；ECS：Ubuntu 24.04.4 LTS、x86_64、Go 1.26.9 |
| 合并方式 | 仅允许 squash；禁止 merge commit 和 rebase merge；合并后自动删除分支 |
| 分支模型 | main 为唯一长期分支，改动经短期分支、PR、检查和其他成员审批后合入 |

## 2. 任务 0：本地基线

本阶段使用教师模板克隆 `/root/lab1-local`；后续实验使用小组仓库 `/root/ops-lab1-11`，两者的提交历史不可混用。

第一次以 `VERSION=1.0.0 BUILD_NO=1` 打包得到的 BUILD_INFO：

```text
name=enroll-conflict
version=1.0.0
build_no=1
run_id=local
commit=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
go=go1.26.9
```

修改 conflict.go 中一行注释但不提交，再用相同参数打包，diff 显示 BUILD_INFO 没有变化，二进制哈希发生变化：

| 构建 | 二进制 SHA-256 |
|---|---|
| 第一次，干净工作区 | `f23447d406c42af7dd35e23c2d1b2f15484cfb709cc9b4bb1262643694b8528e` |
| 第二次，有未提交注释改动 | `d1003950147da8d8f06b725d882de3687db2bdcb06acf37e11ec56a505f21d35` |

第一次 go version -m 的 VCS 信息：

```text
build   vcs=git
build   vcs.revision=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
build   vcs.time=2026-10-12T16:26:07Z
build   vcs.modified=false
```

第二次：

```text
build   vcs=git
build   vcs.revision=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
build   vcs.time=2026-10-12T16:26:07Z
build   vcs.modified=true
```

BUILD_INFO 能说明名称、版本、构建号、声明的提交号和 Go 版本，但未记录构建机器或构建时间；这里的 run_id=local 也不能定位 Actions 运行。vcs.time 是提交时间，不是构建时间。未提交修改存在时，commit 仍指向原提交；若把它理解为构建时完整源码状态，就会出错。此次注释不改变业务逻辑，但 Go 自动嵌入的 dirty 标志发生变化，因此二进制也发生变化；该标志只能提示存在修改，不能恢复具体差异。

## 3. 证据

| 项 | 链接或文本 |
|---|---|
| 规则集 | [protect-main，25826960](https://github.com/eternalhonor/ops-lab1-11/rules/25826960)：active；目标为默认分支 main；bypass_actors 为空 |
| 启用的规则 | Restrict deletions、Block force pushes、Require a pull request（至少 1 人批准）、Require status checks（ci，integration ID 15368） |
| 规则适用说明 | 组长直接推送也被拒绝；未启用“必须跟上最新主干”和“新提交后撤销旧批准”等额外选项。规则集内允许多种合并方式，但仓库级设置仅允许 squash，实际使用 squash |
| 第一轮 PR | [#1](https://github.com/eternalhonor/ops-lab1-11/pull/1)：新增双周班同一单节次冲突测试，并演练错误周次判断的失败与修复 |
| 第一轮：失败的 PR 运行 | [ci #2](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212235906)，TestConflicts 和 TestConflictEndpoint 发现单周与双周被错误判为冲突 |
| 第一轮：修复后通过 | [ci #3](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551) |
| 第一轮：合入 main 后的运行（含制品） | [ci #4](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212947752) |
| 第二轮 PR | [#2](https://github.com/eternalhonor/ops-lab1-11/pull/2) |
| 第二轮：缺陷和测试遗漏 | 将 Valid() 的非法星期判断从 Weekday > 7 改为 Weekday > 6，错误拒绝星期日；原测试检查了合法星期三及非法星期 0、8，却没检查合法星期 7，因此仍通过 |
| 第二轮：仅含缺陷时仍然通过 | [ci #5](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40214666542) |
| 第二轮：补测试后失败 | [ci #6](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40215797101)，TestValidSunday 失败 |
| 第二轮：撤销缺陷后通过 | [ci #7](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40216036223)，将 > 6 改回 > 7，保留测试 |
| 第二轮：合入 main 后通过 | [ci #8](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221353173) |
| 发布工作流 PR | [#3](https://github.com/eternalhonor/ops-lab1-11/pull/3)，[PR CI](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221652315)，[主干 CI](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221981064) |
| release gate 拒绝未合入的提交 | [release #2](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40223221802)：v0.0.1 指向 1ad35f18faadd579b17bebcd54a2d01e43fa0ca0，在门禁步骤失败 |
| 测试引用清理 | 远端 v0.0.1 和 scratch/gate-test、本地 v0.0.1 已删除；失败运行记录保留。门禁在创建 Release 步骤之前失败 |
| Release v1.0.0 | [v1.0.0](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.0) |
| Release v1.0.1 | [v1.0.1](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.1) |

组长直接推送 main 被拒的终端输出：

```text
[main b90b9e7] test: direct push
remote: error: GH013: Repository rule violations found for refs/heads/main.
remote: Review all repository rules at https://github.com/eternalhonor/ops-lab1-11/rules?ref=refs%2Fheads%2Fmain
remote:
remote: - Changes must be made through a pull request.
remote:
remote: - Required status check "ci" is expected.
remote:
To https://github.com/eternalhonor/ops-lab1-11
 ! [remote rejected] main -> main (push declined due to repository rule violations)
error: failed to push some refs to 'https://github.com/eternalhonor/ops-lab1-11'
```

发布门禁失败日志：

```text
Release blocked: tag commit is not in origin/main history.
Error: Process completed with exit code 1.
```

网络情况说明：ECS 的 Git HTTPS 连接多次超时。PR #2 的新增测试和修复通过网页提交，证据以远端记录为准，本地未成功推送的 4eaf384 不作为该 PR 的提交。v1.0.1 推送超时后，通过 API 确认远端无标签，核对标签提交与远端 main 一致，再通过 API 创建附注标签及引用，触发原有发布工作流，没有手工上传或替换附件。远端标签对象 `cf78994f44235552b1ceb771235683e84c6f093d` 指向提交 `5db1ee804e9641d6ea9433c4ccd0c4de3fab234e`；标签对象 SHA 与提交 SHA 不同是正常的。

## 4. 追溯与运行

| 项 | v1.0.0 | v1.0.1 |
|---|---|---|
| 版本号 | 1.0.0 | 1.0.1 |
| 构建号 | 1 | 3 |
| 提交号（SHA） | `364dc77e65a27384e45e139a19ca97d505184895` | `5db1ee804e9641d6ea9433c4ccd0c4de3fab234e` |
| 引入该提交的 PR | [#3](https://github.com/eternalhonor/ops-lab1-11/pull/3) | [#4](https://github.com/eternalhonor/ops-lab1-11/pull/4) |
| 构建它的 Actions 运行 | [release #1，40222522245](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40222522245) | [release #3，40225915528](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40225915528) |
| 制品 SHA-256（已与 Release 的 .sha256 校验一致） | `27d2fb4cb6c79d4763b48412cc64d946748ee5936e3cb24bc010b28e06218956` | `d2e6322b744ebfec650158315e09cbc09ab902c8c4e72385746788cf7eb5e1ad` |

构建号不能唯一确定运行，因为 ci 和 release 分别编号；v1.0.0 对应 [release #1](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40222522245)，而 [ci #1](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40204984418) 也是编号 1，但构建的是另一提交的主干制品。run_id 定位具体运行，run_attempt 区分同一运行的不同尝试。v1.0.1 为 release 构建号 3，是因为 release #2 被门禁拒绝，失败也占用编号。

发布 v1.0.1 后，重新下载 v1.0.0 到独立目录，与首次保存的哈希对照：

```text
enroll-conflict_1.0.0_linux_amd64.tar.gz: OK
27d2fb4cb6c79d4763b48412cc64d946748ee5936e3cb24bc010b28e06218956  enroll-conflict_1.0.0_linux_amd64.tar.gz
```

重新下载的旧版本与发布时记录的值一致。

trace.sh 输出（v1.0.0）：

```text
== 1. 校验和
enroll-conflict_1.0.0_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.0
build_no=1
run_id=40222522245
commit=364dc77e65a27384e45e139a19ca97d505184895
go=go1.26.9
== 3. 提交 364dc77e65a27384e45e139a19ca97d505184895
作者: chen osamu <C.amellia@sjtu.edu.cn>
时间: 2026-10-26T02:35:30+08:00
标题: ci: add gated release workflow (#3)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.0
版本号: 1.0.0    构建号: 1
提交页面（含关联的 PR）: https://github.com/eternalhonor/ops-lab1-11/commit/364dc77e65a27384e45e139a19ca97d505184895
构建它的运行: https://github.com/eternalhonor/ops-lab1-11/actions/runs/40222522245
```

trace.sh 输出（v1.0.1）：

```text
== 1. 校验和
enroll-conflict_1.0.1_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.1
build_no=3
run_id=40225915528
commit=5db1ee804e9641d6ea9433c4ccd0c4de3fab234e
go=go1.26.9
== 3. 提交 5db1ee804e9641d6ea9433c4ccd0c4de3fab234e
作者: eternalhonor <2749393080@qq.com>
时间: 2026-10-26T03:13:21+08:00
标题: test: cover all valid weekdays (#4)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.1
版本号: 1.0.1    构建号: 3
提交页面（含关联的 PR）: https://github.com/eternalhonor/ops-lab1-11/commit/5db1ee804e9641d6ea9433c4ccd0c4de3fab234e
构建它的运行: https://github.com/eternalhonor/ops-lab1-11/actions/runs/40225915528
```

两次脚本均完成校验和提交追溯，未报错。运行时解压已下载的 Release 压缩包，先启动 1.0.1 的二进制，监听 127.0.0.1:18080，查询 /version 后停止进程，再启动 1.0.0；没有重新编译，演练结束后两个服务均已停止。

```text
=== 当前运行版本：1.0.1 ===
{"version":"1.0.1","commit":"5db1ee804e9641d6ea9433c4ccd0c4de3fab234e","build_no":"3","go":"go1.26.9"}

=== 当前运行版本：1.0.0 ===
{"version":"1.0.0","commit":"364dc77e65a27384e45e139a19ca97d505184895","build_no":"1","go":"go1.26.9"}
```

## 5. 数据

### 第一轮 gate 演练时间线

以下均为 2026-10-26，北京时间（UTC+8）；检查结束以 ci 作业完成时间为准。

| 事件 | 时间 | 链接 |
|---|---|---|
| 发起 PR #1 | 00:33:15 | [PR #1](https://github.com/eternalhonor/ops-lab1-11/pull/1) |
| 第一次 CI 运行创建 | 00:33:18 | [ci #2](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212235906) |
| 第一次 ci 开始 | 00:33:20 | [作业](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212235906/job/121225736036) |
| 第一次 ci 失败，作业结束 | 00:33:34 | [失败作业](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212235906/job/121225736036) |
| 创建修复提交 | 00:35:58 | [0eb3a5e](https://github.com/eternalhonor/ops-lab1-11/commit/0eb3a5ea2a159c800d9282ed6d34ed3ea7cd4333) |
| 推送修复提交：以新 CI 创建时间作为接收时间代理 | 00:36:04 | [ci #3](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551) |
| 修复后的 ci 开始 | 00:36:06 | [作业](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551/job/121226425652) |
| ci 通过，作业结束 | 00:36:20 | [通过作业](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551/job/121226425652) |
| ChenJewel 批准 | 00:38:08 | [审批记录](https://github.com/eternalhonor/ops-lab1-11/pull/1#pullrequestreview-5644161720) |
| 合入 main | 00:42:14 | [PR #1](https://github.com/eternalhonor/ops-lab1-11/pull/1) |
| main 的 ci 通过 | 00:42:32 | [ci #4](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212947752) |

没有单独记录客户端推送开始和 GitHub 接收 push 的精确时间。提交对象时间是 00:35:58，随后运行创建于 00:36:04；两者不能当作同一事件。下文用运行创建时间作可观察的分界，并明确其局限。

从首次失败结束到修复后的检查通过，共 **166 秒（2 分 46 秒）**：失败结束至修复提交为 144 秒，提交至新运行创建为 6 秒，新运行创建至作业开始为 2 秒，作业执行为 14 秒。按运行创建时间划分，前 150 秒是人工处理、提交及传输/事件延迟的合计，后 16 秒是排队与检查；不能声称前 150 秒全是编码时间。PR 创建到合并共 539 秒，修复检查通过到合并另有 354 秒，用于审批和后续操作等待。

### 提交到发布的间隔

以下均为 2026-10-26、UTC+8；起点取 PR 第一个实际提交的提交者时间，终点取首次包含该改动的 Release 的 published_at，不是 PR 创建时间或 tag 时间。

| PR | 第一个提交的时间 | 合入时间 | 对应 Release 的发布时间 | 提交到发布的间隔 |
|---|---|---|---|---|
| [#1](https://github.com/eternalhonor/ops-lab1-11/pull/1) | 00:33:11 | 00:42:14 | [v1.0.0](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.0)，02:42:48 | 2 小时 9 分 37 秒（7777 秒） |
| [#2](https://github.com/eternalhonor/ops-lab1-11/pull/2) | 01:04:00 | 02:27:43 | [v1.0.0](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.0)，02:42:48 | 1 小时 38 分 48 秒（5928 秒） |
| [#3](https://github.com/eternalhonor/ops-lab1-11/pull/3) | 02:31:00 | 02:35:30 | [v1.0.0](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.0)，02:42:48 | 11 分 48 秒（708 秒） |
| [#4](https://github.com/eternalhonor/ops-lab1-11/pull/4) | 03:11:24 | 03:13:21 | [v1.0.1](https://github.com/eternalhonor/ops-lab1-11/releases/tag/v1.0.1)，03:25:12 | 13 分 48 秒（828 秒） |

### CI 时长

所有 12 次初始运行均列出，另列重跑；口径为 ci 作业执行时长，不包含该作业启动前的排队。API 原时间为 UTC，下表已转为 UTC+8。

| 运行 | 事件 | 开始时间 | 结束时间 | 时长 | 备注 |
|---|---|---|---|---|---|
| [ci #1，attempt 1](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40204984418/attempts/1) | push | 10-25 23:01:47 | 10-25 23:02:34 | 47 秒 | 初始 CI，成功 |
| [ci #2](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212235906) | PR | 10-26 00:33:20 | 10-26 00:33:34 | 14 秒 | 第一轮，失败 |
| [ci #3](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551) | PR | 10-26 00:36:06 | 10-26 00:36:20 | 14 秒 | 第一轮，修复通过 |
| [ci #4](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212947752) | push | 10-26 00:42:18 | 10-26 00:42:32 | 14 秒 | PR #1 合并后，成功 |
| [ci #5](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40214666542) | PR | 10-26 01:04:11 | 10-26 01:04:29 | 18 秒 | 第二轮，缺陷未被发现，成功 |
| [ci #6](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40215797101) | PR | 10-26 01:18:09 | 10-26 01:18:25 | 16 秒 | 第二轮，补测试后失败 |
| [ci #7](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40216036223) | PR | 10-26 01:21:08 | 10-26 01:21:24 | 16 秒 | 第二轮，修复通过 |
| [ci #8](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221353173) | push | 10-26 02:27:48 | 10-26 02:28:01 | 13 秒 | PR #2 合并后，成功 |
| [ci #9](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221652315) | PR | 10-26 02:31:29 | 10-26 02:31:39 | 10 秒 | 发布工作流 PR，成功 |
| [ci #10](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40221981064) | push | 10-26 02:35:35 | 10-26 02:35:54 | 19 秒 | PR #3 合并后，成功 |
| [ci #11](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40224842316) | PR | 10-26 03:11:34 | 10-26 03:11:51 | 17 秒 | 所有合法星期测试 PR，成功 |
| [ci #12](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40224993404) | push | 10-26 03:13:26 | 10-26 03:13:42 | 16 秒 | PR #4 合并后，成功 |
| [ci #1，attempt 2](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40204984418/attempts/2) | push 事件的重跑 | 10-26 04:44:54 | 10-26 04:45:09 | 15 秒 | ChenJewel 发起；成功；页面工作流总时长 20 秒，与本列口径不同 |

12 次首次运行的平均作业时长为 **17.83 秒**，中位数 **16 秒**，范围 **10～47 秒**；包含重跑的 13 次平均为 **17.62 秒**。PR 检查与主干打包步骤不同，成功和失败运行执行的步骤也不同；共享 runner、缓存等可能影响时长，本实验未控制这些因素，不能把差异都归因于缓存，也不能以小样本代表团队的长期表现。

## 6. 分析问题

### 问题 1：gate 的收益与成本

所用数据：第 3 节两轮 PR 运行和发布门禁记录，第 5 节第一轮时间线及全部 CI 时长；第一轮失败后恢复通过为 166 秒，12 次初始 CI 作业平均 17.83 秒，PR #1 提交到发布为 7777 秒。

第一轮中必需检查拦住了单周与双周被错误判为冲突的提交，第二轮却让错误拒绝星期日的代码得到绿色结果，说明门禁的有效性依赖测试覆盖而不等于程序正确性。补充星期日测试后第二轮也被拦截，人工审批要求则在两轮检查通过但尚未批准时继续阻止合并，发布门禁另外拒绝了非主干历史上的提交。以“首次失败结束到修复后检查通过”为本次可观察的恢复成本，额外经历了 166 秒，其中新 CI 排队和执行为 16 秒，前面的 150 秒含人工处理、提交及传输延迟，不能精确等同于纯编码时间，也不是严格的无缺陷反事实成本。第一轮失败检查本身执行了 14 秒，而修复检查通过后到合并还等待了 354 秒，可见自动检查并非整个等待过程的唯一来源。平均一次 CI 作业约 17.83 秒，占 PR #1 的 7777 秒提交到发布时间约 0.23%，这里只比较单次作业与整个间隔，不把它当成该 PR 所有 CI 的总成本。补测仍无法保证发现未覆盖的业务组合、错误需求、并发和真实环境问题，人工审查也不能保证完全发现这些问题。上述收益成立的条件是规则保持启用、无人绕过、测试确实表达正确需求且失败能阻止合并；本组样本少，并受演练、审批安排和网络超时影响，不能推广为生产团队的效率结论。

### 问题 2：PR 中的构建与 main 上的构建

所用数据：[PR #1 修复通过的 ci #3](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40212457551) 中实际打印的值，以及在本地对 `git rev-list origin/main` 的精确匹配结果（采集时主干为 5db1ee8…）：

```text
github.sha=701a7ec1eab0795fbec4ef7c0c20d73e092fa78c
在当前 origin/main 历史中：False
pull_request.head.sha=0eb3a5ea2a159c800d9282ed6d34ed3ea7cd4333
在当前 origin/main 历史中：False
```

| 对象 | SHA |
|---|---|
| PR 检查使用的临时合并提交 | `701a7ec1eab0795fbec4ef7c0c20d73e092fa78c` |
| PR 分支头提交 | `0eb3a5ea2a159c800d9282ed6d34ed3ea7cd4333` |
| squash 后 main 上的提交 | [`9f2849f303eeb62611606e19c07af24a56e5dfe5`](https://github.com/eternalhonor/ops-lab1-11/commit/9f2849f303eeb62611606e19c07af24a56e5dfe5) |

PR 运行中的两个 SHA 实测不同：github.sha 对应 GitHub 为 PR 检查生成的临时合并提交，而 pull_request.head.sha 对应源分支头提交。默认 checkout 在本次 pull_request 工作流中检查合并结果，因此有利于验证变更与目标分支的组合。本组采用 squash 合并，在 main 上重新生成了 9f2849f…，所以原分支头提交与测试合并提交都没有成为当前主干的祖先，这与只比较文件内容不同。若使用其他合并方式，源分支提交是否保留在主干历史中可能不同，不能把本组结论推广到所有合并方式。PR 检查只能证明当时测试用的候选状态通过检查，不能直接把其制品当成已批准并合入的正式制品。本实验因此仅在 push main 时产出主干包，并在位于 main 历史中的 tag 上构建正式 Release，使包内提交号对应已接受的主干历史与 PR；单独创建任意 tag 还不够，必须通过发布门禁。

### 问题 3：可重复构建

本地五次构建均使用提交 `5db1ee804e9641d6ea9433c4ccd0c4de3fab234e`、Go 1.26.9、VERSION=1.0.1、RUN_ID=local-repro，并显式传入该 COMMIT。①、②、④、⑤ 的 BUILD_NO=900，③ 为 901。④、⑤ 临时删除 Makefile 中 tar 的 --mtime=@0，两次间隔 3 秒；完成后恢复 Makefile，git status -sb 显示干净且与主干同步，没有提交实验改动。

| 构建 | 条件 | 二进制 SHA-256 | tar.gz SHA-256 |
|---|---|---|---|
| ① | 参数固定，工作区干净 | `a4f0cb9f2c849b9c2f65eb5748d415fc7640a7c85bdba91e39ef39effc372aba` | `83179483bc285985fd0c792dfb63ea03efe842bb1136461aad531bf1ef0c7be8` |
| ② | 与①相同 | `a4f0cb9f2c849b9c2f65eb5748d415fc7640a7c85bdba91e39ef39effc372aba` | `83179483bc285985fd0c792dfb63ea03efe842bb1136461aad531bf1ef0c7be8` |
| ③ | 仅 BUILD_NO 改为 901 | `6f7796311289f617d9f0f95304e511420db41452a69fbb2614763c0e2a075545` | `0347f5cd2eb2fe79d3c7c83b7761b46cdf1aa53e51e96615372cc97e5c47c987` |
| ④ | 去掉 --mtime=@0，其余同① | `29d0c654547d7ff48de19462401af1dea5521ae5d9f118eb29ff6d3086cb29d9` | `bd48cec9ef075f722cb4fe318431bbfd5be6be4a48891e9206d6c74c3e56529c` |
| ⑤ | 同④，间隔 3 秒 | `29d0c654547d7ff48de19462401af1dea5521ae5d9f118eb29ff6d3086cb29d9` | `9ce389bbcf670b22324a9752d00d13c0550465aaf54e34a189330d6b3b641081` |
| Actions：ci #1 attempt 1 | main push，原始尝试 | `3453c5480252b3639c93187f956dcacb56ef175919858028934f95b344f2af0c` | `0672a7b8a684bcf56452b5537071b79b7f09fd83abf126168959b357a68ef5db` |
| Actions：ci #1 attempt 2 | 对同一次运行 Re-run all jobs | `3453c5480252b3639c93187f956dcacb56ef175919858028934f95b344f2af0c` | `0672a7b8a684bcf56452b5537071b79b7f09fd83abf126168959b357a68ef5db` |

Actions 数据来自 [ci #1 第一次尝试](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40204984418/attempts/1)和[第二次尝试](https://github.com/eternalhonor/ops-lab1-11/actions/runs/40204984418/attempts/2)的原始日志压缩包，读取 package 步骤的哈希，不是上传 artifact ZIP 的哈希。它们的源码为 `1f385585b270f97de8d6034630ee1a5a5c99e444`，版本为 0.0.0-main、构建号为 1、run_id 为 40204984418，与本地五次构建的提交和参数不同，不能直接跨组比较哈希。

本地①的 VCS 信息：

```text
build   vcs=git
build   vcs.revision=5db1ee804e9641d6ea9433c4ccd0c4de3fab234e
build   vcs.time=2026-10-25T19:13:21Z
build   vcs.modified=false
```

本地④的 VCS 信息：

```text
build   vcs=git
build   vcs.revision=5db1ee804e9641d6ea9433c4ccd0c4de3fab234e
build   vcs.time=2026-10-25T19:13:21Z
build   vcs.modified=true
```

①、②的二进制及压缩包均一致，说明在本次固定源码、参数和环境条件下构建可重复，而③改变了写入二进制及 BUILD_INFO 的构建号，导致两种哈希都改变。④的二进制与①不同，结合 vcs.modified 从 false 变为 true 可知，临时修改 Makefile 使源码树变脏并影响 Go 的嵌入元数据，不能把这项差异归因于 tar 参数直接改变了程序。④、⑤的二进制一致，但压缩包不同，是因为去掉时间归一化后重新生成的 BUILD_INFO 等文件的修改时间进入了 tar，即使内容一致也会改变归档字节。Actions 两次尝试的源码、版本、构建号和 run_id 不变，实测两种哈希也相同，虽然尝试编号和 artifact 名称变化，但尝试编号未写入包内容。制品应保留名称、版本、准确源码身份、工具链及追溯标识，正式构建应拒绝或完整描述 dirty 状态；机器绝对路径、当前墙钟时间和未归一化文件时间等非必要易变信息应避免进入可重复比较的内容，实际构建机器与时间可放在外部流水线记录中。构建号和 run_id 有追溯价值，但它们也是构建输入，若跨运行改变这些值，就不能同时要求整个包保持字节一致；固定 Go 补丁版本也不足以永久固定 runner 镜像和 Actions 主版本引用。release 在同一提交上重编译，是为了将正式版本号从 tag 注入二进制，而主干 CI 包使用 0.0.0-main，这确实偏离严格的“一次构建、多处部署”，代价是额外构建时间及环境漂移风险，因此应验证正式制品，并在部署和回退时复用该正式制品；本次重跑一致也不能保证任意未来环境下重编译都一致。

### 问题 4：发布与回退

所用数据：第 4 节两次 /version 输出、原始及重新下载的 v1.0.0 哈希；问题 3 的参数与环境对哈希的影响。机制参考：[GitHub Immutable releases](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases)。

若 v1.0.1 上线后发现缺陷，应先评估影响并切回保存且验证过的 v1.0.0 制品，再通过 PR 修复、审批和测试后发布新版本号，例如 v1.0.2，不能删除旧 tag 后复用 v1.0.1，否则同名版本会对应不同内容，使缓存、审计和追溯失去确定性。本次回退使用 enroll-conflict_1.0.0_linux_amd64.tar.gz 解压出的二进制，没有重新构建，接口从版本 1.0.1／提交 5db1ee8…／构建号 3，切换为版本 1.0.0／提交 364dc77…／构建号 1，直接验证了实际运行对象的变化。重新从 v1.0.0 源码编译不保证得到原附件，因为构建号、run_id、dirty 状态、工具链及归档元数据等都可能影响字节，即使本实验在固定条件下曾得到一致结果也不足以保证未来重建一致。若攻击者或误操作同时替换压缩包和配套的 .sha256，并保留可通过检查的 BUILD_INFO，trace.sh 的配对校验和源码追溯可能仍通过，因此该脚本不能单独证明二进制确由声明源码构建，也不能把两份可同时修改的文件当成可信签名。将新下载文件与本组预先独立保存的 27d2fb4…18956 比较，则能发现字节变化，前提是这份基准记录本身未被一并篡改。GitHub 的 Immutable releases 可锁定已发布的附件和关联 tag，并提供 release attestation，加强对发布后替换的防护，但本文没有把这一设置当作本组已经启用并验证的控制。最后，这次是本机无数据库迁移的最小回退演示，生产回退还需满足数据格式和配置兼容性，并检查健康状态，不能仅以版本接口改变断言业务已完全恢复。

## 7. 选做记录

| 选做 | 完成情况与证据 |
|---|---|
| O1 容器镜像 | 未做 |
| O2 缓存对反馈时间的影响 | 未做；本报告的一次 CI 重跑用于问题 3 的可重复构建对比，不作为缓存对照实验 |
| O3 迁移到小组项目 | 未做 |

## 8. 成员贡献

以下为报告提交前已核实的功能与流水线 PR；每个 PR 都由不同于发起人的成员批准。

| GitHub 用户名 | 发起并合入的 PR 编号 | 审批过的 PR 编号 |
|---|---|---|
| eternalhonor | [#1](https://github.com/eternalhonor/ops-lab1-11/pull/1)、[#4](https://github.com/eternalhonor/ops-lab1-11/pull/4) | [#2](https://github.com/eternalhonor/ops-lab1-11/pull/2#pullrequestreview-5644647122)、[#3](https://github.com/eternalhonor/ops-lab1-11/pull/3#pullrequestreview-5644670998) |
| HleTopZeus | [#2](https://github.com/eternalhonor/ops-lab1-11/pull/2) | 无（本次已采集的 #1～#4 记录中） |
| ChenJewel | [#3](https://github.com/eternalhonor/ops-lab1-11/pull/3) | [#1](https://github.com/eternalhonor/ops-lab1-11/pull/1#pullrequestreview-5644161720)、[#4](https://github.com/eternalhonor/ops-lab1-11/pull/4#pullrequestreview-5644817573) |

证据保留：仓库保持公开，PR、Release 和 Actions 记录保留至检查结束且成绩公布。ECS 的 `/root/lab1-evidence/` 保存原始日志、哈希、追溯和版本输出；正式附件另存于 `/root/ops-lab1-11/dist/release/`。释放 ECS 前需要下载证据备份，不要执行会删除该附件目录的 make clean。
