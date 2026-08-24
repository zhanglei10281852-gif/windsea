# Bug Reproduction

## 包的性质

当前 tasks/ws-001-campaign-audit/green 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。最终 baseline 分支已在下面固定的 parent SHA 上单独提交任务测试，可直接复核 red；tasks/ws-001-campaign-audit/green 包含同一测试并应得到 green。完整验证日志仍只在本地留存。

## 问题现象

Fix the offshore wind campaign publish flow: when audit persistence fails, return the error and keep the reviewed campaign in review instead of leaving a published campaign without its audit record.

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/windsea
- 仓库地址：https://github.com/zhanglei10281852-gif/windsea.git
- parent SHA：0c5c547e2a521e477822ece3bdca1d296aeb4438

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/windsea.git bug-repro
cd bug-repro
git checkout --detach 0c5c547e2a521e477822ece3bdca1d296aeb4438
go test ./internal/service -run TestCampaignPublishKeepsReviewWhenAuditFails -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/service -run TestCampaignPublishKeepsReviewWhenAuditFails -count=1
--- FAIL: TestCampaignPublishKeepsReviewWhenAuditFails (0.01s)
    campaign_audit_atomicity_task_test.go:33: campaign state=published
FAIL
FAIL	github.com/zhanglei10281852-gif/windsea/internal/service	0.016s
FAIL

```

stderr：

```text
/tmp/gomark-trusted-test.patch:7: trailing whitespace.
package service_test
/tmp/gomark-trusted-test.patch:8: trailing whitespace.

/tmp/gomark-trusted-test.patch:9: trailing whitespace.
import (
/tmp/gomark-trusted-test.patch:10: trailing whitespace.
	"context"
/tmp/gomark-trusted-test.patch:11: trailing whitespace.
	"testing"
warning: squelched 43 whitespace errors
warning: 48 lines add whitespace errors.

```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/service -run TestCampaignPublishKeepsReviewWhenAuditFails -count=1
--- FAIL: TestCampaignPublishKeepsReviewWhenAuditFails (0.21s)
    campaign_audit_atomicity_task_test.go:33: campaign state=published
FAIL
FAIL	github.com/zhanglei10281852-gif/windsea/internal/service	0.453s
FAIL

```

stderr：

```text
/tmp/gomark-trusted-test.patch:7: trailing whitespace.
package service_test
/tmp/gomark-trusted-test.patch:8: trailing whitespace.

/tmp/gomark-trusted-test.patch:9: trailing whitespace.
import (
/tmp/gomark-trusted-test.patch:10: trailing whitespace.
	"context"
/tmp/gomark-trusted-test.patch:11: trailing whitespace.
	"testing"
warning: squelched 43 whitespace errors
warning: 48 lines add whitespace errors.

```

## 通过条件

When audit persistence fails the campaign remains in review and returns an error, while a valid audit allows the campaign to publish.
