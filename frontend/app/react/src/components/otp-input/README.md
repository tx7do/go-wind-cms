# VerificationCodeInput 验证码输入组件

基于 Ant Design Input.OTP 实现的验证码输入组件，支持数字和文本类型。

## 特性

- ✅ 支持自定义验证码长度（默认 6 位）
- ✅ 支持受控和非受控模式
- ✅ 支持数字和文本两种类型
- ✅ 完成输入自动触发回调
- ✅ 支持禁用状态
- ✅ 基于 Ant Design，样式美观

## 安装使用

```typescript
import {VerificationCodeInput} from '@/components/otp-input';
```

## API

### Props

| 属性 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| length | 验证码长度 | `number` | `6` |
| value | 受控模式的值 | `string` | - |
| onChange | 值变化时的回调 | `(value: string) => void` | - |
| onComplete | 完成输入时的回调 | `(value: string) => void` | - |
| disabled | 是否禁用 | `boolean` | `false` |
| className | 自定义类名 | `string` | - |
| type | 输入框类型 | `'number' \| 'text'` | `'number'` |

## 使用示例

### 基本用法（6 位数字）

```tsx
<VerificationCodeInput />
```

### 自定义长度

```tsx
// 4 位验证码
<VerificationCodeInput length={4} />
```

### 受控模式

```tsx
const [code, setCode] = useState('');

<VerificationCodeInput
  value={code}
  onChange={setCode}
/>
```

### 完成回调

```tsx
<VerificationCodeInput
  length={6}
  onComplete={(value) => {
    console.log('验证码输入完成:', value);
    // 在这里进行验证逻辑
  }}
/>
```

### 文本类型（支持字母和数字）

```tsx
<VerificationCodeInput
  type="text"
  value={code}
  onChange={setCode}
/>
```

### 禁用状态

```tsx
<VerificationCodeInput disabled />
```

## 实际场景示例

### 手机验证

```tsx
import React, {useState} from 'react';
import {VerificationCodeInput} from '@/components/otp-input';
import {Button} from 'antd';

const PhoneVerification: React.FC = () => {
    const [code, setCode] = useState('');
    const [loading, setLoading] = useState(false);

    const handleVerify = async (verificationCode: string) => {
        setLoading(true);
        try {
            // 调用验证 API
            await verifySmsCode(verificationCode);
            console.log('验证成功');
        } catch (error) {
            console.error('验证失败', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h3>请输入发送到您手机的 6 位验证码</h3>
            <VerificationCodeInput
                length={6}
                value={code}
                onChange={setCode}
                onComplete={handleVerify}
            />
            <Button 
                type="primary" 
                loading={loading}
                onClick={() => handleVerify(code)}
            >
                验证
            </Button>
        </div>
    );
};
```

### 邮箱验证

```tsx
import React, {useState} from 'react';
import {VerificationCodeInput} from '@/components/otp-input';

const EmailVerification: React.FC = () => {
    const [code, setCode] = useState('');

    return (
        <div>
            <h3>邮箱验证</h3>
            <p>请输入发送到您邮箱的验证码</p>
            <VerificationCodeInput
                type="text"
                length={6}
                value={code}
                onChange={setCode}
                onComplete={(value) => {
                    console.log('开始验证邮箱:', value);
                    // 执行邮箱验证逻辑
                }}
            />
        </div>
    );
};
```

## 注意事项

1. **受控模式**：如果使用 `value` 属性，必须同时提供 `onChange` 回调
2. **自动提交**：使用 `onComplete` 可以在用户输入完成后自动触发验证
3. **类型选择**：
   - `type="number"`：仅允许数字输入（适合短信验证码）
   - `type="text"`：允许字母和数字（适合邮箱验证码或混合验证码）
4. **长度限制**：验证码长度建议设置为 4-8 位，过短不安全，过长用户体验差

## 相关文件

- `VerificationCodeInput.tsx` - 组件实现
- `VerificationCodeInput.module.css` - 组件样式
- `index.ts` - 组件导出

## 依赖

- Ant Design (`Input.OTP`)
- next-intl (国际化)
