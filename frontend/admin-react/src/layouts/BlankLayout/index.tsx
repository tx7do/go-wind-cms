import {Outlet} from 'react-router-dom';

/**
 * 空白布局 - 用于 404/403/500 等错误页面
 * 仅提供 Outlet，不包含侧边栏、导航等
 * 主题配置由外层 ThemeProvider 统一提供
 */
const BlankLayout = () => {
    return <Outlet/>;
};

export default BlankLayout;