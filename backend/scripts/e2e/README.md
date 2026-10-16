# E2E 脚本（admin REST / app REST 级断言）

可重跑的业务流断言脚本，无需浏览器。浏览器级 GUI 验证仍走 Playwright 手工流程。

## 前置

- admin(:6600)/core/app(:6700) 已启动，且**清代理启动**（`HTTP_PROXY='' NO_PROXY='*'`）
- redis 容器名 `redis-server`（密码 `*Abcd123456`）
- `pip install pycryptodome`
- 管理员账号 admin/admin

## 脚本

| 脚本 | 覆盖流 |
|---|---|
| `test_page_review.py` | 页面发布审批全流：创建→提审→通过→发布→重复提审拒绝→下架（自动清审核态）→再提审→驳回 |
| `test_experiment_conversion.py` | 实验转化归因：建实验上线→匿名开页（参与+gwc_vid cookie）→带 cookie 提交 about 页表单→统计出转化 |

## 用法

```bash
cd backend/scripts/e2e
python test_page_review.py
python test_experiment_conversion.py
```

脚本自建自清数据，退出码 0=全过 / 1=有失败。

## 注意

- app 侧请求走 `curl --noproxy '*'`（common.app_curl），规避本机代理劫持
- 提交表单请求体为**顶层** `pageId/sectionId/data`（无 CRUD 的 data 包裹），`data` 值必须全字符串
- 新增端点部署后需在管理页点一次「接口同步」并绑角色，否则 E2E 401/403
