# GoWind Content Hub｜风行，开箱即用的企业级前后端一体内容平台

风行（GoWind HCH）是一款开箱即用的企业级Golang全栈Headless内容平台（HCH=Headless Content
Hub，无头内容中枢），为企业提供灵活、可扩展的全域内容管理与分发解决方案。

[English](./README.en-US.md) | **中文** | [日本語](./README.ja-JP.md)

## 演示地址

| 演示类型          | 访问地址                                                                                 |
|---------------|--------------------------------------------------------------------------------------|
| 后端管理前端        | [https://admin.cms.gowind.cloud](https://admin.cms.gowind.cloud)                     |
| 后端API Swagger | [https://api.admin.cms.gowind.cloud/docs/](https://api.admin.cms.gowind.cloud/docs/) |
| 默认账号密码        | `admin` / `admin`（所有演示地址通用）                                                          |
| 前台API Swagger | [https://api.cms.gowind.cloud/docs/](https://api.cms.gowind.cloud/docs/)             |
| 前台Vue端        | [https://cms.gowind.cloud](https://cms.gowind.cloud)                                 |
| 前台React端        | [https://react.cms.gowind.cloud](https://react.cms.gowind.cloud)                                 |
| 前台Taro端        | [https://taro.cms.gowind.cloud](https://taro.cms.gowind.cloud)                                 |

## 风行·核心技术栈

秉持高效、稳定、可扩展的技术选型理念，支持多前端技术栈适配，系统核心技术栈如下：

- 后端基于 [Golang](https://go.dev/) + [go-kratos](https://go-kratos.dev/) + [wire](https://github.com/google/wire) + [ent](https://entgo.io/docs/getting-started/)
- 管理后台前端基于 [Vue](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) + [Ant Design Vue](https://antdv.com/) + [Vben Admin](https://doc.vben.pro/)
- 前台展示前端支持多技术栈，当前已适配 [Vue3](https://vuejs.org/) + [Naive UI](https://www.naiveui.com/)、[React](https://react.dev/) + [Next.js](https://nextjs.org/) + [Ant Design](https://ant.design/) 、[Taro](https://docs.taro.zone/en/docs/)。

## 风行·核心功能列表

| 功能   | 说明                                                                       |
|------|--------------------------------------------------------------------------|
|多租户管理|企业级多租户架构，支持租户新增、启用 / 禁用、套餐配置与隔离管理；新增租户自动初始化部门、默认角色与管理员账号，支持一键登录租户后台。|
|用户管理|系统用户全生命周期管理，支持用户新增、编辑、启用 / 禁用、重置密码；支持绑定多角色、多部门与上级主管，可设置 / 取消主管身份，支持一键代登录指定用户，支持高级条件查询。|
|角色管理|角色与角色分组统一管理，支持按角色关联用户，精细化配置菜单权限、接口权限与数据权限；支持批量添加 / 移除员工，灵活适配团队权限划分。|
|权限管理|统一管理权限分组、菜单节点与权限点，树形结构展示权限体系，支持按钮级细粒度权限控制，与接口、菜单联动实现完整权限闭环。|
|菜单管理|可视化配置系统菜单目录、菜单页面与功能按钮，支持自定义菜单图标、排序、路由与权限标识，前端菜单自动根据权限动态渲染。|
|部门管理|组织架构部门树形管理，支持多级部门创建、编辑、排序与联动绑定用户，清晰划分企业组织层级。|
|内容建模|可视化自定义内容模型与字段类型，支持文本、数字、富文本、图片、文件、关联等字段，灵活适配文章、产品、公告、素材等各类业务内容结构。|
|内容管理|统一管理各类内容数据，支持内容新增、编辑、发布 / 下架、置顶、排序、回收站与批量操作；支持富文本 / Markdown 编辑，图片与附件一键上传。|
|分类管理|多级内容分类树形管理，支持分类新增、编辑、排序、禁用，支持绑定内容模型，前台可按分类快速筛选内容。|
|标签管理|内容标签统一管理，支持标签新增、编辑、删除与关联内容，支持按标签检索与聚合展示。|
|评论管理|管理用户评论与互动内容，支持查看、审核、删除、回复、屏蔽违规评论，可按内容、用户、时间筛选查询。|
|多语言管理|原生多语言国际化支持，支持语种新增、启用 / 禁用，内容、菜单、提示信息统一翻译管理，无缝支撑出海与跨境业务。|
|站点管理|支持多站点独立配置，每个租户可创建多个站点，独立配置域名、标题、Logo、SEO 信息与展示风格。|
|站点配置|可视化站点系统配置，支持基础信息、SEO 优化、上传限制、缓存策略、邮件 / 短信等全局参数配置，配置实时生效。|
|文件资源管理|统一管理图片、文档、视频等文件资源，支持上传至本地或 OSS 云存储，支持文件预览、复制地址、下载、删除、分组管理与大图查看。|
|字典管理|系统数据字典大类与子项管理，支持联动查询、服务端排序、导入导出，用于下拉选项、状态标识等全局通用数据维护。|
|接口管理|后端接口统一维护与自动同步，支持树形展示接口列表，可配置接口请求参数、响应结果与操作日志记录，用于权限点绑定。|
|任务调度|定时任务管理，支持任务新增、编辑、删除、启动、暂停、立即执行，查看任务执行记录与运行日志，保障定时业务自动运行。|
|消息通知|消息分类与消息推送管理，支持多级消息分类，可向指定用户发送消息，查看消息已读状态与已读时间。|
|站内信|个人站内消息中心，支持消息查看、删除、单条标为已读、一键全部已读，实现系统与用户间消息触达。|
|个人中心|个人信息查看与编辑，头像、昵称修改，登录密码重置，查看登录记录与账号安全信息。|
|缓存管理|系统缓存实时查询与管理，支持按缓存键精准清除、批量清理，保障系统配置与数据实时更新。|
|登录日志|记录所有用户登录成功 / 失败日志，包含登录账号、IP、归属地、设备、时间，支持查询与导出，便于安全审计。|
|操作日志|全链路用户操作日志记录，包含正常与异常操作，记录操作人、IP、归属地、请求参数与结果，支持详情查看与追溯。|
|Headless API|API 优先设计，提供完整前后台 OpenAPI 接口，支持内容查询、创建、更新、删除，适配 Vue、React、Taro、小程序等多端调用。|

## 风行·后台截图展示

<table>
    <tr>
        <td><img src="./docs/images/admin_login.png" alt="后台用户登录界面"/></td>
        <td><img src="./docs/images/admin_register.png" alt="后台用户注册界面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_post_list.png" alt="后台帖子列表"/></td>
        <td><img src="./docs/images/admin_post_edit.png" alt="后台帖子编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_category_list.png" alt="后台分类列表"/></td>
        <td><img src="./docs/images/admin_category_edit.png" alt="后台分类编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_tag_list.png" alt="后台标签列表"/></td>
        <td><img src="./docs/images/admin_tag_edit.png" alt="后台标签编辑"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_comment_list.png" alt="后台评论列表"/></td>
        <td><img src="./docs/images/admin_site_list.png" alt="后台站点列表"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_site_setting_list.png" alt="后台站点配置列表"/></td>
    </tr>
</table>

## 风行·前台截图展示

<table>
    <tr>
        <td><img src="./docs/images/react_app_login.png" alt="React前台用户登录界面"/></td>
        <td><img src="./docs/images/react_app_register.png" alt="React前台用户注册界面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_homepage.png" alt="React前台主页"/></td>
        <td><img src="./docs/images/react_app_about.png" alt="React前台关于页面"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_post_list.png" alt="React前台帖子列表"/></td>
        <td><img src="./docs/images/react_app_post_detail.png" alt="React前台帖子详情页"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_category_list.png" alt="React前台分类列表"/></td>
        <td><img src="./docs/images/react_app_category_detail.png" alt="React前台分类详情页"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_tag_list.png" alt="React前台标签列表"/></td>
        <td><img src="./docs/images/react_app_tag_detail.png" alt="React前台标签详情页"/></td>
    </tr>
</table>

## 联系我们

- 微信个人号：`yang_lin_bo`（备注：`go-wind-cms`）
- 掘金专栏：[go-wind-cms](https://juejin.cn/column/7541283508041826367)

## [感谢JetBrains提供的免费GoLand & WebStorm](https://jb.gg/OpenSource)

[![avatar](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://jb.gg/OpenSource)
