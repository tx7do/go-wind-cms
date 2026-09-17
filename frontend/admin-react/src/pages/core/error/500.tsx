import Fallback from './Fallback.tsx';

/**
 * 500 服务器错误页面
 */
const InternalError = () => {
    return <Fallback status="500" />;
};

export default InternalError;
