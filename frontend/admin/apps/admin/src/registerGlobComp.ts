import type { App } from 'vue';

import {
  Button,
  Card,
  Col,
  Divider,
  Dropdown,
  Form,
  Input,
  InputNumber,
  Layout,
  Menu,
  Modal,
  Popconfirm,
  Row,
  Select,
  Space,
  Switch,
  Tabs,
  Tag,
  Tree,
} from 'ant-design-vue';

/**
 * 注册全局组件
 * @param app
 */
export function registerGlobComp(app: App) {
  app
    .use(Input)
    .use(Col)
    .use(Row)
    .use(InputNumber)
    .use(Button)
    .use(Layout)
    .use(Space)
    .use(Card)
    .use(Form)
    .use(Switch)
    .use(Select)
    .use(Popconfirm)
    .use(Dropdown)
    .use(Tag)
    .use(Tabs)
    .use(Divider)
    .use(Menu)
    .use(Modal)
    .use(Tree);
}
