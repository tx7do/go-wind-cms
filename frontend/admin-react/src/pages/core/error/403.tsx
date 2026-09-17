import Fallback from './Fallback.tsx';

/**
 * 403 禁止访问页面
 */
const Forbidden = () => {
    return <Fallback status="403" />;
};

export default Forbidden;
