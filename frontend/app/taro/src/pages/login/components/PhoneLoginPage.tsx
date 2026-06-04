import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Input, Button, Text} from '@tarojs/components';export default function PhoneLoginPage() {
  const {t} = useTranslation();

  const [phone, setPhone] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [codeSent, setCodeSent] = useState(false);
  const [countdown, setCountdown] = useState(0);

  useEffect(() => {
    let timer: number;
    if (countdown > 0) {
      timer = window.setTimeout(() => setCountdown(countdown - 1), 1000);
    } else if (codeSent && countdown === 0) {
      // 使用 setTimeout 延迟执行，避免在 effect 中同步调用 setState
      setTimeout(() => setCodeSent(false));
    }
    return () => clearTimeout(timer);
  }, [countdown, codeSent]);

  const handleSendCode = () => {
    if (!phone) return;
    setCodeSent(true);
    setCountdown(60);
  };

  const handleLogin = () => {
    if (!phone || !verificationCode) return;
    console.log('登录信息:', {phone, verificationCode});
  };

  return (
    <View className='login-form'>
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.phone')}</Text>
        <Input
          type='number'
          value={phone}
          onInput={(e) => setPhone(e.detail.value)}
          placeholder={t('authentication.login.placeholder_phone')}
          className='input-field'
        />
      </View>

      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.code')}</Text>
        <View className='code-input-row'>
          <Input
            type='text'
            value={verificationCode}
            onInput={(e) => setVerificationCode(e.detail.value)}
            placeholder={t('authentication.login.placeholder_code')}
            className='input-field'
            maxlength={6}
          />
          <Button
            className='send-code-btn'
            disabled={codeSent}
            onClick={handleSendCode}
          >
            {codeSent ? `${countdown}s` : t('authentication.register.send_code')}
          </Button>
        </View>
      </View>

      <Button className='login-button' onClick={handleLogin}>
        {t('authentication.login.login')}
      </Button>
    </View>
  );
}
