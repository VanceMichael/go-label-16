# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

灾后演练导入一份带“ Barn ”和“Source”标签的饲喂计划后，调用方手里的原始标签立刻被改成了小写并去掉空格；恢复人员又在已恢复记录上加本地备注，校验原快照时发现快照标签也跟着变化，哈希随之不再匹配。请先不要修改代码，查清标签从构建快照到生成恢复结果期间为什么没有保持隔离，说明两次污染各自经过的对象边界。

## 含 Bug 版本

- 仓库：VanceMichael/go-label-16
- 仓库地址：https://github.com/VanceMichael/go-label-16.git
- parent SHA：4bfe59160fa1ef29af8ff07cafef4b0054ffedbf

## 复现步骤

```bash
git clone -- https://github.com/VanceMichael/go-label-16.git bug-repro
cd bug-repro
git checkout --detach 4bfe59160fa1ef29af8ff07cafef4b0054ffedbf
go test ./internal/recovery -run ^TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore$ -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/recovery -run ^TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore$ -count=1
--- FAIL: TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore (0.00s)
    labels_test.go:21: snapshot build changed caller labels: map[barn:north-7 source:mobile]
FAIL
FAIL	go-base/internal/recovery	0.032s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/recovery -run ^TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore$ -count=1
--- FAIL: TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore (0.00s)
    labels_test.go:21: snapshot build changed caller labels: map[barn:north-7 source:mobile]
FAIL
FAIL	go-base/internal/recovery	0.001s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

根因结论必须分别定位 Build 调用 normalizeLabels 前后的调用方 map 共享，以及 Apply 经 cloneRecord 生成恢复结果时仍与 Snapshot.Records 共享 Labels 的两段失效机制，并解释第二段修改为何会改变 Hash、使 Verify 失败；同时指出 Payload 已有复制不属于异常范围。以 TestSnapshotLabelsRemainIsolatedAcrossBuildAndRestore 的红测复核，目标仓库代码、测试和配置零改动，不得实施修复或只写成“map 是引用类型”。
