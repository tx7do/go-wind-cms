// ⚠️ 占位语言资源（阿拉伯语 / Arabic）
//
// 当前文件中的文案为 en-US 的逐字副本，仅用于：
//   1. 让 ar locale 的 key 结构与其它语言一致，避免 missing-key 警告；
//   2. 在真实阿拉伯语翻译到位前，让 RTL 布局可以联调验证。
//
// 待真实阿拉伯语翻译就绪后，逐个替换下面的值为阿拉伯语译文即可。
// 注意：本目录下每个 *.json 的 key 必须与 en-US / zh-CN 保持一致。

import app from './app.json'
import authentication from './authentication.json'
import cms from './cms.json'
import comment from './comment.json'
import common from './common.json'
import component from './component.json'
import enum_ from './enum.json'
import menu from './menu.json'
import navbar from './navbar.json'
import page from './page.json'
import preferences from './preferences.json'
import settings from './settings.json'
import ui from './ui.json'

export default {
    app,
    authentication,
    cms,
    comment,
    common,
    component,
    enum: enum_,
    menu,
    navbar,
    page,
    preferences,
    settings,
    ui,
}
