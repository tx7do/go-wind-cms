interface Option {
    content: string; // 内容
    height: number; // 水印行高
    width: number; // 水印宽度
    rotate: number; // 旋转度数（可为负值）
    color: string; // 水印字体颜色
    fontSize: number; // 水印字体的大小
    opacity: number; // 水印透明度（0~1之间取值）
}

/**
 * 水印
 *
 * @deprecated 此 hook 当前在全项目内无任何调用方，且实现存在隐患：
 * 1) 未监听 window resize，窗口尺寸变化后水印不会重绘；
 * 2) 尺寸基于 document.body 计算，body 高度为 0 时（如路由切换瞬间）覆盖不全；
 * 3) 作为普通函数而非 React hook，调用方需自行管理生命周期。
 * 启用前请先补全上述问题，否则按需删除。
 */
export function useWatermark() {
    /**
     * 水印
     * @param options - 操作值
     */
    const Watermark = (options: Option) => {
        const {content, height, width, rotate, color, fontSize, opacity} = options;

        // 判断水印是否存在，如果存在，那么不执行
        if (document.getElementById('south_watermark') !== null) {
            return;
        }
        const TpLine = Math.floor(document.body?.clientWidth / width) * 2; // 一行显示几列

        // 安全构建单行水印：使用 createElement + textContent 避免 content 中的 HTML 注入（XSS）
        const lineDiv = document.createElement('div');
        lineDiv.style.whiteSpace = 'nowrap';
        for (let i = 0; i < TpLine; i++) {
            const span = document.createElement('span');
            span.style.cssText = `
        display: inline-block;
        line-height: ${height}px;
        width: ${width}px;
        text-align: center;
        transform:rotate( ${rotate}deg);
        color: ${color};
        font-size: ${fontSize}px;
        opacity: ${opacity};
        `;
            // textContent 不会被解析为 HTML，杜绝 XSS
            span.textContent = content;
            lineDiv.appendChild(span);
        }

        const TpColumn = Math.floor(document.body.clientHeight / height) * 2 || 4; // 一列显示几行

        const DivLayer = document.createElement('div');
        DivLayer.id = 'south_watermark'; // 给水印盒子添加类名
        DivLayer.style.position = 'fixed';
        DivLayer.style.top = '0px'; // 整体水印距离顶部距离
        DivLayer.style.left = '-100px'; // 改变整体水印的left值
        DivLayer.style.zIndex = '999999'; // 水印页面层级
        DivLayer.style.pointerEvents = 'none';
        DivLayer.style.userSelect = 'none';

        for (let i = 0; i < TpColumn; i++) {
            // 每行使用克隆节点，content 经 textContent 已安全
            DivLayer.appendChild(lineDiv.cloneNode(true));
        }

        document.body.appendChild(DivLayer); // 到页面中
    };

    /** 删除水印 */
    const RemoveWatermark = () => {
        // 判断水印是否存在，如果存在，那么执行
        if (document.getElementById('south_watermark') === null) {
            return;
        }
        if (document.getElementById('south_watermark') !== null) {
            const element = document.getElementById('south_watermark');
            document.body.removeChild(element as HTMLElement);
        }
    };

    return [Watermark, RemoveWatermark] as const;
}
