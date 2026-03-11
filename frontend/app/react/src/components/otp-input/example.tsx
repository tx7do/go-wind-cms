/**
 * VerificationCodeInput 使用示例
 *
 * @example
 * // 基本用法
 * <VerificationCodeInput />
 *
 * @example
 * // 自定义长度
 * <VerificationCodeInput length={4} />
 *
 * @example
 * // 受控模式
 * <VerificationCodeInput
 *   value={code}
 *   onChange={setCode}
 * />
 *
 * @example
 * // 完成回调
 * <VerificationCodeInput
 *   length={6}
 *   onComplete={(value) => console.log('验证码输入完成:', value)}
 * />
 *
 * @example
 * // 文本类型（支持字母和数字）
 * <VerificationCodeInput type="text" />
 *
 * @example
 * // 禁用状态
 * <VerificationCodeInput disabled />
 */

import React, {useState} from 'react';
import {VerificationCodeInput} from '@/components/VerificationCodeInput';

export const VerificationCodeInputExample: React.FC = () => {
    const [code1, setCode1] = useState('');
    const [code2, setCode2] = useState('');
    const [code3, setCode3] = useState('');

    return (
        <div style={{display: 'flex', flexDirection: 'column', gap: '24px'}}>
            {/* 基本用法 */}
            <div>
                <h3>基本用法（6 位数字）</h3>
                <VerificationCodeInput/>
            </div>

            {/* 自定义长度 */}
            <div>
                <h3>自定义长度（4 位）</h3>
                <VerificationCodeInput length={4}/>
            </div>

            {/* 受控模式 */}
            <div>
                <h3>受控模式</h3>
                <p>当前值：{code1 || '空'}</p>
                <VerificationCodeInput
                    value={code1}
                    onChange={setCode1}
                />
            </div>

            {/* 完成回调 */}
            <div>
                <h3>完成回调</h3>
                <VerificationCodeInput
                    length={6}
                    onComplete={(value) => {
                        console.log('验证码输入完成:', value);
                        alert(`验证码输入完成：${value}`);
                    }}
                />
            </div>

            {/* 文本类型 */}
            <div>
                <h3>文本类型（支持字母和数字）</h3>
                <VerificationCodeInput
                    type="text"
                    value={code2}
                    onChange={setCode2}
                />
            </div>

            {/* 禁用状态 */}
            <div>
                <h3>禁用状态</h3>
                <VerificationCodeInput disabled/>
            </div>

            {/* 实际场景示例 */}
            <div>
                <h3>实际场景：手机验证</h3>
                <p>请输入发送到您手机的 6 位验证码</p>
                <VerificationCodeInput
                    length={6}
                    value={code3}
                    onChange={setCode3}
                    onComplete={(value) => {
                        // 自动提交验证
                        console.log('自动验证:', value);
                    }}
                />
            </div>
        </div>
    );
};

export default VerificationCodeInputExample;
